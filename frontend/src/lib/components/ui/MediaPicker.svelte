<script lang="ts">
	import { PUBLIC_API_URL } from "$env/static/public";
	import { mediaApi, type Media, getMediaUrl } from "$api/media";
	import {
		ImageUp,
		CheckCircle,
		AlertCircle,
		Loader2,
		Plus,
		Loader,
		Funnel,
	} from "lucide-svelte";
	import Pagination from "./Pagination.svelte";

	interface Props {
		show?: boolean;
		onselect?: (mediaId: string, media: Media) => void;
		onclose?: () => void;
	}

	let { show, onselect, onclose }: Props = $props();

	// Media
	let media = $state<Media[]>([]);
	let loading = $state(false);
	let searchQuery = $state("");

	// Uploading
	let uploadFile = $state<File | null>(null);
	let uploading = $state(false);
	let uploadProgress = $state(0);
	let uploadStatus = $state<"uploading" | "processing" | "success" | "error">(
		"uploading",
	);
	let uploadError = $state<string>("");

	// View Mode & Pagination
	let viewMode = $state<"grid" | "list">("grid");
	let currentPage = $state(1);
	let totalPages = $state(1);
	let limit = $state(20);

	let total = $state(0);
	let typeFilter = $state<string>("");
	let sortBy = $state<string>("created_at");
	let sortOrder = $state<string>("desc");
	let showFilters = $state(false);

	const ACCEPTED_FILE_TYPES =
		"image/*,video/*,video/mp4,video/quicktime,.mp4,.mov,.avi,.gif";

	// Load media when modal opens
	$effect(() => {
		if (show && media.length === 0) {
			loadMedia();
		}
	});

	async function loadMedia() {
		loading = true;
		try {
			const params = new URLSearchParams({
				page: currentPage.toString(),
				limit: limit.toString(),
				sort: sortBy,
				order: sortOrder,
			});

			if (searchQuery.trim()) {
				params.append("search", searchQuery.trim());
			}

			if (typeFilter) {
				params.append("type", typeFilter);
			}

			const response = await fetch(
				`${PUBLIC_API_URL}/api/v1/media?${params.toString()}`,
				{
					credentials: "include",
				},
			);

			if (!response.ok) {
				throw new Error("Failed to load media");
			}

			const data = await response.json();
			media = data.data || [];
			if (data.pagination) {
				totalPages = data.pagination.total_pages || 1;
				total = data.pagination.total || 0;
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

	async function applyFilters() {
		currentPage = 1;
		await loadMedia();
	}

	async function changeLimit(newLimit: number) {
		limit = newLimit;
		currentPage = 1;
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
		const maxSize = uploadFile.type.startsWith("video/") ? 200 : 50;
		const fileSizeMB = uploadFile.size / (1024 * 1024);

		if (fileSizeMB > maxSize) {
			alert(
				`File size (${fileSizeMB.toFixed(1)}MB) exceeds ${maxSize}MB limit`,
			);
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
					if (uploadProgress === 100) uploadStatus = "processing";
				}
			});

			const uploaded = await new Promise<Media>((resolve, reject) => {
				xhr.addEventListener("load", () => {
					if (xhr.status === 201) {
						uploadStatus = "success";
						resolve(JSON.parse(xhr.responseText));
					} else {
						uploadStatus = "error";
						reject(
							new Error(JSON.parse(xhr.responseText).error || "Upload failed"),
						);
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
			uploadFile = null;
			selectMedia(uploaded);
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
		onselect?.(m.id, m);
		closeModal();
	}

	function closeModal() {
		show = false;
		onclose?.();
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

{#if show}
	<div class="fixed inset-0 z-50 overflow-y-auto" onclick={closeModal}>
		<div class="flex min-h-screen items-center justify-center p-4">
			<div class="fixed inset-0 bg-black/50 transition-opacity"></div>

			<div class="relative bg-surface rounded-lg shadow-xl max-w-6xl w-full max-h-[90vh] overflow-hidden" onclick={(e) => e.stopPropagation()}>
				<div class="border-b border-gray-200 px-6 py-4">
					<div class="flex items-center justify-between">
						<div class="flex-1">
							<h2 class="text-xl font-semibold text-gray-900">Select Media</h2>
							<p class="text-sm text-gray-500 mt-1">Upload new or select existing media</p>
						</div>
            <!-- Filters Panel -->
            {#if showFilters}
              <div class="px-6 py-4 absolute z-50 h-auto bg-surface border-b border-input-border inset-2">
                <div class="grid grid-cols-1 gap-4">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">Type</label>
                    <select
                      bind:value={typeFilter}
                      onchange={applyFilters}
                      class="form-input"
                    >
                      <option value="">All Types</option>
                      <option value="image">Images</option>
                      <option value="video">Videos</option>
                    </select>
                  </div>

                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">Sort By</label>
                    <select
                      bind:value={sortBy}
                      onchange={applyFilters}
                      class="form-input"
                    >
                      <option value="created_at">Date Uploaded</option>
                      <option value="filename">Filename</option>
                      <option value="size">File Size</option>
                    </select>
                  </div>

                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">Order</label>
                    <select
                      bind:value={sortOrder}
                      onchange={applyFilters}
                      class="form-input"
                    >
                      <option value="desc">Descending</option>
                      <option value="asc">Ascending</option>
                    </select>
                  </div>
                </div>
              </div>
            {/if}

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

						<button onclick={closeModal} class="text-gray-400 hover:text-gray-600">
							<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
							</svg>
						</button>
					</div>

					<div class="flex gap-3 mt-4">
						<input type="text" bind:value={searchQuery} placeholder="Search media..." class="form-input grow-0 py-1" />
            <button
              onclick={() => (showFilters = !showFilters)}
              class="px-3 py-2 rounded-lg hover:bg-gray-100 {showFilters ? 'bg-primary text-accent' : 'text-gray-700'}"
              title="Toggle filters"
            >
              <Funnel class="w-5 h-5" />
            </button>
						
						<div class="flex">
							<button onclick={() => (viewMode = "grid")} class="px-3 py-2 rounded-l-md {viewMode === 'grid' ? 'bg-primary text-white' : 'bg-white text-gray-700 hover:bg-gray-50'}">
								<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
								</svg>
							</button>
							<button onclick={() => (viewMode = "list")} class="px-3 py-2 rounded-r-md {viewMode === 'list' ? 'bg-primary text-white' : 'bg-white text-gray-700 hover:bg-gray-50'}">
								<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
								</svg>
							</button>
						</div>

						<label class="px-4 py-2 bg-primary text-white rounded-lg font-medium hover:bg-primary/80 cursor-pointer flex items-center gap-2">
							{#if uploading}
								<Loader class="w-5 h-5 animate-spin" />
								Uploading...
							{:else}
								<ImageUp class="w-5 h-5" />
								Upload
							{/if}
							<input type="file" accept={ACCEPTED_FILE_TYPES} onchange={handleFileSelect} class="hidden" disabled={uploading} />
						</label>
					</div>
				</div>

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
								<button onclick={() => selectMedia(m)} class="group relative aspect-square rounded-lg overflow-hidden bg-gray-100 hover:ring-2 hover:ring-blue-500 transition-all">
									<img src={m.thumbnail_url || getMediaUrl(m)} alt={m.original_filename} class="w-full h-full object-cover" />
									<div class="absolute top-2 right-2 flex gap-1">
										{#if m.mime_type === 'image/webp'}
											<span class="px-2 py-1 bg-green-500 text-white text-xs rounded-full font-medium">WebP</span>
										{/if}
                    {console.log(m.mime_type)}
                    {#if m.mime_type === 'image/gif'}
											<span class="px-2 py-1 bg-blue-500 text-white text-xs rounded-full font-medium">GIF</span>
										{/if}
										{#if isVideo(m)}
											<span class="px-2 py-1 bg-purple-500 text-white text-xs rounded-full font-medium">Video</span>
										{/if}
									</div>
									<div class="absolute inset-0 bg-black/0 group-hover:bg-black/50 transition-opacity flex items-end">
										<div class="w-full p-3 text-white opacity-0 group-hover:opacity-100 transition-opacity">
											<p class="text-sm font-medium truncate">{m.original_filename}</p>
											<p class="text-xs">{formatFileSize(m.size_bytes)}</p>
										</div>
									</div>
								</button>
							{/each}
						</div>
					{:else}
						<div class="bg-surface rounded-lg shadow overflow-hidden">
							<table class="min-w-full divide-y divide-gray-200">
								<thead class="bg-gray-600">
									<tr>
										<th class="px-6 py-3 text-left text-xs font-medium text-foreground uppercase">Preview</th>
										<th class="px-6 py-3 text-left text-xs font-medium text-foreground uppercase">Filename</th>
										<th class="px-6 py-3 text-left text-xs font-medium text-foreground uppercase">Size</th>
										<th class="px-6 py-3 text-left text-xs font-medium text-foreground uppercase">Type</th>
										<th class="px-6 py-3 text-left text-xs font-medium text-foreground uppercase">Action</th>
									</tr>
								</thead>
								<tbody class="divide-y divide-gray-200">
									{#each filteredMedia as m (m.id)}
										<tr class="hover:bg-gray-700 bg-neutral-800">
											<td class="px-6 py-4">
												<img src={m.thumbnail_url || getMediaUrl(m)} alt={m.original_filename} class="h-12 w-12 object-cover rounded" />
											</td>
											<td class="px-6 py-4">
												<div class="text-sm font-medium text-foreground truncate max-w-xs">{m.original_filename}</div>
											</td>
											<td class="px-6 py-4 text-sm text-foreground">{formatFileSize(m.size_bytes)}</td>
											<td class="px-6 py-4">
												{#if m.mime_type === 'image/webp'}
													<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">WebP</span>
												{:else if isVideo(m)}
													<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-800">Video</span>
                        {:else if m.mime_type === 'image/gif'}
													<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-500 text-white">{m.mime_type.split('/')[1].toUpperCase()}</span>
												{/if}
											</td>
											<td class="px-6 py-4">
												<button onclick={() => selectMedia(m)} class="text-white font-medium text-sm btn btn-primary bg-transparent hover:bg-primary px-4 py-2"><Plus size={16} /></button>
											</td>
										</tr>
									{/each}
								</tbody>
							</table>
						</div>
					{/if}

          <!-- Enhanced Pagination Footer -->
          {#if !loading && media.length > 0}
            <div class="px-6 py-4 border-t border-gray-200 flex items-center justify-between">
              <div class="flex items-center gap-4 w-full">
                <span class="text-sm text-gray-300">
                  Showing {(currentPage - 1) * limit + 1} to {Math.min(currentPage * limit, total)} of {total} results
                </span>
                <select
                  bind:value={limit}
                  onchange={() => changeLimit(limit)}
                  class="form-input max-w-40 py-1"
                >
                  <option value={10}>10 per page</option>
                  <option value={20}>20 per page</option>
                  <option value={50}>50 per page</option>
                  <option value={100}>100 per page</option>
                </select>
              </div>

              <!-- Pagination -->
              <Pagination 
                currentPage={currentPage}
                totalPages={totalPages} 
                onPageChange={changePage} 
              />



            </div>
          {/if}			
        </div>
			</div>
		</div>
	</div>
{/if}
