<script lang="ts">
import { PUBLIC_API_URL } from "$env/static/public";
import { mediaApi, type Media, getMediaUrl } from "$api/media";
import { onMount } from "svelte";

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

// View mode: 'grid' or 'list'
let viewMode = $state<"grid" | "list">("grid");

// Pagination
let currentPage = $state(1);
let totalPages = $state(1);
let limit = $state(20);

onMount(async () => {
	if (value) {
		await loadSelectedMedia();
	}
});

async function loadSelectedMedia() {
	try {
		const response = await fetch(
			`${PUBLIC_API_URL || "localhost:8080"}/api/v1/media/${value}`,
			{
				credentials: "include",
			},
		);
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

		// Debug: Log the first media item to see what we're getting
		if (media.length > 0) {
			console.log("MediaPicker - First media item:", media[0]);
			console.log("MediaPicker - URL:", media[0].url);
			console.log("MediaPicker - Storage path:", media[0].storage_path);
			console.log("MediaPicker - Generated URL:", getMediaUrl(media[0]));
		}

		// Calculate total pages from response if available
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

	uploading = true;
	try {
		const uploaded = await mediaApi.uploadMedia(uploadFile);
		media = [uploaded, ...media];
		selectMedia(uploaded);
		uploadFile = null;
	} catch (err) {
		console.error("Upload failed:", err);
		alert("Failed to upload file");
	} finally {
		uploading = false;
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

const filteredMedia = $derived(
	media.filter((m) => {
		if (!searchQuery) return true;
		const query = searchQuery.toLowerCase();
		return (
			m.filename.toLowerCase().includes(query) ||
			m.original_filename.toLowerCase().includes(query)
		);
	}),
);
</script>

<div class="media-picker">
	{#if label}
		<label class="block text-sm font-medium text-gray-700 mb-2">
			{label}
			{#if required}
				<span class="text-red-500">*</span>
			{/if}
		</label>
	{/if}

	{#if selectedMedia}
		<div class="relative group border-2 border-gray-300 rounded-lg overflow-hidden">
			<img
				src={getMediaUrl(selectedMedia)}
				alt={selectedMedia.original_filename}
				class="w-full h-48 group-hover:h-80 object-cover"
			/>
			<div class="absolute inset-0 bg-black/20 bg-opacity-0 group-hover:bg-opacity-50 transition-all flex items-center justify-center gap-2">
				<button
					type="button"
					onclick={openPicker}
					class="opacity-0 group-hover:opacity-100 transition-opacity px-4 py-2 bg-white text-gray-900 rounded-lg font-medium hover:bg-gray-100"
				>
					Change
				</button>
				<button
					type="button"
					onclick={clearSelection}
					class="opacity-0 group-hover:opacity-100 transition-opacity px-4 py-2 bg-red-600 text-white rounded-lg font-medium hover:bg-red-700"
				>
					Remove
				</button>
			</div>
		</div>
		<p class="mt-2 text-sm text-gray-500">{selectedMedia.original_filename}</p>
	{:else}
		<button
			type="button"
			onclick={openPicker}
			class="w-full h-48 border-2 border-dashed border-gray-300 rounded-lg hover:border-gray-400 transition-colors flex flex-col items-center justify-center gap-3 text-gray-500 hover:text-gray-700"
		>
			<svg class="w-12 h-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
			</svg>
			<span class="font-medium">Select Image</span>
		</button>
	{/if}

	{#if error}
		<p class="mt-2 text-sm text-red-600">{error}</p>
	{/if}
</div>

{#if showModal}
	<div class="fixed inset-0 z-50 overflow-y-auto">
		<div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
			<!-- Overlay -->
			<div 
				class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"
				onclick={() => showModal = false}
				aria-hidden="true"
			></div>

			<!-- Modal Content - added 'relative' class -->
			<div class="relative inline-block w-full max-w-6xl my-8 overflow-hidden text-left align-middle transition-all transform bg-white rounded-lg shadow-xl">
				<!-- Header -->
				<div class="flex items-center justify-between px-6 py-4 border-b border-gray-200">
					<h3 class="text-lg font-semibold text-gray-900">Select Media</h3>
					<button
						type="button"
						onclick={() => showModal = false}
						class="text-gray-400 hover:text-gray-600"
					>
						<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
						</svg>
					</button>
				</div>

				<!-- Toolbar -->
				<div class="px-6 py-4 bg-gray-50 border-b border-gray-200">
					<div class="flex gap-4 items-center">
						<!-- Search -->
						<div class="relative flex-1">
							<svg class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
							</svg>
							<input
								type="text"
								bind:value={searchQuery}
								placeholder="Search media..."
								class="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
							/>
						</div>

						<!-- View Toggle -->
						<div class="flex border border-gray-300 rounded-lg overflow-hidden">
							<button
								type="button"
								onclick={() => viewMode = 'grid'}
								class="px-3 py-2 {viewMode === 'grid' ? 'bg-blue-600 text-white' : 'bg-white text-gray-700 hover:bg-gray-50'}"
							>
								<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
								</svg>
							</button>
							<button
								type="button"
								onclick={() => viewMode = 'list'}
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
								<svg class="animate-spin h-5 w-5" fill="none" viewBox="0 0 24 24">
									<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
									<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
								</svg>
								Uploading...
							{:else}
								<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
								</svg>
								Upload
							{/if}
							<input
								type="file"
								accept="image/*"
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
							<svg class="animate-spin h-8 w-8 text-gray-400" fill="none" viewBox="0 0 24 24">
								<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
								<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
							</svg>
						</div>
					{:else if filteredMedia.length === 0}
						<div class="text-center py-12 text-gray-500">
							<svg class="w-12 h-12 mx-auto mb-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
							</svg>
							<p>No media found</p>
						</div>
					{:else if viewMode === 'grid'}
						<!-- Grid View -->
						<div class="grid grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
							{#each filteredMedia as m (m.id)}
								<button
									type="button"
									onclick={() => selectMedia(m)}
									class="relative aspect-square rounded-lg overflow-hidden border-2 border-gray-200 hover:border-blue-500 transition-colors group"
								>
									<img
										src={getMediaUrl(m)}
										alt={m.original_filename}
										class="w-full h-full object-cover"
									/>
									<div class="absolute inset-0 bg-black/0 group-hover:bg-black/30 transition-all flex items-center justify-center">
										<svg class="w-8 h-8 text-white opacity-0 group-hover:opacity-100 transition-opacity" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
										</svg>
									</div>
								</button>
							{/each}
						</div>
					{:else}
						<!-- List View -->
						<div class="bg-white rounded-lg shadow overflow-hidden">
							<table class="min-w-full divide-y divide-gray-200">
								<thead class="bg-gray-50">
									<tr>
										<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Preview</th>
										<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Filename</th>
										<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Size</th>
										<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Uploaded</th>
										<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Action</th>
									</tr>
								</thead>
								<tbody class="bg-white divide-y divide-gray-200">
									{#each filteredMedia as m (m.id)}
										<tr class="hover:bg-gray-50 transition-colors">
											<td class="px-6 py-4 whitespace-nowrap">
												<img
													src={getMediaUrl(m)}
													alt={m.original_filename}
													class="h-12 w-12 object-cover rounded"
												/>
											</td>
											<td class="px-6 py-4">
												<div class="text-sm font-medium text-gray-900 truncate max-w-xs">
													{m.original_filename}
												</div>
											</td>
											<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
												{formatFileSize(m.size_bytes)}
											</td>
											<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
												{formatDate(m.created_at)}
											</td>
											<td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
												<button
													type="button"
													onclick={() => selectMedia(m)}
													class="text-blue-600 hover:text-blue-900"
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
				</div>

				<!-- Pagination -->
				{#if totalPages > 1}
					<div class="px-6 py-4 border-t border-gray-200 flex items-center justify-between">
						<button
							type="button"
							onclick={() => changePage(currentPage - 1)}
							disabled={currentPage === 1}
							class="px-4 py-2 border border-gray-300 rounded-lg bg-white text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed font-medium transition-colors"
						>
							Previous
						</button>
						<span class="text-sm text-gray-600">
							Page {currentPage} of {totalPages}
						</span>
						<button
							type="button"
							onclick={() => changePage(currentPage + 1)}
							disabled={currentPage === totalPages}
							class="px-4 py-2 border border-gray-300 rounded-lg bg-white text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed font-medium transition-colors"
						>
							Next
						</button>
					</div>
				{/if}
			</div>
		</div>
	</div>
{/if}
