<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchMedia, deleteMedia, updateMedia } from '$lib/api/media';
	import MediaGrid from '$src/components/MediaGrid.svelte';
	import MediaList from '$src/components/MediaList.svelte';

	let view = $state('grid'); // 'list' or 'grid'
	let media = $state([]);
	let page = $state(1);
	let totalPages = $state(1);

	onMount(async () => {
		await loadMedia();
	});

	async function loadMedia() {
		const res = await fetchMedia(page);
		media = res.items;
		totalPages = res.total_pages;
	}
</script>

<div class="mb-4 flex items-center justify-between">
	<h1 class="text-xl font-bold">Media Library</h1>
	<div class="flex gap-2">
		<button onclick={() => (view = 'grid')}>🔲 Grid</button>
		<button onclick={() => (view = 'list')}>📋 List</button>
	</div>
</div>

{#if view === 'grid'}
	<MediaGrid {media} refresh={loadMedia} />
{:else}
	<MediaList {media} refresh={loadMedia} />
{/if}
