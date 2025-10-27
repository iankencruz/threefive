<!-- frontend/src/routes/admin/pages/[id]/preview/+page.svelte -->
<script lang="ts">
	import BlockRenderer from "$lib/components/blocks/BlockRenderer.svelte";
	import { goto } from "$app/navigation";
	import type { PageData } from "./$types";
	import { Eye } from "lucide-svelte";

	let { data }: { data: PageData } = $props();

	const formatDate = (dateString: string) => {
		return new Date(dateString).toLocaleDateString("en-US", {
			year: "numeric",
			month: "long",
			day: "numeric",
		});
	};
</script>

<svelte:head>
	<title>Preview: {data.page.title}</title>
</svelte:head>

<!-- Preview Banner - Sticky at top -->
<div class="w-full mx-auto relative top-0">
<div class="@container w-full absolute top-0 z-50 bg-stone-400 text-black px-4 py-1 shadow-lg">
	<div class="max-w-7xl mx-auto flex items-center justify-between">
		<div class="flex items-center gap-3">
			<Eye class="text-white" size={16}/>
			<div>
				<span class="font-semibold">Preview Mode</span>
				<span class="mx-2">•</span>
				<span class="text-sm">
					Status: <strong class={["capitalize ml-2", data.page.status === "published" ? "text-green-600" : "text-amber-800"]} >{data.page.status}</strong>
				</span>
			</div>
		</div>
		
		<div class="flex items-center gap-3">
			<button
				onclick={() => goto(`/admin/pages/${data.page.id}/edit`)}
				class="px-4 py-1 cursor-pointer  text-black rounded-sm font-medium hover:bg-gray-100 transition-colors text-sm"
			>
				Edit
			</button>
			<button
				onclick={() => goto('/admin/pages')}
				class="px-4 py-1 cursor-pointer  text-black rounded-sm font-medium hover:bg-gray-100 transition-colors text-sm"
			>
				Close
			</button>
		</div>
	</div>
</div>
</div>
<!-- Page Content - Same as public pages -->
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
