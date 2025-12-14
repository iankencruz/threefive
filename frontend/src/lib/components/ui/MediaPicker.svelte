<script lang="ts">
	import { PUBLIC_API_URL } from '$env/static/public';
	import { type Media, getMediaUrl } from '$api/media';
	import { ImageUp, Plus, Loader, Funnel, CircleCheck, CircleAlert } from 'lucide-svelte';
	import Pagination from './Pagination.svelte';

	interface Props {
		show?: boolean;
		onselect?: (mediaId: string, media: Media) => void;
		onclose?: () => void;
	}

	let { show, onselect, onclose }: Props = $props();

	// Media
	let media = $state<Media[]>([]);
	let loading = $state(false);
	let searchQuery = $state('');

	// Uploading
	let uploadFile = $state<File | null>(null);
	let uploading = $state(false);
	let uploadProgress = $state(0);
	let uploadStatus = $state<'uploading' | 'processing' | 'success' | 'error'>('uploading');
	let uploadError = $state<string>('');

	// View Mode & Pagination
	let viewMode = $state<'grid' | 'list'>('grid');
	let currentPage = $state(1);
	let totalPages = $state(1);
	let limit = $state(20);

	let total = $state(0);
	let typeFilter = $state<string>('');
	let sortBy = $state<string>('created_at');
	let sortOrder = $state<string>('desc');
	let showFilters = $state(false);

	const ACCEPTED_FILE_TYPES = 'image/*,video/*,video/mp4,video/quicktime,.mp4,.mov,.avi,.gif';

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
				order: sortOrder
			});

			if (searchQuery.trim()) {
				params.append('search', searchQuery.trim());
			}

			if (typeFilter) {
				params.append('type', typeFilter);
			}

			const response = await fetch(`${PUBLIC_API_URL}/api/v1/media?${params.toString()}`, {
				credentials: 'include'
			});

			if (!response.ok) {
				throw new Error('Failed to load media');
			}

			const data = await response.json();
			media = data.data || [];
			if (data.pagination) {
				totalPages = data.pagination.total_pages || 1;
				total = data.pagination.total || 0;
			}
		} catch (err) {
			console.error('Failed to load media:', err);
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
		const maxSize = uploadFile.type.startsWith('video/') ? 200 : 50;
		const fileSizeMB = uploadFile.size / (1024 * 1024);

		if (fileSizeMB > maxSize) {
			alert(`File size (${fileSizeMB.toFixed(1)}MB) exceeds ${maxSize}MB limit`);
			return;
		}

		uploading = true;
		uploadProgress = 0;
		uploadStatus = 'uploading';
		uploadError = '';

		try {
			const formData = new FormData();
			formData.append('file', uploadFile);
			const xhr = new XMLHttpRequest();

			xhr.upload.addEventListener('progress', (e) => {
				if (e.lengthComputable) {
					uploadProgress = Math.round((e.loaded / e.total) * 100);
					if (uploadProgress === 100) uploadStatus = 'processing';
				}
			});

			const uploaded = await new Promise<Media>((resolve, reject) => {
				xhr.addEventListener('load', () => {
					if (xhr.status === 201) {
						uploadStatus = 'success';
						resolve(JSON.parse(xhr.responseText));
					} else {
						uploadStatus = 'error';
						reject(new Error(JSON.parse(xhr.responseText).error || 'Upload failed'));
					}
				});
				xhr.addEventListener('error', () => {
					uploadStatus = 'error';
					reject(new Error('Network error during upload'));
				});
				xhr.open('POST', `${PUBLIC_API_URL}/api/v1/media/upload`);
				xhr.withCredentials = true;
				xhr.send(formData);
			});

			media = [uploaded, ...media];
			uploadFile = null;
			selectMedia(uploaded);
		} catch (err) {
			console.error('Upload failed:', err);
			uploadError = err instanceof Error ? err.message : 'Upload failed';
			uploadStatus = 'error';
		} finally {
			setTimeout(() => {
				uploading = false;
				uploadProgress = 0;
				uploadStatus = 'uploading';
				uploadError = '';
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

	function getProcessingMessage(): string {
		if (!uploadFile) return 'Processing...';
		if (isImage({ mime_type: uploadFile.type } as Media)) {
			return 'Converting to WebP and generating thumbnail...';
		} else if (isVideo({ mime_type: uploadFile.type } as Media)) {
			return 'Optimizing video and extracting thumbnail...';
		}
		return 'Processing file...';
	}

	const filteredMedia = $derived(
		media.filter((m) => {
			if (!searchQuery.trim()) return true;
			return m.original_filename.toLowerCase().includes(searchQuery.toLowerCase());
		})
	);
</script>

{#if show}
	<div
		class="fixed inset-0 z-50 overflow-y-auto"
		role="button"
		tabindex="0"
		onkeydown={(e) => {
			if (e.key === 'Enter' || e.key === ' ') {
				closeModal();
			}
		}}
		onclick={closeModal}
	>
		<div class="flex min-h-screen items-center justify-center p-4">
			<div class="fixed inset-0 bg-black/50 transition-opacity"></div>

			<div
				class="relative max-h-[90vh] w-full max-w-6xl overflow-hidden rounded-lg bg-surface shadow-xl"
				onclick={(e) => e.stopPropagation()}
				role="button"
				tabindex="0"
				onkeydown={(e) => {
					// Execute the click logic on Enter (13) or Space (32) key press
					if (e.key === 'Enter' || e.key === ' ') {
						e.stopPropagation();
						// Call the same function/logic as the click handler
						// For example: handleSelectMedia(e);
					}
				}}
			>
				<div class="border-b border-gray-200 px-6 py-4">
					<div class="flex items-center justify-between">
						<div class="flex-1">
							<h2 class="text-xl font-semibold">Select Media</h2>
							<p class="mt-1 text-sm text-gray-500">Upload new or select existing media</p>
						</div>
						<!-- Filters Panel -->
						{#if showFilters}
							<div
								class="absolute inset-2 z-50 h-auto border-b border-input-border bg-surface px-6 py-4"
							>
								<div class="grid grid-cols-1 gap-4">
									<div>
										<label for="filter-type" class="mb-2 block text-sm font-medium text-gray-700"
											>Type</label
										>
										<select
											name="filter-type"
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
										<label for="sort-by" class="mb-2 block text-sm font-medium text-gray-700"
											>Sort By</label
										>
										<select
											name="sort-by"
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
										<label for="sort-order" class="mb-2 block text-sm font-medium text-gray-700"
											>Order</label
										>
										<select
											name="sort-order"
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
							<div class="mr-4 flex items-center gap-3 rounded-lg bg-blue-50 px-4 py-2">
								{#if uploadStatus === 'uploading'}
									<Loader class="h-5 w-5 animate-spin text-blue-600" />
									<div>
										<p class="text-sm font-medium text-blue-900">Uploading...</p>
										<p class="text-xs text-blue-600">{uploadProgress}%</p>
									</div>
								{:else if uploadStatus === 'processing'}
									<Loader class="h-5 w-5 animate-spin text-yellow-600" />
									<div>
										<p class="text-sm font-medium text-yellow-900">Processing</p>
										<p class="text-xs text-yellow-600">{getProcessingMessage()}</p>
									</div>
								{:else if uploadStatus === 'success'}
									<CircleCheck class="h-5 w-5 text-green-600" />
									<p class="text-sm font-medium text-green-900">Success!</p>
								{:else if uploadStatus === 'error'}
									<CircleAlert class="h-5 w-5 text-red-600" />
									<div>
										<p class="text-sm font-medium text-red-900">Error</p>
										<p class="text-xs text-red-600">{uploadError}</p>
									</div>
								{/if}
							</div>
						{/if}

						<button
							aria-label="Close Media Picker Modal"
							onclick={closeModal}
							class="text-gray-400 hover:text-gray-600"
						>
							<svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M6 18L18 6M6 6l12 12"
								/>
							</svg>
						</button>
					</div>

					<div class="mt-4 flex gap-3">
						<input
							type="text"
							bind:value={searchQuery}
							placeholder="Search media..."
							class="form-input grow-0 py-1"
						/>
						<button
							onclick={() => (showFilters = !showFilters)}
							class="rounded-lg px-3 py-2 hover:bg-gray-100 {showFilters
								? 'bg-primary text-accent'
								: 'text-gray-700'}"
							title="Toggle filters"
						>
							<Funnel class="h-5 w-5" />
						</button>

						<div class="flex">
							<button
								aria-label="Grid View Mode"
								type="button"
								onclick={() => (viewMode = 'grid')}
								class="rounded-l-md px-3 py-2 {viewMode === 'grid'
									? 'bg-primary text-white'
									: 'bg-white text-gray-700 hover:bg-gray-50'}"
							>
								<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"
									/>
								</svg>
							</button>
							<button
								aria-label="List View Mode"
								type="button"
								onclick={() => (viewMode = 'list')}
								class="rounded-r-md px-3 py-2 {viewMode === 'list'
									? 'bg-primary text-white'
									: 'bg-white text-gray-700 hover:bg-gray-50'}"
							>
								<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M4 6h16M4 12h16M4 18h16"
									/>
								</svg>
							</button>
						</div>

						<label
							class="flex cursor-pointer items-center gap-2 rounded-lg bg-primary px-4 py-2 font-medium text-white hover:bg-primary/80"
						>
							{#if uploading}
								<Loader class="h-5 w-5 animate-spin" />
								Uploading...
							{:else}
								<ImageUp class="h-5 w-5" />
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

				<div class="max-h-[600px] overflow-y-auto px-6 py-6">
					{#if loading}
						<div class="flex items-center justify-center py-12">
							<Loader class="h-8 w-8 animate-spin text-gray-400" />
						</div>
					{:else if filteredMedia.length === 0}
						<div class="py-12 text-center">
							<ImageUp class="mx-auto mb-4 h-12 w-12 text-gray-400" />
							<p class="text-gray-500">No media found. Upload your first file!</p>
						</div>
					{:else if viewMode === 'grid'}
						<div class="grid grid-cols-2 gap-4 md:grid-cols-3 lg:grid-cols-4">
							{#each filteredMedia as m (m.id)}
								<button
									onclick={() => selectMedia(m)}
									class="group relative aspect-square overflow-hidden rounded-lg bg-gray-100 transition-all hover:ring-2 hover:ring-blue-500"
								>
									<img
										src={m.thumbnail_url || getMediaUrl(m)}
										alt={m.original_filename}
										class="h-full w-full object-cover"
									/>
									<div class="absolute top-2 right-2 flex gap-1">
										{#if m.mime_type === 'image/webp'}
											<span
												class="rounded-full bg-green-500 px-2 py-1 text-xs font-medium text-white"
												>WebP</span
											>
										{/if}
										{console.log(m.mime_type)}
										{#if m.mime_type === 'image/gif'}
											<span
												class="rounded-full bg-blue-500 px-2 py-1 text-xs font-medium text-white"
												>GIF</span
											>
										{/if}
										{#if isVideo(m)}
											<span
												class="rounded-full bg-purple-500 px-2 py-1 text-xs font-medium text-white"
												>Video</span
											>
										{/if}
									</div>
									<div
										class="absolute inset-0 flex items-end bg-black/0 transition-opacity group-hover:bg-black/50"
									>
										<div
											class="w-full p-3 text-white opacity-0 transition-opacity group-hover:opacity-100"
										>
											<p class="truncate text-sm font-medium">{m.original_filename}</p>
											<p class="text-xs">{formatFileSize(m.size_bytes)}</p>
										</div>
									</div>
								</button>
							{/each}
						</div>
					{:else}
						<div class="overflow-hidden rounded-lg bg-surface shadow">
							<table class="min-w-full divide-y divide-gray-200">
								<thead class="bg-gray-600">
									<tr>
										<th class="px-6 py-3 text-left text-xs font-medium text-foreground uppercase"
											>Preview</th
										>
										<th class="px-6 py-3 text-left text-xs font-medium text-foreground uppercase"
											>Filename</th
										>
										<th class="px-6 py-3 text-left text-xs font-medium text-foreground uppercase"
											>Size</th
										>
										<th class="px-6 py-3 text-left text-xs font-medium text-foreground uppercase"
											>Type</th
										>
										<th class="px-6 py-3 text-left text-xs font-medium text-foreground uppercase"
											>Action</th
										>
									</tr>
								</thead>
								<tbody class="divide-y divide-gray-200">
									{#each filteredMedia as m (m.id)}
										<tr class="bg-neutral-800 hover:bg-gray-700">
											<td class="px-6 py-4">
												<img
													src={m.thumbnail_url || getMediaUrl(m)}
													alt={m.original_filename}
													class="h-12 w-12 rounded object-cover"
												/>
											</td>
											<td class="px-6 py-4">
												<div class="max-w-xs truncate text-sm font-medium text-foreground">
													{m.original_filename}
												</div>
											</td>
											<td class="px-6 py-4 text-sm text-foreground"
												>{formatFileSize(m.size_bytes)}</td
											>
											<td class="px-6 py-4">
												{#if m.mime_type === 'image/webp'}
													<span
														class="inline-flex items-center rounded-full bg-green-100 px-2.5 py-0.5 text-xs font-medium text-green-800"
														>WebP</span
													>
												{:else if isVideo(m)}
													<span
														class="inline-flex items-center rounded-full bg-purple-100 px-2.5 py-0.5 text-xs font-medium text-purple-800"
														>Video</span
													>
												{:else if m.mime_type === 'image/gif'}
													<span
														class="inline-flex items-center rounded-full bg-blue-500 px-2.5 py-0.5 text-xs font-medium text-white"
														>{m.mime_type.split('/')[1].toUpperCase()}</span
													>
												{/if}
											</td>
											<td class="px-6 py-4">
												<button
													onclick={() => selectMedia(m)}
													class="btn btn-primary bg-transparent px-4 py-2 text-sm font-medium text-white hover:bg-primary"
													><Plus size={16} /></button
												>
											</td>
										</tr>
									{/each}
								</tbody>
							</table>
						</div>
					{/if}

					<!-- Enhanced Pagination Footer -->
					{#if !loading && media.length > 0}
						<div class="flex items-center justify-between border-t border-gray-200 px-6 py-4">
							<div class="flex w-full items-center gap-4">
								<span class="text-sm text-gray-300">
									Showing {(currentPage - 1) * limit + 1} to {Math.min(currentPage * limit, total)} of
									{total} results
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
							<Pagination {currentPage} {totalPages} onPageChange={changePage} />
						</div>
					{/if}
				</div>
			</div>
		</div>
	</div>
{/if}
