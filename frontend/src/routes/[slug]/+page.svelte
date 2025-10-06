<!-- frontend/src/routes/[slug]/+page.svelte -->
<script lang="ts">
import BlockRenderer from "$lib/components/blocks/BlockRenderer.svelte";
import type { PageData } from "./$types";

let { data }: { data: PageData } = $props();

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
	<!-- Render all blocks -->
	<BlockRenderer blocks={data.page.blocks || []} />
	
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
							<h3 class="text-sm font-medium text-gray-500 mb-2">Website</h3>
							<a 
								href={data.page.project_data.project_url} 
								target="_blank" 
								rel="noopener noreferrer"
								class="text-lg font-semibold text-blue-600 hover:text-blue-700 hover:underline"
							>
								Visit Site →
							</a>
						</div>
					{/if}
					
					{#if data.page.project_data.project_status}
						<div class="bg-white p-6 rounded-lg shadow-sm">
							<h3 class="text-sm font-medium text-gray-500 mb-2">Status</h3>
							<p class="text-lg font-semibold text-gray-900 capitalize">{data.page.project_data.project_status}</p>
						</div>
					{/if}
				</div>
				
				{#if data.page.project_data.technologies && data.page.project_data.technologies.length > 0}
					<div class="mt-8 bg-white p-6 rounded-lg shadow-sm">
						<h3 class="text-sm font-medium text-gray-500 mb-4">Technologies Used</h3>
						<div class="flex flex-wrap gap-2">
							{#each data.page.project_data.technologies as tech}
								<span class="px-3 py-1 bg-blue-100 text-blue-800 rounded-full text-sm font-medium">
									{tech}
								</span>
							{/each}
						</div>
					</div>
				{/if}
			</div>
		</section>
	{/if}
	
	<!-- Optional: Blog-specific footer section -->
	{#if data.page.page_type === 'blog' && data.page.blog_data}
		<section class="py-16 bg-gray-50 border-t border-gray-200">
			<div class="container mx-auto px-4 max-w-4xl">
				<div class="bg-white p-8 rounded-lg shadow-sm">
					{#if data.page.blog_data.excerpt}
						<p class="text-lg text-gray-600 mb-6 italic">"{data.page.blog_data.excerpt}"</p>
					{/if}
					
					<div class="flex flex-wrap gap-6 text-sm text-gray-500">
						{#if data.page.published_at}
							<div class="flex items-center gap-2">
								<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
								</svg>
								<span>Published {formatDate(data.page.published_at)}</span>
							</div>
						{/if}
						
						{#if data.page.blog_data.reading_time}
							<div class="flex items-center gap-2">
								<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
								</svg>
								<span>{data.page.blog_data.reading_time} min read</span>
							</div>
						{/if}
					</div>
				</div>
			</div>
		</section>
	{/if}
</div>
