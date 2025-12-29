<!-- frontend/src/lib/components/ui/form/fields/MultiMediaField.svelte -->
<script lang="ts">
	import { ImageUp, GripVertical, ArrowUp, ArrowDown } from 'lucide-svelte';
	import type { MediaFieldProps } from '../types';
	import { PUBLIC_API_URL } from '$env/static/public';
	import type { Media } from '$api/media';

	interface Props extends MediaFieldProps {
		value: string[];
	}

	let {
		field,
		value = $bindable([]),
		error,
		disabled = false,
		mediaCache,
		onMediaPickerOpen,
		onMediaRemove,
		onMediaMove,
		onchange
	}: Props = $props();

	const mediaList = $derived(value.map((id) => mediaCache.get(id)).filter((m) => m !== undefined));

	// Drag state for reordering
	let draggedIndex = $state<number | null>(null);
	let dragOverIndex = $state<number | null>(null);

	// Drag and drop upload state
	let isDraggingFile = $state(false);
	let isUploading = $state(false);
	let uploadProgress = $state(0);
	let uploadingCount = $state(0);
	let uploadedCount = $state(0);

	function handleRemove(mediaId: string) {
		onMediaRemove?.(field.name, mediaId);
	}

	function handleMove(fromIndex: number, toIndex: number) {
		if (toIndex < 0 || toIndex >= value.length) return;
		onMediaMove?.(field.name, fromIndex, toIndex);
	}

	// Drag and drop handlers for reordering
	function handleDragStart(e: DragEvent, index: number) {
		if (disabled) return;
		draggedIndex = index;
		if (e.dataTransfer) {
			e.dataTransfer.effectAllowed = 'move';
			e.dataTransfer.setData('text/plain', index.toString());
		}
	}

	function handleDragOver(e: DragEvent, index: number) {
		if (disabled || draggedIndex === null) return;
		e.preventDefault();
		if (e.dataTransfer) {
			e.dataTransfer.dropEffect = 'move';
		}
		dragOverIndex = index;
	}

	function handleDragLeave() {
		dragOverIndex = null;
	}

	function handleDrop(e: DragEvent, dropIndex: number) {
		if (disabled || draggedIndex === null) return;
		e.preventDefault();

		if (draggedIndex !== dropIndex) {
			handleMove(draggedIndex, dropIndex);
		}

		draggedIndex = null;
		dragOverIndex = null;
	}

	function handleDragEnd() {
		draggedIndex = null;
		dragOverIndex = null;
	}

	// Drag and drop handlers for file upload
	function handleFileDragEnter(e: DragEvent) {
		e.preventDefault();
		e.stopPropagation();
		if (!disabled && e.dataTransfer?.types.includes('Files')) {
			isDraggingFile = true;
		}
	}

	function handleFileDragOver(e: DragEvent) {
		e.preventDefault();
		e.stopPropagation();
		if (e.dataTransfer) {
			e.dataTransfer.dropEffect = 'copy';
		}
	}

	function handleFileDragLeave(e: DragEvent) {
		e.preventDefault();
		e.stopPropagation();
		const rect = (e.currentTarget as HTMLElement).getBoundingClientRect();
		const x = e.clientX;
		const y = e.clientY;
		if (x <= rect.left || x >= rect.right || y <= rect.top || y >= rect.bottom) {
			isDraggingFile = false;
		}
	}

	async function handleFileDrop(e: DragEvent) {
		e.preventDefault();
		e.stopPropagation();
		isDraggingFile = false;

		if (disabled || !e.dataTransfer?.files.length) return;

		const files = Array.from(e.dataTransfer.files);
		await uploadFiles(files);
	}

	async function uploadFiles(files: File[]) {
		if (!files.length) return;

		isUploading = true;
		uploadingCount = files.length;
		uploadedCount = 0;
		uploadProgress = 0;

		for (const file of files) {
			try {
				// Check file type
				const isImage = file.type.startsWith('image/');
				const isVideo = file.type.startsWith('video/');

				if (!isImage && !isVideo) {
					console.warn(`Skipping ${file.name}: not an image or video`);
					continue;
				}

				// Check file size
				const maxSize = isVideo ? 200 : 50;
				const fileSizeMB = file.size / (1024 * 1024);

				if (fileSizeMB > maxSize) {
					console.warn(`Skipping ${file.name}: file size exceeds ${maxSize}MB`);
					continue;
				}

				const uploadedMedia = await uploadFile(file);

				// Add to value array and cache
				const newValue = [...value, uploadedMedia.id];
				value = newValue;
				mediaCache.set(uploadedMedia.id, uploadedMedia);

				// ✅ CRITICAL: Notify parent form of the change after each upload
				onchange(newValue);

				uploadedCount++;
				uploadProgress = Math.round((uploadedCount / uploadingCount) * 100);
			} catch (err) {
				console.error(`Failed to upload ${file.name}:`, err);
			}
		}

		// Reset upload state after a short delay
		setTimeout(() => {
			isUploading = false;
			uploadProgress = 0;
			uploadingCount = 0;
			uploadedCount = 0;
		}, 1000);
	}

	async function uploadFile(file: File): Promise<Media> {
		const formData = new FormData();
		formData.append('file', file);

		const response = await fetch(`${PUBLIC_API_URL}/api/v1/media/upload`, {
			method: 'POST',
			credentials: 'include',
			body: formData
		});

		if (!response.ok) {
			const error = await response.json();
			throw new Error(error.error || 'Upload failed');
		}

		return response.json();
	}
