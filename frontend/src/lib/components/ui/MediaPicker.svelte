<script lang="ts">
import { PUBLIC_API_URL } from "$env/static/public";
import { mediaApi, type Media, getMediaUrl } from "$api/media";
import { onMount } from "svelte";
import { ImageUp, CheckCircle, AlertCircle, Loader2 } from "lucide-svelte";

interface Props {
	value?: string;
	label?: string;
	required?: boolean;
	error?: string;
	onchange?: (mediaId: string | null) => void;
}

let {
	value = $bindable(""),
	label,
	required = false,
	error,
	onchange,
}: Props = $props();

let showModal = $state(false);
let media = $state<Media[]>([]);
let selectedMedia = $state<Media | null>(null);
let loading = $state(false);
let searchQuery = $state("");
let uploadFile = $state<File | null>(null);
let uploading = $state(false);
let uploadProgress = $state(0);
let uploadStatus = $state<"uploading" | "processing" | "success" | "error">(
	"uploading",
);
let uploadError = $state<string>("");

// View mode: 'grid' or 'list'
let viewMode = $state<"grid" | "list">("grid");

// Pagination
let currentPage = $state(1);
let totalPages = $state(1);
let limit = $state(20);

// ✅ Accept images, videos, and GIFs
const ACCEPTED_FILE_TYPES =
	"image/*,video/*,video/mp4,video/quicktime,.mp4,.mov,.avi,.gif";

onMount(async () => {
	if (value) {
		await loadSelectedMedia();
	}
});

// ✅ Watch for value changes
$effect(() => {
	if (value && (!selectedMedia || selectedMedia.id !== value)) {
		loadSelectedMedia();
	} else if (!value) {
		selectedMedia = null;
	}
});

async function loadSelectedMedia() {
	try {
		const response = await fetch(`${PUBLIC_API_URL}/api/v1/media/${value}`, {
			credentials: "include",
		});
		if (response.ok) {
			selectedMedia = await response.json();
		}
	} catch (err) {
		console.error("Failed to load selected media:", err);
	}
}

async function openPicker() {
	showModal = true;
	currentPage = 1;
	await loadMedia();
}

async function loadMedia() {
	loading = true;
	try {
		const response = await mediaApi.listMedia(currentPage, limit);
		media = response.data || [];

		if (response.pagination) {
			totalPages = response.pagination.total_pages || 1;
		}
	} catch (err) {
		console.error("Failed to load media:", err);
	} finally {
		loading = false;
	}
}

async function changePage(newPage: number) {
	if (newPage < 1 || newPage > totalPages) return;
	currentPage = newPage;
	await loadMedia();
}

async function handleFileSelect(e: Event) {
	const input = e.target as HTMLInputElement;
	if (input.files && input.files[0]) {
		uploadFile = input.files[0];
		await handleUpload();
	}
}

async function handleUpload() {
	if (!uploadFile) return;

	const maxSize = uploadFile.type.startsWith("video/") ? 200 : 50; // MB
	const fileSizeMB = uploadFile.size / (1024 * 1024);

	if (fileSizeMB > maxSize) {
		alert(`File size (${fileSizeMB.toFixed(1)}MB) exceeds ${maxSize}MB limit`);
		return;
	}

	uploading = true;
	uploadProgress = 0;
	uploadStatus = "uploading";
	uploadError = "";

	try {
		const formData = new FormData();
		formData.append("file", uploadFile);

		const xhr = new XMLHttpRequest();

		xhr.upload.addEventListener("progress", (e) => {
			if (e.lengthComputable) {
				uploadProgress = Math.round((e.loaded / e.total) * 100);

				// When upload completes, show processing status
				if (uploadProgress === 100) {
					uploadStatus = "processing";
				}
			}
		});

		const uploaded = await new Promise<Media>((resolve, reject) => {
			xhr.addEventListener("load", () => {
				if (xhr.status === 201) {
					uploadStatus = "success";
					resolve(JSON.parse(xhr.responseText));
				} else {
					uploadStatus = "error";
					const error = JSON.parse(xhr.responseText);
					reject(new Error(error.error || "Upload failed"));
				}
			});

			xhr.addEventListener("error", () => {
				uploadStatus = "error";
				reject(new Error("Network error during upload"));
			});

			xhr.open("POST", `${PUBLIC_API_URL}/api/v1/media/upload`);
			xhr.withCredentials = true;
			xhr.send(formData);
		});

		media = [uploaded, ...media];
		selectMedia(uploaded);
		uploadFile = null;
	} catch (err) {
		console.error("Upload failed:", err);
		uploadError = err instanceof Error ? err.message : "Upload failed";
		uploadStatus = "error";
	} finally {
		setTimeout(() => {
			uploading = false;
			uploadProgress = 0;
			uploadStatus = "uploading";
			uploadError = "";
		}, 2000);
	}
}

