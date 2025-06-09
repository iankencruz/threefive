<script lang="ts">
	import MediaItem from './MediaItem.svelte';

	let { media, refresh } = $props();
	console.log('Media List View');
</script>

<div class="flex flex-col gap-4">
	{#each media as item (item.id)}
		<MediaItem {item} onrefresh={refresh} />
	{/each}
</div>
