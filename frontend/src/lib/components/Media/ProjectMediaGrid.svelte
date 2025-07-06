<script lang="ts">
	import type { MediaItem } from '$src/lib/types';
	import { Link, Unlink } from '@lucide/svelte';
	import { draggable, droppable, type DragDropState } from '@thisux/sveltednd';
	import { flip } from 'svelte/animate';
	import { toast } from 'svelte-sonner';
	import LinkMediaModal from './LinkMediaModal.svelte';

	let { projectSlug, media, onremove, onrefresh, onsort } = $props<{
		projectSlug: string;
		media: MediaItem[];
		onremove?: (id: string) => void;
		onrefresh?: () => void;
		onsort?: (items: MediaItem[]) => void;
	}>();

	let modalOpen = $state(false);

	async function handleLink(media: MediaItem) {
		try {
			const res = await fetch(`/api/v1/admin/projects/${projectSlug}/media`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ media_id: media.id, sort_order: 0 })
			});

			if (!res.ok) throw new Error('Failed to link media');
			toast.success(`Linked media "${media.title || 'Untitled'}"`);
			onrefresh?.();
		} catch (err) {
			console.error('❌ Failed to link media:', err);
			toast.error('Failed to link media');
		}
	}

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
	}
</script>

<div class="mt-10 space-y-4">
	<div class="flex items-center justify-between">
		<h2 class="text-lg font-semibold text-gray-900">Linked Media</h2>

		<button
			onclick={() => (modalOpen = true)}
			class="inline-flex items-center gap-1 rounded-md border border-gray-300 bg-white px-3 py-1.5 text-sm text-gray-700 shadow-sm hover:bg-gray-50"
		>
			<Link class="size-4" />
			<span>Link Media</span>
		</button>
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

					<button
						onclick={(e) => handleRemoveClick(e, item.id)}
						class="absolute top-1 right-1 inline-flex cursor-pointer items-center rounded-full bg-white p-1 text-gray-500 shadow hover:bg-gray-100 hover:text-red-600"
					>
						<Unlink class="size-4" />
						<span class="sr-only">Remove</span>
					</button>
				</li>
			{/each}
		</ul>
	{/if}

	<LinkMediaModal
		open={modalOpen}
		context={{ type: 'project', id: projectSlug }}
		linkedMediaIds={media.map((m: MediaItem) => m.id)}
		onclose={() => (modalOpen = false)}
		onlinked={handleLink}
	/>
</div>