</script>

<div class="space-y-2">
	{#if field.label}
		<label for={field.name} class="block text-sm font-medium">
			{field.label}
			{#if field.required}
				<span class="text-red-500">*</span>
			{/if}
		</label>
	{/if}

	{#if mediaList.length > 0}
		<div class="space-y-2" role="list">
			{#each mediaList as media, index}
				{@const mediaId = value[index]}
				{@const isDragging = draggedIndex === index}
				{@const isDragOver = dragOverIndex === index}

				<div
					role="listitem"
					draggable={!disabled}
					ondragstart={(e) => handleDragStart(e, index)}
					ondragover={(e) => handleDragOver(e, index)}
					ondragleave={handleDragLeave}
					ondrop={(e) => handleDrop(e, index)}
					ondragend={handleDragEnd}
					class="flex items-center gap-3 rounded-lg border p-3 transition-all {isDragging
						? 'border-primary bg-gray-700 opacity-50'
						: isDragOver
							? 'border-primary bg-gray-700'
							: 'border-gray-700 bg-gray-800'}"
				>
					<button
						type="button"
						{disabled}
						class="cursor-move text-gray-400 hover:text-gray-300 disabled:cursor-not-allowed disabled:opacity-50"
						title="Drag to reorder"
						tabindex="-1"
					>
						<GripVertical class="h-5 w-5" />
					</button>

					<img
						src={media.thumbnail_url || media.url}
						alt={media.original_filename}
						class="h-16 w-16 rounded object-cover"
					/>

					<div class="flex-1">
						<p class="text-sm font-medium text-gray-100">{media.original_filename}</p>
						<p class="mt-1 text-xs text-gray-500">{media.mime_type}</p>
					</div>

					<div class="flex gap-1">
						<button
							type="button"
							onclick={() => handleMove(index, index - 1)}
							class="px-2 py-1 text-sm text-gray-300 hover:text-primary disabled:cursor-not-allowed disabled:opacity-30"
							disabled={disabled || index === 0}
							title="Move up"
						>
							<ArrowUp size={18} />
						</button>
						<button
							type="button"
							onclick={() => handleMove(index, index + 1)}
							class="px-2 py-1 text-sm text-gray-300 hover:text-primary disabled:cursor-not-allowed disabled:opacity-30"
							disabled={disabled || index === mediaList.length - 1}
							title="Move down"
						>
							<ArrowDown size={18} />
						</button>
						<button
							type="button"
							{disabled}
							onclick={() => handleRemove(mediaId)}
							class="px-3 py-1 text-sm font-medium text-red-600 hover:text-red-700 disabled:cursor-not-allowed disabled:opacity-50"
						>
							Remove
						</button>
					</div>
				</div>
			{/each}
		</div>
	{/if}

	<div
		role="button"
		tabindex="0"
		id={field.name}
		aria-label={field.label || 'Upload media files'}
		onclick={() => !disabled && onMediaPickerOpen(field.name)}
		onkeydown={(e) => {
			if (!disabled && (e.key === 'Enter' || e.key === ' ')) {
				e.preventDefault();
				onMediaPickerOpen(field.name);
			}
		}}
		ondragenter={handleFileDragEnter}
		ondragover={handleFileDragOver}
		ondragleave={handleFileDragLeave}
		ondrop={handleFileDrop}
		class="w-full cursor-pointer rounded-lg border-2 p-4 text-center transition-colors {disabled
			? 'cursor-not-allowed opacity-50'
			: isDraggingFile
				? 'border-solid border-primary bg-primary/10'
				: 'border-dashed border-gray-300 hover:border-gray-400'}"
	>
		{#if isUploading}
			<div class="pointer-events-none space-y-2">
				<div
					class="mx-auto mb-2 h-8 w-8 animate-spin rounded-full border-4 border-gray-300 border-t-primary"
				></div>
				<p class="text-sm text-gray-600">
					Uploading {uploadedCount} of {uploadingCount}... {uploadProgress}%
				</p>
				<div class="mx-auto h-2 w-48 overflow-hidden rounded-full bg-gray-200">
					<div
						class="h-full bg-primary transition-all duration-300"
						style="width: {uploadProgress}%"
					></div>
				</div>
			</div>
		{:else}
			<div class="pointer-events-none">
				<ImageUp class="mx-auto mb-2 h-8 w-8 text-gray-400" />
				<p class="text-sm text-gray-600">
					{isDraggingFile
						? 'Drop files to upload'
						: mediaList.length > 0
							? 'Add more media'
							: 'Click to select or drag & drop multiple files'}
				</p>
				<p class="mt-1 text-xs text-gray-400">
					{field.required && mediaList.length === 0 ? 'At least one file is required' : 'Optional'} •
					Images (50MB) or Videos (200MB)
				</p>
			</div>
		{/if}
	</div>

	{#if error}
		<p class="text-sm text-red-600">{error}</p>
	{/if}
</div>
