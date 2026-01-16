<!-- frontend/src/routes/[slug]/+page.svelte -->
<script lang="ts">
	import BlockRenderer from '$lib/components/blocks/BlockRenderer.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	// Format date helper
	const formatDate = (dateString: string) => {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'long',
			day: 'numeric'
		});
	};
</script>

<svelte:head>
	<!-- SEO Meta Tags -->
	<title>TFP - {data.page.title || data.page.seo?.meta_title}</title>
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
		<meta
			name="robots"
			content="{data.page.seo.robots_index ? 'index' : 'noindex'}, {data.page.seo.robots_follow
				? 'follow'
				: 'nofollow'}"
		/>
	{/if}

	<!-- Canonical URL -->
	{#if data.page.seo?.canonical_url}
		<link rel="canonical" href={data.page.seo.canonical_url} />
	{/if}
</svelte:head>

<!-- Page Content -->
<div class="min-h-screen bg-white">
	<!-- ✨ Pass mediaMap to BlockRenderer -->
	<BlockRenderer blocks={data.page.blocks || []} mediaMap={data.mediaMap || {}} />

	<!-- Optional: Project-specific footer section -->
	{#if data.page.page_type === 'project' && data.page.project_data}
		<section class="border-t border-gray-200 bg-gray-50 py-16">
			<div class="container mx-auto max-w-4xl px-4">
				<h2 class="mb-8 text-center text-2xl font-bold text-gray-900">Project Details</h2>

				<div class="grid grid-cols-1 gap-6 md:grid-cols-2">
					{#if data.page.project_data.client_name}
						<div class="rounded-lg bg-white p-6 shadow-sm">
							<h3 class="mb-2 text-sm font-medium text-gray-500">Client</h3>
							<p class="text-lg font-semibold text-gray-900">
								{data.page.project_data.client_name}
							</p>
						</div>
					{/if}

					{#if data.page.project_data.project_year}
						<div class="rounded-lg bg-white p-6 shadow-sm">
							<h3 class="mb-2 text-sm font-medium text-gray-500">Year</h3>
							<p class="text-lg font-semibold text-gray-900">
								{data.page.project_data.project_year}
							</p>
						</div>
					{/if}

					{#if data.page.project_data.project_url}
						<div class="rounded-lg bg-white p-6 shadow-sm">
							<h3 class="mb-2 text-sm font-medium text-gray-500">Live Site</h3>
							<a
								href={data.page.project_data.project_url}
								target="_blank"
								rel="noopener noreferrer"
								class="text-lg font-semibold text-blue-600 hover:text-blue-700"
							>
								View Project →
							</a>
						</div>
					{/if}

					{#if data.page.project_data.project_status}
						<div class="rounded-lg bg-white p-6 shadow-sm">
							<h3 class="mb-2 text-sm font-medium text-gray-500">Status</h3>
							<p class="text-lg font-semibold text-gray-900 capitalize">
								{data.page.project_data.project_status}
							</p>
						</div>
					{/if}
				</div>

				{#if data.page.project_data.technologies && data.page.project_data.technologies.length > 0}
					<div class="mt-8 rounded-lg bg-white p-6 shadow-sm">
						<h3 class="mb-4 text-sm font-medium text-gray-500">Technologies Used</h3>
						<div class="flex flex-wrap gap-2">
							{#each data.page.project_data.technologies as tech}
								<span class="rounded-full bg-blue-100 px-3 py-1 text-sm font-medium text-blue-700">
									{tech}
								</span>
							{/each}
						</div>
					</div>
				{/if}
			</div>
		</section>
	{/if}

	<!-- Blog metadata -->
	{#if data.page.page_type === 'blog' && data.page.blog_data}
		<section class="border-t border-gray-200 py-8">
			<div class="container mx-auto max-w-4xl px-4">
				<div class="flex items-center justify-between text-sm text-gray-600">
					{#if data.page.blog_data.reading_time}
						<span>{data.page.blog_data.reading_time} min read</span>
					{/if}
					<span>Published {formatDate(data.page.created_at)}</span>
				</div>
			</div>
		</section>
	{/if}
</div>
