<script lang="ts">
import { PUBLIC_API_URL } from "$env/static/public";
import {
	ImageUp,
	Video,
	Download,
	Trash2,
	X,
	CheckCircle,
	Loader2,
} from "lucide-svelte";
import { onMount } from "svelte";

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
	status: "uploading" | "processing" | "success" | "error";
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
		const response = await fetch(
			`${PUBLIC_API_URL}/api/v1/media?page=1&limit=100`,
			{
				credentials: "include",
			},
		);

		if (response.ok) {
			const data = await response.json();
			media = data.data || [];
		}
	} catch (error) {
		console.error("Failed to load media:", error);
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
			status: "uploading",
		};
		uploads = [...uploads, uploadState];
		uploadFile(uploadState);
	});
}

async function uploadFile(uploadState: UploadState) {
	const formData = new FormData();
	formData.append("file", uploadState.file);

	try {
		const xhr = new XMLHttpRequest();

		xhr.upload.addEventListener("progress", (e) => {
			if (e.lengthComputable) {
				uploadState.progress = Math.round((e.loaded / e.total) * 100);

				if (uploadState.progress === 100) {
					uploadState.status = "processing";
				}
			}
		});

		const uploadedMedia = await new Promise<Media>((resolve, reject) => {
			xhr.addEventListener("load", () => {
				if (xhr.status === 201) {
					resolve(JSON.parse(xhr.responseText));
				} else {
					const error = JSON.parse(xhr.responseText);
					reject(new Error(error.error || "Upload failed"));
				}
			});

			xhr.addEventListener("error", () => {
				reject(new Error("Network error during upload"));
			});

			xhr.open("POST", `${PUBLIC_API_URL}/api/v1/media/upload`);
			xhr.withCredentials = true;
			xhr.send(formData);
		});

		uploadState.status = "success";
		uploadState.media = uploadedMedia;
		media = [uploadedMedia, ...media];
	} catch (error) {
		uploadState.status = "error";
		uploadState.error =
			error instanceof Error ? error.message : "Upload failed";
	}
}

async function deleteMedia(m: Media) {
	if (!confirm(`Delete ${m.original_filename}?`)) return;

	try {
		const response = await fetch(`${PUBLIC_API_URL}/api/v1/media/${m.id}`, {
			method: "DELETE",
			credentials: "include",
		});

		if (response.ok) {
			media = media.filter((item) => item.id !== m.id);
			if (selectedMedia?.id === m.id) {
				selectedMedia = null;
			}
		}
	} catch (error) {
		console.error("Delete failed:", error);
		alert("Failed to delete media");
	}
}

function formatFileSize(bytes: number): string {
	if (bytes < 1024) return bytes + " B";
	if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB";
	return (bytes / (1024 * 1024)).toFixed(1) + " MB";
}

function isImage(m: Media): boolean {
	return m.mime_type.startsWith("image/");
}

function isVideo(m: Media): boolean {
	return m.mime_type.startsWith("video/");
}

function getProcessingMessage(file: File): string {
	if (file.type.startsWith("image/")) {
		return "Converting to WebP and generating thumbnail...";
	} else if (file.type.startsWith("video/")) {
		return "Optimizing video and extracting thumbnail...";
	}
	return "Processing file...";
}

function clearUploads() {
	uploads = [];
}
</script>

