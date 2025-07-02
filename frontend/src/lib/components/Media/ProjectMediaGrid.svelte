<script lang="ts">
	import type { MediaItem } from '$src/lib/types';
	import { Link, Unlink } from '@lucide/svelte';
	import { draggable, droppable, type DragDropState } from '@thisux/sveltednd';
	import { flip } from 'svelte/animate';
	import { toast } from 'svelte-sonner';

	let { media, onremove, onrefresh, onlink, onsort } = $props<{
		media: MediaItem[];
		onremove?: (id: string) => void;
		onrefresh?: () => void;
		onlink?: () => void;
		onsort?: (items: MediaItem[]) => void;
	}>();

	function handleDrop(state: DragDropState<MediaItem>) {
		const { draggedItem, sourceContainer, targetContainer } = state;
		if (!targetContainer || sourceContainer === targetContainer) return;

		const draggedId = draggedItem.id;
		const fromIndex = media.findIndex((m: MediaItem) => m.id === draggedId);
		const toIndex = parseInt(targetContainer);

		if (fromIndex === -1 || toIndex === -1 || fromIndex === toIndex) return;

		const updated = [...media];
		const [moved] = updated.splice(fromIndex, 1);
		updated.splice(toIndex, 0, moved);

		// Reassign sort_order
		const reordered = updated.map((item, index) => ({
			...item,
			sort_order: index
		}));

		media = reordered;
		toast.success(`Moved "${moved.title}" to position ${toIndex + 1}`);
		onsort?.(reordered);
	}

	function handleRemoveClick(e: MouseEvent, id: string) {
		e.stopPropagation();
		e.preventDefault();
		onremove?.(id);
		console.log('remove clicked', id);
	}
</script>

<div class="mt-10 space-y-4">
	<div class="flex items-center justify-between">
		<h2 class="text-lg font-semibold text-gray-900">Linked Media</h2>

		{#if onlink}
			<button
				onclick={onlink}
				class="inline-flex items-center gap-1 rounded-md border border-gray-300 bg-white px-3 py-1.5 text-sm text-gray-700 shadow-sm hover:bg-gray-50 focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 focus:outline-hidden"
			>
				<Link class="size-4" />
				<span>Link Media</span>
			</button>
		{/if}
	</div>

	{#if media.length === 0}
		<p class="text-sm text-gray-500">No media linked to this project yet.</p>
	{:else}
		<ul class="grid grid-cols-2 gap-x-4 gap-y-8 sm:grid-cols-3 md:grid-cols-4 xl:grid-cols-5">
			{#each media as item, index (item.id)}
				<li
					use:droppable={{ container: index.toString(), callbacks: { onDrop: handleDrop } }}
					class="group relative block rounded-lg bg-gray-100 ring-1 ring-gray-200"
					animate:flip={{ duration: 300 }}
				>
					<!-- ✅ Move draggable here, so it's only the image that's draggable -->
					<div use:draggable={{ container: index.toString(), dragData: item }} class="cursor-grab">
						<img
							src={item.thumbnail_url || item.url}
							alt={item.alt_text || item.title || 'Media'}
							class="pointer-events-none aspect-video w-full object-cover transition group-hover:opacity-75"
						/>

						<div class="p-2">
							<p class="truncate text-sm font-medium text-gray-900">
								{item.title || 'Untitled'}
							</p>
							<p class="truncate text-xs text-gray-500">
								{item.alt_text || 'No alt text'}
							</p>
						</div>
					</div>

					<!-- ✅ Properly working remove button -->
					{#if onremove}
						<button
							onclick={(e) => handleRemoveClick(e, item.id)}
							class="absolute top-1 right-1 inline-flex cursor-pointer items-center rounded-full bg-white p-1 text-gray-500 shadow hover:bg-gray-100 hover:text-red-600"
						>
							<Unlink class="size-4" />
							<span class="sr-only">Remove</span>
						</button>
					{/if}
				</li>
			{/each}
		</ul>
	{/if}
</div>
