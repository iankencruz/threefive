<script lang="ts">
	import MediaItem from './MediaItem.svelte';
	import { Link, Link2Off, Unlink } from '@lucide/svelte';

	let { media, onremove, onrefresh, onlink } = $props<{
		media: any[];
		onremove?: (id: string) => void;
		onrefresh?: () => void;
		onlink?: () => void;
	}>();
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
			{#each media as item (item.id)}
				<li class="group relative block rounded-lg bg-gray-100 ring-1 ring-gray-200">
					<img
						src={item.thumbnail_url || item.url}
						alt={item.alt_text || item.title || 'Media'}
						class="pointer-events-none aspect-video w-full object-cover transition group-hover:opacity-75"
					/>

					<div class="p-2">
						<p class="truncate text-sm font-medium text-gray-900">{item.title || 'Untitled'}</p>
						<p class="truncate text-xs text-gray-500">{item.alt_text || 'No alt text'}</p>
					</div>

					<!-- Remove button -->
					{#if onremove}
						<button
							onclick={() => onremove?.(item.id)}
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
