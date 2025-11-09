<!-- frontend/src/routes/preview/pages/[id]/+page.svelte -->
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
<div class="w-full mx-auto relative top-0 ">
<div class="@container w-full absolute top-0 z-50 bg-stone-400/50 pointer-events-none text-black px-4 py-1 shadow-lg">
	<div class="max-w-7xl mx-auto flex items-center justify-between">
		<div class="flex flex-row items-center gap-3">
			<Eye class="hidden lg:block text-white" size={16}/>
			<div class="flex flex-row gap-2">
				<span class="hidden lg:block font-semibold">Preview Mode</span>
				<span class="hidden lg:block mx-2">•</span>
        <div class="inline-flex items-center">
          <strong class={["capitalize text-base font-normal", data.page.status === "published" ? "text-green-600" : "text-amber-800"]} >{data.page.status}</strong>
        </div>
			</div>
		</div>
		
		<div class="flex items-center gap-3">
			<button
				onclick={() => goto(`/admin/pages/${data.page.id}/edit`)}
				class="px-4 py-1 cursor-pointer text-black rounded-sm font-medium hover:bg-gray-100 transition-colors text-sm"
			>
				Edit
			</button>
			<button
				onclick={() => goto('/admin/pages')}
				class="px-4 py-1 cursor-pointer text-black rounded-sm font-medium hover:bg-gray-100 transition-colors text-sm"
			>
				Close
			</button>
		</div>
	</div>
</div>
</div>

<!-- Page Content - Same as public pages -->
<div class="min-h-screen bg-white">
	<!-- ✨ CRITICAL: Pass mediaMap to BlockRenderer (same as public pages) -->
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
