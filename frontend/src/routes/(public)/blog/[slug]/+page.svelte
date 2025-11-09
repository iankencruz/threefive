<!-- routes/(public)/blog/[slug]/+page.svelte -->
<script lang="ts">
	import { getContext, onMount } from "svelte";
	import BlockRenderer from "$lib/components/blocks/BlockRenderer.svelte";
	import { getNavbarVariant } from "$lib/utils/navbar";
	import type { PageData } from "./$types";

	let { data }: { data: PageData } = $props();

	// Get navbar context
	const navbar = getContext<{
		variant: string;
		setVariant: (v: "transparent" | "opaque") => void;
	}>("navbar");

	// Set navbar variant based on blocks
	onMount(() => {
		const variant = getNavbarVariant(data.page.blocks || []);
		navbar.setVariant(variant);
	});
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
	<!-- Render blocks (including hero if present) -->
	<BlockRenderer blocks={data.page.blocks || []} mediaMap={data.mediaMap || {}} />
	
	<!-- Blog-specific metadata section (optional) -->
	{#if data.page.blog_data}
		<section class="py-8 bg-gray-50 border-t border-gray-200">
			<div class="container mx-auto px-4 max-w-4xl">
				<div class="flex items-center justify-between text-sm text-gray-600">
					{#if data.page.blog_data.reading_time}
						<span>{data.page.blog_data.reading_time} min read</span>
					{/if}
					<span>{new Date(data.page.created_at).toLocaleDateString()}</span>
				</div>
			</div>
		</section>
	{/if}
</div>
