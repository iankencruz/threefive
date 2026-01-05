<!-- frontend/src/lib/components/projects/ProjectMediaGallery.svelte -->
<script lang="ts">
	import { toast } from 'svelte-sonner';
	import { X, Star, ImageUp, Video, Plus, Image } from 'lucide-svelte';
	import MediaPicker from '$lib/components/ui/MediaPicker.svelte';
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

	let showMediaPicker = $state(false);
	let draggedIndex = $state<number | null>(null);

	// Handle media selection from picker
	function handleMediaSelect(mediaId: string, selectedMedia: Media) {
		const exists = media.some((m) => m.id === mediaId);
		if (!exists) {
			media = [...media, selectedMedia];
			onMediaChange(media.map((m) => m.id));
			toast.success('Media added');
		}
		showMediaPicker = false;
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

	function moveMediaUp(index: number) {
		if (index > 0) {
			const newMedia = [...media];
			[newMedia[index - 1], newMedia[index]] = [newMedia[index], newMedia[index - 1]];
			media = newMedia;
			onMediaChange(media.map((m) => m.id));
		}
	}

	function moveMediaDown(index: number) {
		if (index < media.length - 1) {
			const newMedia = [...media];
			[newMedia[index], newMedia[index + 1]] = [newMedia[index + 1], newMedia[index]];
			media = newMedia;
			onMediaChange(media.map((m) => m.id));
		}
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
	<!-- Header with Add Button -->
	<div class="flex items-center justify-between">
		<h3 class="text-lg font-semibold">Project Media ({media.length})</h3>
		<button
			type="button"
			onclick={() => (showMediaPicker = true)}
			class="flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary/90"
		>
			<Plus class="h-4 w-4" />
			Add Media
		</button>
	</div>

	<!-- Media Grid -->
	{#if media.length === 0}
		<div class="rounded-lg border-2 border-dashed border-gray-300 p-12 text-center">
			<div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-gray-100">
				<Image class="h-8 w-8 text-gray-400" />
			</div>
			<p class="mb-2 text-gray-600">No media added yet</p>
			<p class="mb-4 text-sm text-gray-500">
				Click "Add Media" to select images and videos from your media library
			</p>
			<button
				type="button"
				onclick={() => (showMediaPicker = true)}
				class="inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary/90"
			>
				<Plus class="h-4 w-4" />
				Add Media
			</button>
		</div>
	{:else}
		<div>
			<p class="mb-4 text-sm text-gray-500">Drag to reorder • Click star to set featured image</p>
			<div class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-4">
				{#each media as item, index (item.id)}
					<div
						draggable="true"
						ondragstart={() => handleDragStart(index)}
						ondragover={(e) => handleDragOver(e, index)}
						ondragend={handleDragEnd}
						class="group relative cursor-move overflow-hidden rounded-lg border border-gray-200 bg-white shadow-sm transition-all hover:border-primary hover:shadow-md"
						class:ring-2={draggedIndex === index}
						class:ring-primary={draggedIndex === index}
					>
						<!-- Featured Badge -->
						{#if featuredImageId === item.id}
							<div class="absolute top-2 left-2 z-10 rounded-full bg-yellow-400 p-1.5 shadow-md">
								<Star class="h-4 w-4 fill-white text-white" />
							</div>
						{/if}

						<!-- Position Badge -->
						<div
							class="absolute top-2 right-2 z-10 flex h-6 w-6 items-center justify-center rounded-full bg-primary text-xs font-bold text-white shadow-md"
						>
							{index + 1}
						</div>

						<!-- Image/Video Preview -->
						<div class="aspect-square bg-gray-100">
							{#if isImage(item)}
								<img
									src={item.thumbnail_url || item.medium_url || item.url}
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

						<!-- Controls Overlay -->
						<div
							class="absolute inset-x-2 top-10 z-10 flex flex-col items-end gap-1 opacity-0 transition-opacity group-hover:opacity-100"
						>
							<button
								type="button"
								onclick={(e) => {
									e.stopPropagation();
									removeMedia(index);
								}}
								class="flex h-7 w-7 items-center justify-center rounded bg-red-600 text-white hover:bg-red-700"
								title="Remove"
							>
								<X class="h-4 w-4" />
							</button>
							<button
								type="button"
								onclick={() => moveMediaUp(index)}
								disabled={index === 0}
								class="flex h-7 w-7 items-center justify-center rounded bg-white text-sm font-bold text-gray-900 hover:bg-gray-100 disabled:hidden"
								title="Move up"
							>
								↑
							</button>
							<button
								type="button"
								onclick={() => moveMediaDown(index)}
								disabled={index === media.length - 1}
								class="flex h-7 w-7 items-center justify-center rounded bg-white text-sm font-bold text-gray-900 hover:bg-gray-100 disabled:hidden"
								title="Move down"
							>
								↓
							</button>
						</div>

						<!-- Info -->
						<div class="p-2">
							<p class="truncate text-xs font-medium" title={item.original_filename}>
								{item.original_filename}
							</p>
							<p class="text-xs text-gray-500">
								{formatFileSize(item.size_bytes)}
								{#if item.width && item.height}
									• {item.width}×{item.height}
								{/if}
							</p>

							<!-- Action Buttons -->
							<div class="mt-2 flex gap-1">
								<button
									type="button"
									onclick={(e) => {
										e.stopPropagation();
										setFeaturedImage(item.id);
									}}
									class="group flex-1 rounded bg-yellow-300/30 px-2 py-1 text-xs transition-colors hover:bg-yellow-500"
									class:bg-yellow-100={featuredImageId === item.id}
									class:text-yellow-700={featuredImageId === item.id}
									title="Set as featured"
								>
									<Star
										class={`mx-auto h-3 w-3 group-hover:fill-yellow-700 ${featuredImageId === item.id ? 'fill-current' : 'text-yellow-700'}`}
									/>
								</button>
								<button
									type="button"
									onclick={(e) => {
										e.stopPropagation();
										removeMedia(index);
									}}
									class="flex-1 rounded bg-red-100 px-2 py-1 text-xs text-red-600 transition-colors hover:bg-red-200"
									title="Remove"
								>
									<X class="mx-auto h-3 w-3" />
								</button>
							</div>
						</div>
					</div>
				{/each}
			</div>
		</div>
	{/if}

	<!-- Validation warning -->
	{#if media.length === 0}
		<p class="flex items-center gap-2 text-sm text-amber-500">
			<span>!</span>
			No media added. Add at least one image to showcase your project.
		</p>
	{/if}
</div>

<!-- Media Picker Modal -->
{#if showMediaPicker}
	<MediaPicker
		show={showMediaPicker}
		onselect={handleMediaSelect}
		onclose={() => (showMediaPicker = false)}
	/>
{/if}
