<script lang="ts">
	import { PUBLIC_API_URL } from '$env/static/public';
	import { toast } from 'svelte-sonner';
	import { Upload, X, Star, ImageUp, Video, Loader } from 'lucide-svelte';
	import type { Media } from '$api/media';

	interface Props {
		media: Media[];
		featuredImageId?: string | null;
		onMediaChange: (mediaIds: string[]) => void;
		onFeaturedImageChange: (mediaId: string | null) => void;
	}

	let {
		media = $bindable([]),
		featuredImageId = $bindable(null),
		onMediaChange,
		onFeaturedImageChange
	}: Props = $props();

	let uploading = $state(false);
	let dragOver = $state(false);
	let draggedIndex = $state<number | null>(null);

	// Handle file upload
	async function handleFileUpload(files: FileList | null) {
		if (!files || files.length === 0) return;

		uploading = true;

		try {
			const uploadPromises = Array.from(files).map(async (file) => {
				const formData = new FormData();
				formData.append('file', file);

				const response = await fetch(`${PUBLIC_API_URL}/api/v1/media/upload`, {
					method: 'POST',
					credentials: 'include',
					body: formData
				});

				if (!response.ok) {
					throw new Error(`Failed to upload ${file.name}`);
				}

				return await response.json();
			});

			const uploadedMedia = await Promise.all(uploadPromises);

			// Add new media to the list
			media = [...media, ...uploadedMedia];

			// Update parent with new media IDs
			onMediaChange(media.map((m) => m.id));

			toast.success(`Uploaded ${uploadedMedia.length} file(s)`);
		} catch (error) {
			console.error('Upload error:', error);
			toast.error('Failed to upload files');
		} finally {
			uploading = false;
		}
	}

	function handleDrop(e: DragEvent) {
		e.preventDefault();
		dragOver = false;
		handleFileUpload(e.dataTransfer?.files || null);
	}

	function handleFileSelect(e: Event) {
		const target = e.target as HTMLInputElement;
		handleFileUpload(target.files);
		target.value = ''; // Reset input
	}

	function removeMedia(index: number) {
		const removedMedia = media[index];
		media = media.filter((_, i) => i !== index);

		// If removed media was featured, clear featured image
		if (featuredImageId === removedMedia.id) {
			featuredImageId = null;
			onFeaturedImageChange(null);
		}

		onMediaChange(media.map((m) => m.id));
		toast.success('Media removed');
	}

	function setFeaturedImage(mediaId: string) {
		featuredImageId = mediaId;
		onFeaturedImageChange(mediaId);
		toast.success('Featured image updated');
	}

	// Drag and drop reordering
	function handleDragStart(index: number) {
		draggedIndex = index;
	}

	function handleDragOver(e: DragEvent, index: number) {
		e.preventDefault();
		if (draggedIndex === null || draggedIndex === index) return;

		const newMedia = [...media];
		const [draggedItem] = newMedia.splice(draggedIndex, 1);
		newMedia.splice(index, 0, draggedItem);

		media = newMedia;
		draggedIndex = index;
		onMediaChange(media.map((m) => m.id));
	}

	function handleDragEnd() {
		draggedIndex = null;
	}

	function formatFileSize(bytes: number): string {
		if (bytes < 1024) return bytes + ' B';
		if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
		return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
	}

	function isImage(m: Media): boolean {
		return m.mime_type.startsWith('image/');
	}

	function isVideo(m: Media): boolean {
		return m.mime_type.startsWith('video/');
	}
</script>

