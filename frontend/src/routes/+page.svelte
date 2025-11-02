<!-- frontend/src/routes/+page.svelte -->
<script lang="ts">
	import BlockRenderer from "$lib/components/blocks/BlockRenderer.svelte";
	import type { PageData } from "./$types";

	let { data }: { data: PageData } = $props();
</script>

<svelte:head>
	<title>{data.page.seo?.meta_title || data.page.title}</title>
	<meta name="description" content={data.page.seo?.meta_description || ''} />
	
	{#if data.page.seo?.og_title}
		<meta property="og:title" content={data.page.seo.og_title} />
	{/if}
	{#if data.page.seo?.og_description}
		<meta property="og:description" content={data.page.seo.og_description} />
	{/if}
	
	{#if data.page.seo}
		<meta name="robots" content="{data.page.seo.robots_index ? 'index' : 'noindex'}, {data.page.seo.robots_follow ? 'follow' : 'nofollow'}" />
	{/if}
</svelte:head>

<div class="min-h-screen bg-white">
	<!-- ✨ Pass mediaMap to BlockRenderer -->
	<BlockRenderer blocks={data.page.blocks || []} mediaMap={data.mediaMap || {}} />
</div>