<div class="min-h-screen bg-gray-50 py-8">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
		<!-- Header -->
		<div class="mb-8">
			<h1 class="text-3xl font-bold text-gray-900">Media Library</h1>
			<p class="mt-2 text-gray-600">
				Upload and manage images and videos with automatic processing
			</p>
		</div>

		<!-- Upload Area -->
		<div class="bg-white rounded-lg shadow-sm p-6 mb-8">
			<h2 class="text-xl font-semibold text-gray-900 mb-4">Upload Media</h2>
			
			<div
				class="relative border-2 border-dashed rounded-lg p-8 text-center transition-colors cursor-pointer"
				class:border-blue-500={dragOver}
				class:bg-blue-50={dragOver}
				class:border-gray-300={!dragOver}
				ondragover={(e) => { e.preventDefault(); dragOver = true; }}
				ondragleave={() => { dragOver = false; }}
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
				
				<ImageUp class="mx-auto h-12 w-12 text-gray-400 mb-4" />
				
				<div class="space-y-2">
					<p class="text-lg font-medium text-gray-900">
						Drop files here or click to upload
					</p>
					<p class="text-sm text-gray-500">
						Images will be converted to WebP • Videos will be optimized
					</p>
					<p class="text-xs text-gray-400">
						Images: 50MB max • Videos: 200MB max
					</p>
				</div>
			</div>

			<!-- Upload Progress -->
			{#if uploads.length > 0}
				<div class="mt-6 space-y-3">
					<div class="flex items-center justify-between">
						<h3 class="text-sm font-medium text-gray-700">
							Uploading {uploads.length} {uploads.length === 1 ? 'file' : 'files'}
						</h3>
						<button
							onclick={clearUploads}
							class="text-sm text-gray-500 hover:text-gray-700"
						>
							Clear all
						</button>
					</div>

					{#each uploads as upload (upload.file.name)}
						<div class="border rounded-lg p-4 bg-white">
							<div class="flex items-start gap-4">
								<div class="flex-shrink-0">
									{#if upload.status === 'success' && upload.media?.thumbnail_url}
										<img 
											src={upload.media.thumbnail_url} 
											alt={upload.file.name}
											class="w-16 h-16 object-cover rounded"
										/>
									{:else}
										<div class="w-16 h-16 bg-gray-100 rounded flex items-center justify-center">
											{#if upload.file.type.startsWith('image/')}
												<ImageUp class="w-8 h-8 text-gray-400" />
											{:else}
												<Video class="w-8 h-8 text-gray-400" />
											{/if}
										</div>
									{/if}
								</div>

								<div class="flex-1 min-w-0">
									<div class="flex items-start justify-between gap-2 mb-2">
										<div class="flex-1 min-w-0">
											<p class="text-sm font-medium text-gray-900 truncate">
												{upload.file.name}
											</p>
											<p class="text-xs text-gray-500">
												{formatFileSize(upload.file.size)}
												{#if upload.media && upload.media.width && upload.media.height}
													• {upload.media.width}×{upload.media.height}
												{/if}
											</p>
										</div>

										<div class="flex-shrink-0">
											{#if upload.status === 'success'}
												<CheckCircle class="w-5 h-5 text-green-500" />
											{:else if upload.status === 'error'}
												<X class="w-5 h-5 text-red-500" />
											{:else}
												<Loader2 class="w-5 h-5 text-blue-500 animate-spin" />
											{/if}
										</div>
									</div>

									{#if upload.status === 'uploading' || upload.status === 'processing'}
										<div class="space-y-1">
											<div class="w-full bg-gray-200 rounded-full h-2 overflow-hidden">
												<div 
													class="h-full transition-all duration-300"
													class:bg-blue-500={upload.status === 'uploading'}
													class:bg-yellow-500={upload.status === 'processing'}
													style="width: {upload.progress}%"
												></div>
											</div>
											<p class="text-xs text-gray-500">
												{#if upload.status === 'uploading'}
													Uploading... {upload.progress}%
												{:else if upload.status === 'processing'}
													{getProcessingMessage(upload.file)}
												{/if}
											</p>
										</div>
									{/if}

									{#if upload.status === 'success' && upload.media}
										<div class="flex items-center gap-2 text-xs text-green-600 font-medium">
											<CheckCircle class="w-4 h-4" />
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
											<X class="w-4 h-4" />
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
			<div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
				<div class="bg-white rounded-lg shadow-sm p-4">
					<div class="text-sm text-gray-500">Total Files</div>
					<div class="text-2xl font-bold text-gray-900">{media.length}</div>
				</div>
				<div class="bg-white rounded-lg shadow-sm p-4">
					<div class="text-sm text-gray-500">Images</div>
					<div class="text-2xl font-bold text-blue-600">
						{media.filter(isImage).length}
					</div>
				</div>
				<div class="bg-white rounded-lg shadow-sm p-4">
					<div class="text-sm text-gray-500">Videos</div>
					<div class="text-2xl font-bold text-purple-600">
						{media.filter(isVideo).length}
					</div>
				</div>
				<div class="bg-white rounded-lg shadow-sm p-4">
					<div class="text-sm text-gray-500">Total Size</div>
					<div class="text-2xl font-bold text-gray-900">
						{formatFileSize(media.reduce((acc, m) => acc + m.size_bytes, 0))}
					</div>
				</div>
			</div>
		{/if}

		<!-- Media Grid -->
		<div class="bg-white rounded-lg shadow-sm p-6">
			<h2 class="text-xl font-semibold text-gray-900 mb-4">Media Library</h2>
			
			{#if loading}
				<div class="flex items-center justify-center py-12">
					<Loader2 class="w-8 h-8 text-gray-400 animate-spin" />
				</div>
			{:else if media.length === 0}
				<div class="text-center py-12">
					<ImageUp class="mx-auto h-16 w-16 text-gray-400 mb-4" />
					<p class="text-gray-500">No media uploaded yet. Upload your first file!</p>
				</div>
			{:else}
				<div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
					{#each media as m (m.id)}
						<button
							onclick={() => selectedMedia = m}
							class="group relative aspect-square rounded-lg overflow-hidden bg-gray-100 hover:ring-2 hover:ring-blue-500 transition-all"
						>
							{#if m.thumbnail_url}
								<img 
									src={m.thumbnail_url} 
									alt={m.original_filename}
									class="w-full h-full object-cover"
								/>
							{:else if isImage(m)}
								<img 
									src={m.url} 
									alt={m.original_filename}
									class="w-full h-full object-cover"
								/>
							{:else if isVideo(m)}
								<div class="w-full h-full flex items-center justify-center bg-gray-900">
									<Video class="w-16 h-16 text-white opacity-50" />
								</div>
							{:else}
								<div class="w-full h-full flex items-center justify-center">
									<ImageUp class="w-16 h-16 text-gray-400" />
								</div>
							{/if}

							<div class="absolute inset-0 bg-black/10  group-hover:bg-opacity-40 transition-all flex items-center justify-center">
								<div class="opacity-0 group-hover:opacity-100 transition-opacity">
									<p class="text-white text-sm font-medium px-2 text-center truncate max-w-full">
										{m.original_filename}
									</p>
								</div>
							</div>

							<div class="absolute top-2 right-2 flex gap-1">
								{#if m.mime_type === 'image/webp'}
									<span class="bg-green-500 text-white text-xs px-2 py-1 rounded-full">
										WebP
									</span>
								{/if}
								{#if isVideo(m)}
									<span class="bg-purple-500 text-white text-xs px-2 py-1 rounded-full">
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
		class="fixed inset-0 bg-black bg-opacity-75 z-50 flex items-center justify-center p-4"
		onclick={() => selectedMedia = null}
	>
		<div 
			class="bg-white rounded-lg max-w-4xl w-full max-h-[90vh] overflow-hidden"
			onclick={(e) => e.stopPropagation()}
		>
			<div class="flex items-center justify-between p-4 border-b">
				<h3 class="text-lg font-semibold text-gray-900 truncate flex-1 mr-4">
					{selectedMedia.original_filename}
				</h3>
				<button 
					onclick={() => selectedMedia = null}
					class="text-gray-400 hover:text-gray-600 flex-shrink-0"
				>
					<X class="w-6 h-6" />
				</button>
			</div>

			<div class="p-6 overflow-y-auto max-h-[calc(90vh-140px)]">
				<div class="mb-6">
					{#if isImage(selectedMedia)}
						<img 
							src={selectedMedia.url} 
							alt={selectedMedia.original_filename}
							class="w-full rounded-lg"
						/>
					{:else if isVideo(selectedMedia)}
						<video 
							src={selectedMedia.url} 
							controls
							class="w-full rounded-lg"
						>
							<track kind="captions" />
						</video>
					{/if}
				</div>

				<div class="grid grid-cols-2 gap-4 mb-6">
					<div>
						<div class="text-sm text-gray-500">Filename</div>
						<div class="font-medium text-gray-900 break-all">{selectedMedia.filename}</div>
					</div>
					<div>
						<div class="text-sm text-gray-500">Original Name</div>
						<div class="font-medium text-gray-900 break-all">{selectedMedia.original_filename}</div>
					</div>
					<div>
						<div class="text-sm text-gray-500">Type</div>
						<div class="font-medium text-gray-900">{selectedMedia.mime_type}</div>
					</div>
					<div>
						<div class="text-sm text-gray-500">Size</div>
						<div class="font-medium text-gray-900">{formatFileSize(selectedMedia.size_bytes)}</div>
					</div>
					{#if selectedMedia.width && selectedMedia.height}
						<div>
							<div class="text-sm text-gray-500">Dimensions</div>
							<div class="font-medium text-gray-900">
								{selectedMedia.width} × {selectedMedia.height}
							</div>
						</div>
					{/if}
				</div>

				<div class="space-y-3 mb-6">
					<div>
						<div class="text-sm text-gray-500 mb-1">Media URL</div>
						<div class="bg-gray-50 p-2 rounded text-sm font-mono text-gray-700 break-all">
							{selectedMedia.url}
						</div>
					</div>
					{#if selectedMedia.thumbnail_url}
						<div>
							<div class="text-sm text-gray-500 mb-1">Thumbnail URL</div>
							<div class="bg-gray-50 p-2 rounded text-sm font-mono text-gray-700 break-all">
								{selectedMedia.thumbnail_url}
							</div>
						</div>
					{/if}
				</div>

				<div class="flex gap-2">
					<a
						href={selectedMedia.url}
						download={selectedMedia.filename}
						class="flex-1 flex items-center justify-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
					>
						<Download class="w-4 h-4" />
						Download
					</a>
					<button
						onclick={() => deleteMedia(selectedMedia!)}
						class="flex-1 flex items-center justify-center gap-2 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
					>
						<Trash2 class="w-4 h-4" />
						Delete
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}
