<script lang="ts">
	import { X } from '@lucide/svelte';
	import { fetchMedia } from '$lib/api/media';
	import { toast } from 'svelte-sonner';
	import type { MediaItem } from '$lib/types';
	import Pagination from '../Navigation/Pagination.svelte';

	let {
		open,
		projectSlug,
		onclose,
		onlinked,
		linkedMediaIds = [],
		selectOnly = false,
		mediaPool = null
	} = $props<{
		open: boolean;
		projectSlug: string;
		onclose: () => void;
		onlinked: (media: MediaItem) => void;
		linkedMediaIds: string[];
		selectOnly?: boolean;
		mediaPool?: MediaItem[] | null;
	}>();

	let loading = $state(false);
	let page = $state(1);
	let totalPages = $state(1);
	let pageSize = 10;
	let totalMedia = $state(0);
	let allUnlinkedMedia: MediaItem[] = [];
	let filteredMedia = $state<MediaItem[]>([]);

	let hasLoaded = false;

	// Calculate pagination counts
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
				const res = await fetchMedia(1, 1000); // get all for pool pagination
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

	// Runs once on modal open
	$effect(() => {
		if (open && !hasLoaded) {
			page = 1;
			loadMedia();
			hasLoaded = true;
		}
	});

	// Reset flag on close
	$effect(() => {
		if (!open) {
			hasLoaded = false;
		}
	});

	function handlePageChange(newPage: number) {
		if (newPage >= 1 && newPage <= totalPages && newPage !== page) {
			page = newPage;
			paginateMediaPool();
		}
	}

	async function linkMedia(item: MediaItem) {
		if (selectOnly) {
			onlinked(item);
			return;
		}

		try {
			const res = await fetch(`/api/v1/admin/projects/${projectSlug}/media`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ media_id: item.id, sort_order: 0 })
			});

			if (!res.ok) throw new Error('Failed to link media');
			toast.success('Media linked');

			// Remove item from pool
			allUnlinkedMedia = allUnlinkedMedia.filter((m) => m.id !== item.id);

			// Recalculate pagination (handles backfilling)
			const currentStart = (page - 1) * pageSize;
			const currentEnd = page * pageSize;

			filteredMedia = allUnlinkedMedia.slice(currentStart, currentEnd);

			// Update pagination values
			totalMedia = allUnlinkedMedia.length;
			totalPages = Math.max(1, Math.ceil(totalMedia / pageSize));

			// If we overshot the page count (e.g. linked last item on last page)
			if (page > totalPages) {
				page = totalPages;
				filteredMedia = allUnlinkedMedia.slice((page - 1) * pageSize, page * pageSize);
			}

			onlinked(item);
		} catch (err) {
			console.error('❌ Failed to link media:', err);
			toast.error('Failed to link media');
		}
	}
</script>

{#if open}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
		<div class="relative w-full max-w-3xl rounded-lg bg-white shadow-xl">
			<div class="p-4">
				<div class="flex items-center justify-between border-b pb-2">
					<h2 class="text-lg font-semibold">Select Media</h2>
					<button onclick={onclose} class="text-gray-500 hover:text-red-500">
						<X />
					</button>
				</div>

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
			</div>

			<Pagination {page} {totalPages} {pageSize} {totalMedia} onchange={handlePageChange} />
		</div>
	</div>
{/if}