<div class="space-y-6">
	<!-- Upload Area -->
	<div
		aria-pressed="false"
		class="relative cursor-pointer rounded-lg border-2 border-dashed p-8 text-center transition-colors"
		class:border-primary={dragOver}
		class:bg-primary={dragOver}
		class:border-gray-300={!dragOver}
		ondragover={(e) => {
			e.preventDefault();
			dragOver = true;
		}}
		ondragleave={() => {
			dragOver = false;
		}}
		ondrop={handleDrop}
		onclick={() => document.getElementById('file-input')?.click()}
	>
		<input
			id="file-input"
			type="file"
			accept="image/*,video/*"
			multiple
			onchange={handleFileSelect}
			class="hidden"
		/>

		{#if uploading}
			<div class="flex flex-col items-center gap-4">
				<Loader class="h-12 w-12 animate-spin text-primary" />
				<p class="text-lg font-medium">Uploading...</p>
			</div>
		{:else}
			<Upload class="mx-auto mb-4 h-12 w-12 text-gray-400" />
			<div class="space-y-2">
				<p class="text-lg font-medium">Drop files here or click to upload</p>
				<p class="text-sm text-gray-500">
					Images will be converted to WebP • Videos will be optimized
				</p>
				<p class="text-xs text-gray-400">Images: 50MB max • Videos: 200MB max</p>
			</div>
		{/if}
	</div>

	<!-- Media Grid -->
	{#if media.length > 0}
		<div>
			<div class="mb-4 flex items-center justify-between">
				<h3 class="text-lg font-semibold">Project Media ({media.length})</h3>
				<p class="text-sm text-gray-500">Drag to reorder • Click star to set featured</p>
			</div>

			<div class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4">
				{#each media as item, index (item.id)}
					<div
						draggable="true"
						ondragstart={() => handleDragStart(index)}
						ondragover={(e) => handleDragOver(e, index)}
						ondragend={handleDragEnd}
						class="group relative cursor-move overflow-hidden rounded-lg border border-gray-200 bg-white shadow-sm transition-all hover:shadow-md"
						class:ring-2={draggedIndex === index}
						class:ring-primary={draggedIndex === index}
					>
						<!-- Featured Badge -->
						{#if featuredImageId === item.id}
							<div class="absolute top-2 left-2 z-10 rounded-full bg-yellow-400 p-1.5 shadow-md">
								<Star class="h-4 w-4 fill-white text-white" />
							</div>
						{/if}

						<!-- Image/Video Preview -->
						<div class="aspect-square bg-gray-100">
							{#if isImage(item)}
								<img
									src={item.thumbnail_url}
									alt={item.original_filename}
									class="h-full w-full object-cover"
								/>
							{:else if isVideo(item)}
								<div class="flex h-full w-full items-center justify-center">
									<Video class="h-12 w-12 text-gray-400" />
								</div>
							{:else}
								<div class="flex h-full w-full items-center justify-center">
									<ImageUp class="h-12 w-12 text-gray-400" />
								</div>
							{/if}
						</div>

						<!-- Info & Actions -->
						<div class="p-3">
							<p class="truncate text-sm font-medium" title={item.original_filename}>
								{item.original_filename}
							</p>
							<p class="text-xs text-gray-500">
								{formatFileSize(item.size_bytes)}
								{#if item.width && item.height}
									• {item.width}×{item.height}
								{/if}
							</p>

							<!-- Action Buttons -->
							<div class="mt-2 flex gap-2">
								<button
									type="button"
									onclick={(e) => {
										e.stopPropagation();
										setFeaturedImage(item.id);
									}}
									class="flex-1 rounded bg-gray-100 px-2 py-1 text-xs transition-colors hover:bg-gray-200"
									class:bg-yellow-100={featuredImageId === item.id}
									class:text-yellow-700={featuredImageId === item.id}
								>
									<Star
										class={`mx-auto h-3 w-3 ${featuredImageId === item.id ? 'fill-current' : ''}`}
									/>
								</button>
								<button
									type="button"
									onclick={(e) => {
										e.stopPropagation();
										removeMedia(index);
									}}
									class="flex-1 rounded bg-red-100 px-2 py-1 text-xs text-red-600 transition-colors hover:bg-red-200"
								>
									<X class="mx-auto h-3 w-3" />
								</button>
							</div>
						</div>
					</div>
				{/each}
			</div>
		</div>
	{:else}
		<div class="rounded-lg border border-dashed border-gray-300 p-8 text-center">
			<ImageUp class="mx-auto mb-3 h-12 w-12 text-gray-400" />
			<p class="text-sm text-gray-500">No media added yet</p>
		</div>
	{/if}
</div>
