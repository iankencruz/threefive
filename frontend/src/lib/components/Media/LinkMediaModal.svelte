<script lang="ts">
	import { X } from '@lucide/svelte';
	import { fetchMedia } from '$lib/api/media';
	import { toast } from 'svelte-sonner';
	import type { MediaItem } from '$lib/types';

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
	let allMedia = $state<MediaItem[]>([]);
	let filteredMedia = $state<MediaItem[]>([]);

	$effect(() => {
		if (open) {
			(async () => {
				loading = true;
				try {
					if (mediaPool) {
						allMedia = mediaPool;
					} else {
						const data = await fetchMedia();
						allMedia = data.items;
					}
					filteredMedia = selectOnly
						? allMedia
						: allMedia.filter((item: MediaItem) => !linkedMediaIds.includes(item.id));
				} catch (err) {
					console.error('❌ Failed to load media:', err);
					toast.error('Could not load media');
				} finally {
					loading = false;
				}
			})();
		}
	});

	async function linkMedia(item: MediaItem) {
		if (selectOnly) {
			onlinked(item); // just emit and close
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

			// remove from modal grid
			filteredMedia = filteredMedia.filter((m) => m.id !== item.id);
			onlinked(item);
		} catch (err) {
			console.error('❌ Failed to link media:', err);
			toast.error('Failed to link media');
		}
	}
</script>

{#if open}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
		<div class="w-full max-w-3xl rounded-lg bg-white p-4 shadow-xl">
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
	</div>
{/if}
