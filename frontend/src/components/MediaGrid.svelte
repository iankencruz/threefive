<script lang="ts">
	import MediaItem from './MediaItem.svelte';
	let { media, refresh } = $props();
</script>

<div class="grid grid-cols-2 gap-4 md:grid-cols-3 lg:grid-cols-4">
	{#each media as item (item.id)}
		<MediaItem {item} onrefresh={refresh} />
	{/each}
</div>
