<!-- frontend/src/lib/components/ui/form/fields/MediaField.svelte -->
<script lang="ts">
	import { ImageUp } from 'lucide-svelte';
	import type { MediaFieldProps } from '../types';
	import { PUBLIC_API_URL } from '$env/static/public';
	import type { Media } from '$api/media';

	interface Props extends MediaFieldProps {
		value: string | null;
	}

	let {
		field,
		value = $bindable(null),
		error,
		disabled = false,
		mediaCache,
		onMediaPickerOpen,
		onMediaRemove,
		onchange
	}: Props = $props();

	const hasMedia = $derived(value && mediaCache.has(value));
	const media = $derived(hasMedia ? mediaCache.get(value) : null);

	// Drag and drop state
	let isDragging = $state(false);
	let isUploading = $state(false);
	let uploadProgress = $state(0);

	function handleClear() {
		onMediaRemove?.(field.name, value || undefined);
	}

	// Drag and drop handlers
	function handleDragEnter(e: DragEvent) {
		e.preventDefault();
		e.stopPropagation();
		if (!disabled && e.dataTransfer?.types.includes('Files')) {
			isDragging = true;
		}
	}

	function handleDragOver(e: DragEvent) {
		e.preventDefault();
		e.stopPropagation();
		if (e.dataTransfer) {
			e.dataTransfer.dropEffect = 'copy';
		}
	}

	function handleDragLeave(e: DragEvent) {
		e.preventDefault();
		e.stopPropagation();
		// Only set isDragging to false if we're leaving the dropzone entirely
		const rect = (e.currentTarget as HTMLElement).getBoundingClientRect();
		const x = e.clientX;
		const y = e.clientY;
		if (x <= rect.left || x >= rect.right || y <= rect.top || y >= rect.bottom) {
			isDragging = false;
		}
	}

	async function handleDrop(e: DragEvent) {
		e.preventDefault();
		e.stopPropagation();
		isDragging = false;

		if (disabled || !e.dataTransfer?.files.length) return;

		const file = e.dataTransfer.files[0];
		await uploadFile(file);
	}

	async function uploadFile(file: File) {
		if (!file) return;

		// Check file type
		const isImage = file.type.startsWith('image/');
		const isVideo = file.type.startsWith('video/');

		if (!isImage && !isVideo) {
			alert('Please upload an image or video file');
			return;
		}

		// Check file size (50MB for images, 200MB for videos)
		const maxSize = isVideo ? 200 : 50;
		const fileSizeMB = file.size / (1024 * 1024);

		if (fileSizeMB > maxSize) {
			alert(`File size (${fileSizeMB.toFixed(1)}MB) exceeds ${maxSize}MB limit`);
			return;
		}

		isUploading = true;
		uploadProgress = 0;

		try {
			const formData = new FormData();
			formData.append('file', file);

			const xhr = new XMLHttpRequest();

			xhr.upload.addEventListener('progress', (e) => {
				if (e.lengthComputable) {
					uploadProgress = Math.round((e.loaded / e.total) * 100);
				}
			});

			const uploadedMedia = await new Promise<Media>((resolve, reject) => {
				xhr.addEventListener('load', () => {
					if (xhr.status === 201) {
						resolve(JSON.parse(xhr.responseText));
					} else {
						reject(new Error(JSON.parse(xhr.responseText).error || 'Upload failed'));
					}
				});
				xhr.addEventListener('error', () => {
					reject(new Error('Network error during upload'));
				});
				xhr.open('POST', `${PUBLIC_API_URL}/api/v1/media/upload`);
				xhr.withCredentials = true;
				xhr.send(formData);
			});

			// Update the value and cache
			value = uploadedMedia.id;
			mediaCache.set(uploadedMedia.id, uploadedMedia);

			// ✅ CRITICAL: Notify parent form of the change
			onchange(uploadedMedia.id);
		} catch (err) {
			console.error('Upload failed:', err);
			alert(err instanceof Error ? err.message : 'Upload failed');
		} finally {
			isUploading = false;
			uploadProgress = 0;
		}
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

	{#if hasMedia && media && !isUploading}
		<div class="flex items-center gap-3 rounded-lg border border-gray-700 bg-gray-800 p-3">
			<img
				src={media.thumbnail_url || media.url}
				alt={media.original_filename}
				class="h-16 w-16 rounded object-cover"
			/>
			<div class="flex-1">
				<p class="text-sm font-medium text-gray-100">{media.original_filename}</p>
				<p class="mt-1 text-xs text-gray-500">{media.mime_type}</p>
			</div>
			<div class="flex gap-2">
				<button
					type="button"
					{disabled}
					onclick={() => onMediaPickerOpen(field.name)}
					class="px-3 py-1.5 text-sm font-medium text-blue-600 hover:text-blue-700 disabled:opacity-50"
				>
					Change
				</button>
				<button
					type="button"
					{disabled}
					onclick={handleClear}
					class="px-3 py-1.5 text-sm font-medium text-red-600 hover:text-red-700 disabled:opacity-50"
				>
					Remove
				</button>
			</div>
		</div>
	{:else}
		<button
			type="button"
			id={field.name}
			aria-label={field.label || 'Upload media file'}
			{disabled}
			onclick={() => onMediaPickerOpen(field.name)}
			ondragenter={handleDragEnter}
			ondragover={handleDragOver}
			ondragleave={handleDragLeave}
			ondrop={handleDrop}
			class="w-full rounded-lg border-2 p-6 text-center transition-colors disabled:cursor-not-allowed disabled:opacity-50 {isDragging
				? 'border-solid border-primary bg-primary/10'
				: 'border-dashed border-gray-300 hover:border-gray-400'}"
		>
			{#if isUploading}
				<div class="space-y-2">
					<div
						class="mx-auto mb-2 h-12 w-12 animate-spin rounded-full border-4 border-gray-300 border-t-primary"
					></div>
					<p class="text-sm text-gray-600">Uploading... {uploadProgress}%</p>
					<div class="mx-auto h-2 w-48 overflow-hidden rounded-full bg-gray-200">
						<div
							class="h-full bg-primary transition-all duration-300"
							style="width: {uploadProgress}%"
						></div>
					</div>
				</div>
			{:else}
				<ImageUp class="mx-auto mb-2 h-12 w-12 text-gray-400" />
				<p class="text-sm text-gray-600">
					{isDragging ? 'Drop file to upload' : 'Click to select or drag & drop media'}
				</p>
				<p class="mt-1 text-xs text-gray-400">
					{field.required ? 'Required' : 'Optional'} • Images (50MB) or Videos (200MB)
				</p>
			{/if}
		</button>
	{/if}

	{#if error}
		<p class="text-sm text-red-600">{error}</p>
	{/if}
</div>
