<script lang="ts">
	import { PUBLIC_API_URL } from '$env/static/public';
	import { ImageUp, Video, Download, Trash2, X, CheckCircle, Loader2 } from 'lucide-svelte';
	import { onMount } from 'svelte';

	interface Media {
		id: string;
		filename: string;
		original_filename: string;
		mime_type: string;
		url: string;
		thumbnail_url?: string;
		width?: number;
		height?: number;
		size_bytes: number;
		created_at: string;
	}

	interface UploadState {
		file: File;
		progress: number;
		status: 'uploading' | 'processing' | 'success' | 'error';
		error?: string;
		media?: Media;
	}

	let media = $state<Media[]>([]);
	let loading = $state(false);
	let uploads = $state<UploadState[]>([]);
	let selectedMedia = $state<Media | null>(null);
	let dragOver = $state(false);

	onMount(async () => {
		await loadMedia();
	});

	async function loadMedia() {
		loading = true;
		try {
			const response = await fetch(`${PUBLIC_API_URL}/api/v1/media?page=1&limit=100`, {
				credentials: 'include'
			});

			if (response.ok) {
				const data = await response.json();
				media = data.data || [];
			}
		} catch (error) {
			console.error('Failed to load media:', error);
		} finally {
			loading = false;
		}
	}

	function handleFileSelect(e: Event) {
		const input = e.target as HTMLInputElement;
		if (input.files) {
			addFiles(Array.from(input.files));
		}
	}

	function handleDrop(e: DragEvent) {
		e.preventDefault();
		dragOver = false;

		if (e.dataTransfer?.files) {
			addFiles(Array.from(e.dataTransfer.files));
		}
	}

	function addFiles(files: File[]) {
		files.forEach((file) => {
			const uploadState: UploadState = {
				file,
				progress: 0,
				status: 'uploading'
			};
			uploads = [...uploads, uploadState];
			uploadFile(uploadState);
		});
	}

	// Replace the uploadFile function in your admin/media/+page.svelte

	async function uploadFile(uploadState: UploadState) {
		const formData = new FormData();
		formData.append('file', uploadState.file);

		try {
			const xhr = new XMLHttpRequest();

			xhr.upload.addEventListener('progress', (e) => {
				if (e.lengthComputable) {
					const progress = Math.round((e.loaded / e.total) * 100);

					// Find the index and update the specific upload state
					const index = uploads.findIndex((u) => u.file === uploadState.file);
					if (index !== -1) {
						uploads[index].progress = progress;

						if (progress === 100) {
							uploads[index].status = 'processing';
						}

						// Trigger reactivity
						uploads = uploads;
					}
				}
			});

			const uploadedMedia = await new Promise<Media>((resolve, reject) => {
				xhr.addEventListener('load', () => {
					if (xhr.status === 201) {
						resolve(JSON.parse(xhr.responseText));
					} else {
						const error = JSON.parse(xhr.responseText);
						reject(new Error(error.error || 'Upload failed'));
					}
				});

				xhr.addEventListener('error', () => {
					reject(new Error('Network error during upload'));
				});

				xhr.open('POST', `${PUBLIC_API_URL}/api/v1/media/upload`);
				xhr.withCredentials = true;
				xhr.send(formData);
			});

			// Find and update the upload state to success
			const index = uploads.findIndex((u) => u.file === uploadState.file);
			if (index !== -1) {
				uploads[index].status = 'success';
				uploads[index].media = uploadedMedia;
				uploads = uploads;
			}

			// Add to media list
			media = [uploadedMedia, ...media];

			// Remove from uploads after 2 seconds
			setTimeout(() => {
				uploads = uploads.filter((u) => u.file !== uploadState.file);
			}, 2000);
		} catch (error) {
			const index = uploads.findIndex((u) => u.file === uploadState.file);
			if (index !== -1) {
				uploads[index].status = 'error';
				uploads[index].error = error instanceof Error ? error.message : 'Upload failed';
				uploads = uploads;
			}
		}
	}

	async function deleteMedia(m: Media) {
		if (!confirm(`Delete ${m.original_filename}?`)) return;

		try {
			const response = await fetch(`${PUBLIC_API_URL}/api/v1/media/${m.id}`, {
				method: 'DELETE',
				credentials: 'include'
			});

			if (response.ok) {
				media = media.filter((item) => item.id !== m.id);
				if (selectedMedia?.id === m.id) {
					selectedMedia = null;
				}
			}
		} catch (error) {
			console.error('Delete failed:', error);
			alert('Failed to delete media');
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

	function getProcessingMessage(file: File): string {
		if (file.type.startsWith('image/')) {
			return 'Converting to WebP and generating thumbnail...';
		} else if (file.type.startsWith('video/')) {
			return 'Optimizing video and extracting thumbnail...';
		}
		return 'Processing file...';
	}

	function clearUploads() {
		uploads = [];
	}
</script>

<svelte:head>
	<title>Admin: Media Library</title>
</svelte:head>

<div class="min-h-screen bg-background py-8">
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<!-- Header -->
		<div class="mb-8">
			<h1 class="">Media Library</h1>
			<p class="mt-2">Upload and manage images and videos with automatic processing</p>
		</div>

		<!-- Upload Area -->
		<div class="mb-8 rounded-lg bg-surface p-6 shadow-sm">
			<h2 class="mb-4 text-xl font-semibold">Upload Media</h2>

			<div
				class="relative cursor-pointer rounded-lg border-2 border-dashed p-8 text-center transition-colors"
				class:border-blue-500={dragOver}
				class:bg-blue-50={dragOver}
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

				<ImageUp class="mx-auto mb-4 h-12 w-12 text-gray-400" />

				<div class="space-y-2">
					<p class="text-lg font-medium">Drop files here or click to upload</p>
					<p class="text-sm">Images will be converted to WebP • Videos will be optimized</p>
					<p class="text-xs text-gray-400">Images: 50MB max • Videos: 200MB max</p>
				</div>
			</div>

			<!-- Upload Progress -->
			{#if uploads.length > 0}
				<div class="mt-6 space-y-3">
					<div class="flex items-center justify-between">
						<h3 class="text-sm font-medium text-gray-700">
							Uploading {uploads.length}
							{uploads.length === 1 ? 'file' : 'files'}
						</h3>
						<button onclick={clearUploads} class="text-sm hover:text-gray-700"> Clear all </button>
					</div>

					{#each uploads as upload (upload.file.name)}
						<div class="rounded-lg border bg-surface p-4">
							<div class="flex items-start gap-4">
								<div class="flex-shrink-0">
									{#if upload.status === 'success' && upload.media?.thumbnail_url}
										<img
											src={upload.media.thumbnail_url}
											alt={upload.file.name}
											class="h-16 w-16 rounded object-cover"
										/>
									{:else}
										<div class="flex h-16 w-16 items-center justify-center rounded bg-gray-100">
											{#if upload.file.type.startsWith('image/')}
												<ImageUp class="h-8 w-8 text-gray-400" />
											{:else}
												<Video class="h-8 w-8 text-gray-400" />
											{/if}
										</div>
									{/if}
								</div>

								<div class="min-w-0 flex-1">
									<div class="mb-2 flex items-start justify-between gap-2">
										<div class="min-w-0 flex-1">
											<p class="truncate text-sm font-medium">
												{upload.file.name}
											</p>
											<p class="text-xs">
												{formatFileSize(upload.file.size)}
												{#if upload.media && upload.media.width && upload.media.height}
													• {upload.media.width}×{upload.media.height}
												{/if}
											</p>
										</div>

										<div class="flex-shrink-0">
											{#if upload.status === 'success'}
												<CheckCircle class="h-5 w-5 text-green-500" />
											{:else if upload.status === 'error'}
												<X class="h-5 w-5 text-red-500" />
											{:else}
												<Loader2 class="h-5 w-5 animate-spin text-blue-500" />
											{/if}
										</div>
									</div>

									{#if upload.status === 'uploading' || upload.status === 'processing'}
										<div class="space-y-1">
											<div class="h-2 w-full overflow-hidden rounded-full bg-gray-200">
												<div
													class="h-full transition-all duration-300"
													class:bg-blue-500={upload.status === 'uploading'}
													class:bg-yellow-500={upload.status === 'processing'}
													style="width: {upload.progress}%"
												></div>
											</div>
											<p class="text-xs">
												{#if upload.status === 'uploading'}
													Uploading... {upload.progress}%
												{:else if upload.status === 'processing'}
													{getProcessingMessage(upload.file)}
												{/if}
											</p>
										</div>
									{/if}

									{#if upload.status === 'success' && upload.media}
										<div class="flex items-center gap-2 text-xs font-medium text-green-600">
											<CheckCircle class="h-4 w-4" />
											{#if upload.file.type.startsWith('image/')}
												Converted to WebP
											{:else if upload.file.type.startsWith('video/')}
												Video optimized
											{:else}
												Upload complete
											{/if}
											• {formatFileSize(upload.media.size_bytes)}
										</div>
									{/if}

									{#if upload.status === 'error'}
										<div class="flex items-center gap-2 text-xs text-red-600">
											<X class="h-4 w-4" />
											{upload.error}
										</div>
									{/if}
								</div>
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>

		<!-- Stats -->
		{#if media.length > 0}
			<div class="mb-8 grid grid-cols-1 gap-4 md:grid-cols-4">
				<div class="rounded-lg bg-surface p-4 shadow-sm">
					<div class="text-sm">Total Files</div>
					<div class="text-2xl font-bold">{media.length}</div>
				</div>
				<div class="rounded-lg bg-surface p-4 shadow-sm">
					<div class="text-sm">Images</div>
					<div class="text-2xl font-bold text-blue-600">
						{media.filter(isImage).length}
					</div>
				</div>
				<div class="rounded-lg bg-surface p-4 shadow-sm">
					<div class="text-sm">Videos</div>
					<div class="text-2xl font-bold text-purple-600">
						{media.filter(isVideo).length}
					</div>
				</div>
				<div class="rounded-lg bg-surface p-4 shadow-sm">
					<div class="text-sm">Total Size</div>
					<div class="text-2xl font-bold">
						{formatFileSize(media.reduce((acc, m) => acc + m.size_bytes, 0))}
					</div>
				</div>
			</div>
		{/if}

		<!-- Media Grid -->
		<div class="rounded-lg bg-surface p-6 shadow-sm">
			<h2 class="mb-4 text-xl font-semibold">Media Library</h2>

			{#if loading}
				<div class="flex items-center justify-center py-12">
					<Loader2 class="h-8 w-8 animate-spin text-gray-400" />
				</div>
			{:else if media.length === 0}
				<div class="py-12 text-center">
					<ImageUp class="mx-auto mb-4 h-16 w-16 text-gray-400" />
					<p class=" ">No media uploaded yet. Upload your first file!</p>
				</div>
			{:else}
				<div class="grid grid-cols-2 gap-4 md:grid-cols-3 lg:grid-cols-4">
					{#each media as m (m.id)}
						<button
							onclick={() => (selectedMedia = m)}
							class="group relative aspect-square overflow-hidden rounded-lg bg-gray-100 transition-all hover:ring-2 hover:ring-blue-500"
						>
							{#if m.thumbnail_url}
								<img
									src={m.thumbnail_url}
									alt={m.original_filename}
									class="h-full w-full object-cover"
								/>
							{:else if isImage(m)}
								<img src={m.url} alt={m.original_filename} class="h-full w-full object-cover" />
							{:else if isVideo(m)}
								<div class="flex h-full w-full items-center justify-center bg-gray-900">
									<Video class="h-16 w-16 text-white opacity-50" />
								</div>
							{:else}
								<div class="flex h-full w-full items-center justify-center">
									<ImageUp class="h-16 w-16 text-gray-400" />
								</div>
							{/if}

							<div
								class="group-hover:bg-opacity-40 absolute inset-0 flex items-center justify-center bg-black/10 transition-all"
							>
								<div class="opacity-0 transition-opacity group-hover:opacity-100">
									<p class="max-w-full truncate px-2 text-center text-sm font-medium text-white">
										{m.original_filename}
									</p>
								</div>
							</div>

							<div class="absolute top-2 right-2 flex gap-1">
								{#if m.mime_type === 'image/webp'}
									<span class="rounded-full bg-green-500 px-2 py-1 text-xs text-white"> WebP </span>
								{/if}
								{#if isVideo(m)}
									<span class="rounded-full bg-purple-500 px-2 py-1 text-xs text-white">
										Video
									</span>
								{/if}
							</div>
						</button>
					{/each}
				</div>
			{/if}
		</div>
	</div>
</div>

<!-- Media Detail Modal -->
{#if selectedMedia}
	<div
		class="bg-opacity-75 fixed inset-0 z-50 flex items-center justify-center bg-black p-4"
		onclick={() => (selectedMedia = null)}
	>
		<div
			class="max-h-[90vh] w-full max-w-4xl overflow-hidden rounded-lg bg-surface"
			onclick={(e) => e.stopPropagation()}
		>
			<div class="flex items-center justify-between border-b p-4">
				<h3 class="mr-4 flex-1 truncate text-lg font-semibold">
					{selectedMedia.original_filename}
				</h3>
				<button
					onclick={() => (selectedMedia = null)}
					class="flex-shrink-0 text-gray-400 hover:text-gray-600"
				>
					<X class="h-6 w-6" />
				</button>
			</div>

			<div class="max-h-[calc(90vh-140px)] overflow-y-auto p-6">
				<div class="mb-6">
					{#if isImage(selectedMedia)}
						<img
							src={selectedMedia.url}
							alt={selectedMedia.original_filename}
							class="w-full rounded-lg"
						/>
					{:else if isVideo(selectedMedia)}
						<video src={selectedMedia.url} controls class="w-full rounded-lg">
							<track kind="captions" />
						</video>
					{/if}
				</div>

				<div class="mb-6 grid grid-cols-2 gap-4">
					<div>
						<div class="text-sm">Filename</div>
						<div class="font-medium break-all">{selectedMedia.filename}</div>
					</div>
					<div>
						<div class="text-sm">Original Name</div>
						<div class="font-medium break-all">{selectedMedia.original_filename}</div>
					</div>
					<div>
						<div class="text-sm">Type</div>
						<div class="font-medium">{selectedMedia.mime_type}</div>
					</div>
					<div>
						<div class="text-sm">Size</div>
						<div class="font-medium">{formatFileSize(selectedMedia.size_bytes)}</div>
					</div>
					{#if selectedMedia.width && selectedMedia.height}
						<div>
							<div class="text-sm">Dimensions</div>
							<div class="font-medium">
								{selectedMedia.width} × {selectedMedia.height}
							</div>
						</div>
					{/if}
				</div>

				<div class="mb-6 space-y-3">
					<div>
						<div class="mb-1 text-sm">Media URL</div>
						<div class="rounded bg-gray-50 p-2 font-mono text-sm break-all text-gray-700">
							{selectedMedia.url}
						</div>
					</div>
					{#if selectedMedia.thumbnail_url}
						<div>
							<div class="mb-1 text-sm">Thumbnail URL</div>
							<div class="rounded bg-gray-50 p-2 font-mono text-sm break-all text-gray-700">
								{selectedMedia.thumbnail_url}
							</div>
						</div>
					{/if}
				</div>

				<div class="flex gap-2">
					<a
						href={selectedMedia.url}
						download={selectedMedia.filename}
						class="flex flex-1 items-center justify-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-white transition-colors hover:bg-blue-700"
					>
						<Download class="h-4 w-4" />
						Download
					</a>
					<button
						onclick={() => deleteMedia(selectedMedia!)}
						class="flex flex-1 items-center justify-center gap-2 rounded-lg bg-red-600 px-4 py-2 text-white transition-colors hover:bg-red-700"
					>
						<Trash2 class="h-4 w-4" />
						Delete
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}
