<script lang="ts">
	import BlockRenderer from '$lib/components/Public/BlockRenderer.svelte';
	let { page } = $props();
</script>

<svelte:head>
	<title>{page.title}</title>
	<meta name="description" content={page.seo_description} />
</svelte:head>

<article class="prose mx-auto w-full max-w-7xl px-4 py-12">
	<BlockRenderer blocks={page.content} />
</article>
