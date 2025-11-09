<!-- routes/(public)/[slug]/+page.svelte -->
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

	// Format date helper
	const formatDate = (dateString: string) => {
		return new Date(dateString).toLocaleDateString("en-US", {
			year: "numeric",
			month: "long",
			day: "numeric",
		});
	};
</script>

<svelte:head>
	<!-- SEO Meta Tags -->
	<title>{data.page.seo?.meta_title || data.page.title}</title>
	<meta name="description" content={data.page.seo?.meta_description || ''} />
	
	<!-- Open Graph -->
	{#if data.page.seo?.og_title}
		<meta property="og:title" content={data.page.seo.og_title} />
	{/if}
	{#if data.page.seo?.og_description}
		<meta property="og:description" content={data.page.seo.og_description} />
	{/if}
	{#if data.page.seo?.og_image_id}
		<meta property="og:image" content={data.page.seo.og_image_id} />
	{/if}
	
	<!-- Robots -->
	{#if data.page.seo}
		<meta name="robots" content="{data.page.seo.robots_index ? 'index' : 'noindex'}, {data.page.seo.robots_follow ? 'follow' : 'nofollow'}" />
	{/if}
	
	<!-- Canonical URL -->
	{#if data.page.seo?.canonical_url}
		<link rel="canonical" href={data.page.seo.canonical_url} />
	{/if}
</svelte:head>

<!-- Page Content -->
<div class="min-h-screen bg-white">
	<BlockRenderer blocks={data.page.blocks || []} mediaMap={data.mediaMap || {}} />
	
	<!-- Optional: Project-specific footer section -->
	{#if data.page.page_type === 'project' && data.page.project_data}
		<section class="py-16 bg-gray-50 border-t border-gray-200">
			<div class="container mx-auto px-4 max-w-4xl">
				<h2 class="text-2xl font-bold text-gray-900 mb-8 text-center">Project Details</h2>
				
				<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
					{#if data.page.project_data.client_name}
						<div class="bg-white p-6 rounded-lg shadow-sm">
							<h3 class="text-sm font-medium text-gray-500 mb-2">Client</h3>
							<p class="text-lg font-semibold text-gray-900">{data.page.project_data.client_name}</p>
						</div>
					{/if}
					
					{#if data.page.project_data.project_year}
						<div class="bg-white p-6 rounded-lg shadow-sm">
							<h3 class="text-sm font-medium text-gray-500 mb-2">Year</h3>
							<p class="text-lg font-semibold text-gray-900">{data.page.project_data.project_year}</p>
						</div>
					{/if}
					
					{#if data.page.project_data.project_url}
						<div class="bg-white p-6 rounded-lg shadow-sm">
							<h3 class="text-sm font-medium text-gray-500 mb-2">Project URL</h3>
							<a 
								href={data.page.project_data.project_url} 
								target="_blank" 
								rel="noopener noreferrer"
								class="text-lg font-semibold text-blue-600 hover:text-blue-800"
							>
								Visit Project →
							</a>
						</div>
					{/if}
					
					{#if data.page.project_data.technologies && data.page.project_data.technologies.length > 0}
						<div class="bg-white p-6 rounded-lg shadow-sm">
							<h3 class="text-sm font-medium text-gray-500 mb-2">Technologies</h3>
							<div class="flex flex-wrap gap-2 mt-2">
								{#each data.page.project_data.technologies as tech}
									<span class="px-3 py-1 bg-gray-100 text-gray-800 text-sm rounded-full">
										{tech}
									</span>
								{/each}
							</div>
						</div>
					{/if}
				</div>
			</div>
		</section>
	{/if}
</div>