function selectMedia(m: Media) {
	selectedMedia = m;
	value = m.id;
	onchange?.(m.id);
	showModal = false;
}

function clearSelection() {
	selectedMedia = null;
	value = "";
	onchange?.(null);
}

function formatFileSize(bytes: number): string {
	if (bytes < 1024) return bytes + " B";
	if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB";
	return (bytes / (1024 * 1024)).toFixed(1) + " MB";
}

function formatDate(dateString: string): string {
	return new Date(dateString).toLocaleDateString("en-US", {
		year: "numeric",
		month: "short",
		day: "numeric",
	});
}

function isImage(m: Media): boolean {
	return m.mime_type.startsWith("image/");
}

function isVideo(m: Media): boolean {
	return m.mime_type.startsWith("video/");
}

function getProcessingMessage(): string {
	if (!uploadFile) return "Processing...";

	if (isImage({ mime_type: uploadFile.type } as Media)) {
		return "Converting to WebP and generating thumbnail...";
	} else if (isVideo({ mime_type: uploadFile.type } as Media)) {
		return "Optimizing video and extracting thumbnail...";
	}
	return "Processing file...";
}

const filteredMedia = $derived(
	media.filter((m) => {
		if (!searchQuery.trim()) return true;
		return m.original_filename
			.toLowerCase()
			.includes(searchQuery.toLowerCase());
	}),
);
</script>

