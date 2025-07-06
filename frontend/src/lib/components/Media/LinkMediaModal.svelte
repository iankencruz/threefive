<script lang="ts">
	import { X } from '@lucide/svelte';
	import { fetchMedia } from '$lib/api/media';
	import { toast } from 'svelte-sonner';
	import type { MediaItem } from '$lib/types';
	import Pagination from '../Navigation/Pagination.svelte';

	let {
		open,
		context,
		onclose,
		onlinked,
		linkedMediaIds = [],
		selectOnly = false,
		mediaPool = null
	} = $props<{
		open: boolean;
		context: { type: 'project' | 'block' | 'page'; id: string };
		onclose: () => void;
		onlinked: (media: MediaItem) => void;
		linkedMediaIds: string[];
		selectOnly?: boolean;
		mediaPool?: MediaItem[] | null;
	}>();

	let tab = $state<'link' | 'upload'>('link');
	let loading = $state(false);
	let page = $state(1);
	let totalPages = $state(1);
	let pageSize = 10;
	let totalMedia = $state(0);
	let allUnlinkedMedia: MediaItem[] = [];
	let filteredMedia = $state<MediaItem[]>([]);
	let hasLoaded = false;

	function paginateMediaPool() {
		const start = (page - 1) * pageSize;
		const end = page * pageSize;

		filteredMedia = allUnlinkedMedia.slice(start, end);
		totalMedia = allUnlinkedMedia.length;
		totalPages = Math.max(1, Math.ceil(totalMedia / pageSize));
	}

	async function loadMedia() {
		loading = true;
		try {
			if (mediaPool) {
				allUnlinkedMedia = selectOnly
					? [...mediaPool]
					: mediaPool.filter((m: MediaItem) => !linkedMediaIds.includes(m.id));
			} else {
				const res = await fetchMedia(1, 1000);
				allUnlinkedMedia = selectOnly
					? res.items
					: res.items.filter((m: MediaItem) => !linkedMediaIds.includes(m.id));
			}

			paginateMediaPool();
		} catch (err) {
			console.error('❌ Failed to load media:', err);
			toast.error('Could not load media');
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		if (open && !hasLoaded) {
			tab = 'link';
			page = 1;
			loadMedia();
			hasLoaded = true;
		}
	});

	$effect(() => {
		if (!open) hasLoaded = false;
	});

	function handlePageChange(newPage: number) {
		if (newPage >= 1 && newPage <= totalPages && newPage !== page) {
			page = newPage;
			paginateMediaPool();
		}
	}

	function linkMedia(item: MediaItem) {
		onlinked(item);
		if (selectOnly) onclose();
		allUnlinkedMedia = allUnlinkedMedia.filter((m) => m.id !== item.id);
		paginateMediaPool();
	}
</script>

{#if open}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
		<div class="relative w-full max-w-4xl rounded-lg bg-white shadow-xl">
			<div class="p-4">
				<!-- Tabs -->
				<div class="mt-4 flex gap-2 border-b text-sm font-medium">
					<button
						onclick={() => (tab = 'link')}
						class="rounded-t px-4 py-2"
						class:font-bold={tab === 'link'}
						class:border-b-2={tab === 'link'}
					>
						Link Media
					</button>
					<button
						onclick={() => (tab = 'upload')}
						class="rounded-t px-4 py-2"
						class:font-bold={tab === 'upload'}
						class:border-b-2={tab === 'upload'}
					>
						Upload Media
					</button>
				</div>

				<!-- Tab Content -->
				{#if tab === 'link'}
					{#if loading}
						<p class="p-4 text-gray-500">Loading media...</p>
					{:else if filteredMedia.length === 0}
						<p class="p-4 text-gray-500">No media available to link.</p>
					{:else}
						<ul
							class="mt-4 grid grid-cols-2 gap-x-4 gap-y-8 sm:grid-cols-3 md:grid-cols-4 xl:grid-cols-6 xl:gap-x-8"
						>
							{#each filteredMedia as item (item.id)}
								<li
									class="group relative block overflow-hidden rounded-lg bg-gray-100 ring-1 ring-gray-200 hover:cursor-pointer"
								>
									<button class="h-auto w-full" onclick={() => linkMedia(item)}>
										<img
											src={item.thumbnail_url || item.url}
											alt={item.title}
											class="aspect-video w-full object-cover"
										/>
										<div class="p-2">
											<p class="truncate text-sm font-medium text-gray-900">
												{item.title || 'Untitled'}
											</p>
										</div>
									</button>
								</li>
							{/each}
						</ul>
					{/if}
					<Pagination {page} {totalPages} {pageSize} {totalMedia} onchange={handlePageChange} />
				{:else if tab === 'upload'}
					<div class="mt-4">
						<UploadMediaForm
							{context}
							onuploaded={(media: MediaItem) => {
								toast.success('Media uploaded');
								onlinked(media);
								loadMedia(); // Refresh list to exclude newly linked
							}}
						/>
					</div>
				{/if}
			</div>
		</div>
	</div>
{/if}