<div class="space-y-2">
	{#if label}
		<label class="block text-sm font-medium text-gray-700">
			{label}
			{#if required}
				<span class="text-red-500">*</span>
			{/if}
		</label>
	{/if}

	{#if selectedMedia}
		<!-- Selected Media Preview -->
		<div class="relative group border-2 border-gray-200 rounded-lg p-4 bg-gray-50">
			<button
				type="button"
				onclick={clearSelection}
				class="absolute top-2 right-2 p-1 bg-red-500 text-white rounded-full opacity-0 group-hover:opacity-100 transition-opacity"
				title="Remove"
			>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
				</svg>
			</button>

			<div class="flex items-center gap-4">
				<!-- Use thumbnail if available, fallback to main URL -->
				<img
					src={selectedMedia.thumbnail_url || getMediaUrl(selectedMedia)}
					alt={selectedMedia.original_filename}
					class="w-20 h-20 object-cover rounded"
				/>
				<div class="flex-1 min-w-0">
					<p class="text-sm font-medium text-gray-900 truncate">
						{selectedMedia.original_filename}
					</p>
					<p class="text-xs text-gray-500">
						{formatFileSize(selectedMedia.size_bytes)}
						{#if selectedMedia.width && selectedMedia.height}
							• {selectedMedia.width}×{selectedMedia.height}
						{/if}
					</p>
					<!-- Show if it's WebP processed -->
					{#if selectedMedia.mime_type === 'image/webp'}
						<span class="inline-flex items-center gap-1 mt-1 text-xs text-green-600 font-medium">
							<CheckCircle class="w-3 h-3" />
							WebP Optimized
						</span>
					{/if}
				</div>
				<button
					type="button"
					onclick={openPicker}
					class="px-3 py-1.5 text-sm font-medium text-blue-600 hover:text-blue-700"
				>
					Change
				</button>
			</div>
		</div>
	{:else}
		<!-- Select Media Button -->
		<button
			type="button"
			onclick={openPicker}
			class="w-full border-2 border-dashed border-gray-300 rounded-lg p-6 text-center hover:border-gray-400 transition-colors"
		>
			<ImageUp class="mx-auto h-12 w-12 text-gray-400 mb-2" />
			<p class="text-sm text-gray-600">Click to select or upload media</p>
			<p class="text-xs text-gray-400 mt-1">Images will be converted to WebP • Videos will be optimized</p>
		</button>
	{/if}

	{#if error}
		<p class="text-sm text-red-600">{error}</p>
	{/if}
</div>

<!-- Media Picker Modal -->
{#if showModal}
	<div
		class="fixed inset-0 z-50 overflow-y-auto"
		onclick={() => (showModal = false)}
	>
		<div class="flex min-h-screen items-center justify-center p-4">
			<div
				class="fixed inset-0 bg-black/50 transition-opacity"
			></div>

			<div
				class="relative bg-white rounded-lg shadow-xl max-w-6xl w-full max-h-[90vh] overflow-hidden"
				onclick={(e) => e.stopPropagation()}
			>
				<!-- Header -->
				<div class="border-b border-gray-200 px-6 py-4">
					<div class="flex items-center justify-between">
						<div class="flex-1">
							<h2 class="text-xl font-semibold text-gray-900">Select Media</h2>
							<p class="text-sm text-gray-500 mt-1">
								Upload new or select existing media
							</p>
						</div>

						<!-- Upload Status -->
						{#if uploading}
							<div class="flex items-center gap-3 px-4 py-2 bg-blue-50 rounded-lg mr-4">
								{#if uploadStatus === 'uploading'}
									<Loader2 class="w-5 h-5 text-blue-600 animate-spin" />
									<div>
										<p class="text-sm font-medium text-blue-900">Uploading...</p>
										<p class="text-xs text-blue-600">{uploadProgress}%</p>
									</div>
								{:else if uploadStatus === 'processing'}
									<Loader2 class="w-5 h-5 text-yellow-600 animate-spin" />
									<div>
										<p class="text-sm font-medium text-yellow-900">Processing</p>
										<p class="text-xs text-yellow-600">{getProcessingMessage()}</p>
									</div>
								{:else if uploadStatus === 'success'}
									<CheckCircle class="w-5 h-5 text-green-600" />
									<p class="text-sm font-medium text-green-900">Success!</p>
								{:else if uploadStatus === 'error'}
									<AlertCircle class="w-5 h-5 text-red-600" />
									<div>
										<p class="text-sm font-medium text-red-900">Error</p>
										<p class="text-xs text-red-600">{uploadError}</p>
									</div>
								{/if}
							</div>
						{/if}

						<button
							onclick={() => (showModal = false)}
							class="text-gray-400 hover:text-gray-600"
						>
							<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
							</svg>
						</button>
					</div>

					<!-- Search and Actions -->
					<div class="flex gap-3 mt-4">
						<input
							type="text"
							bind:value={searchQuery}
							placeholder="Search media..."
							class="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
						/>
						
						<!-- View Toggle -->
						<div class="flex border border-gray-300 rounded-lg overflow-hidden">
							<button
								onclick={() => (viewMode = "grid")}
								class="px-3 py-2 {viewMode === 'grid' ? 'bg-blue-600 text-white' : 'bg-white text-gray-700 hover:bg-gray-50'}"
							>
								<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
								</svg>
							</button>
							<button
								onclick={() => (viewMode = "list")}
								class="px-3 py-2 {viewMode === 'list' ? 'bg-blue-600 text-white' : 'bg-white text-gray-700 hover:bg-gray-50'}"
							>
								<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
								</svg>
							</button>
						</div>

						<!-- Upload Button -->
						<label class="px-4 py-2 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700 cursor-pointer flex items-center gap-2">
							{#if uploading}
								<Loader2 class="w-5 h-5 animate-spin" />
								Uploading...
							{:else}
								<ImageUp class="w-5 h-5" />
								Upload
							{/if}
							<input
								type="file"
								accept={ACCEPTED_FILE_TYPES}
								onchange={handleFileSelect}
								class="hidden"
								disabled={uploading}
							/>
						</label>
					</div>
				</div>

				<!-- Content -->
				<div class="px-6 py-6 max-h-[600px] overflow-y-auto">
					{#if loading}
						<div class="flex items-center justify-center py-12">
							<Loader2 class="w-8 h-8 text-gray-400 animate-spin" />
						</div>
					{:else if filteredMedia.length === 0}
						<div class="text-center py-12">
							<ImageUp class="mx-auto h-12 w-12 text-gray-400 mb-4" />
							<p class="text-gray-500">No media found. Upload your first file!</p>
						</div>
					{:else if viewMode === "grid"}
						<div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
							{#each filteredMedia as m (m.id)}
								<button
									onclick={() => selectMedia(m)}
									class="group relative aspect-square rounded-lg overflow-hidden bg-gray-100 hover:ring-2 hover:ring-blue-500 transition-all"
								>
									<!-- Use thumbnail if available -->
									<img
										src={m.thumbnail_url || getMediaUrl(m)}
										alt={m.original_filename}
										class="w-full h-full object-cover"
									/>
									
									<!-- Badges -->
									<div class="absolute top-2 right-2 flex gap-1">
										{#if m.mime_type === 'image/webp'}
											<span class="px-2 py-1 bg-green-500 text-white text-xs rounded-full font-medium">
												WebP
											</span>
										{/if}
										{#if isVideo(m)}
											<span class="px-2 py-1 bg-purple-500 text-white text-xs rounded-full font-medium">
												Video
											</span>
										{/if}
									</div>

									<!-- Info Overlay -->
									<div class="absolute inset-0 bg-black/15  group-hover:bg-opacity-50 transition-opacity flex items-end">
										<div class="w-full p-3 text-white opacity-0 group-hover:opacity-100 transition-opacity">
											<p class="text-sm font-medium truncate">{m.original_filename}</p>
											<p class="text-xs">{formatFileSize(m.size_bytes)}</p>
										</div>
									</div>

									<!-- Selected Check -->
									{#if selectedMedia?.id === m.id}
										<div class="absolute inset-0 bg-blue-500 bg-opacity-20 flex items-center justify-center">
											<div class="bg-blue-500 rounded-full p-2">
												<CheckCircle class="w-6 h-6 text-white" />
											</div>
										</div>
									{/if}
								</button>
							{/each}
						</div>
					{:else}
						<!-- List View -->
						<div class="bg-white rounded-lg shadow overflow-hidden">
							<table class="min-w-full divide-y divide-gray-200">
								<thead class="bg-gray-50">
									<tr>
										<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Preview</th>
										<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Filename</th>
										<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Size</th>
										<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Type</th>
										<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Action</th>
									</tr>
								</thead>
								<tbody class="divide-y divide-gray-200">
									{#each filteredMedia as m (m.id)}
										<tr class="hover:bg-gray-50">
											<td class="px-6 py-4">
												<img
													src={m.thumbnail_url || getMediaUrl(m)}
													alt={m.original_filename}
													class="h-12 w-12 object-cover rounded"
												/>
											</td>
											<td class="px-6 py-4">
												<div class="text-sm font-medium text-gray-900 truncate max-w-xs">
													{m.original_filename}
												</div>
											</td>
											<td class="px-6 py-4 text-sm text-gray-500">
												{formatFileSize(m.size_bytes)}
											</td>
											<td class="px-6 py-4">
												{#if m.mime_type === 'image/webp'}
													<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
														WebP
													</span>
												{:else if isVideo(m)}
													<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-800">
														Video
													</span>
												{:else}
													<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
														{m.mime_type.split('/')[1].toUpperCase()}
													</span>
												{/if}
											</td>
											<td class="px-6 py-4">
												<button
													onclick={() => selectMedia(m)}
													class="text-blue-600 hover:text-blue-700 font-medium text-sm"
												>
													Select
												</button>
											</td>
										</tr>
									{/each}
								</tbody>
							</table>
						</div>
					{/if}

					<!-- Pagination -->
					{#if totalPages > 1}
						<div class="flex items-center justify-center gap-2 mt-6">
							<button
								onclick={() => changePage(currentPage - 1)}
								disabled={currentPage === 1}
								class="px-3 py-1 border rounded disabled:opacity-50"
							>
								Previous
							</button>
							<span class="text-sm text-gray-600">
								Page {currentPage} of {totalPages}
							</span>
							<button
								onclick={() => changePage(currentPage + 1)}
								disabled={currentPage === totalPages}
								class="px-3 py-1 border rounded disabled:opacity-50"
							>
								Next
							</button>
						</div>
					{/if}
				</div>
			</div>
		</div>
	</div>
{/if}
