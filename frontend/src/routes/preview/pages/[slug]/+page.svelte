<!-- frontend/src/routes/preview/pages/[slug]/+page.svelte -->
<script lang="ts">
	import BlockRenderer from "$lib/components/blocks/BlockRenderer.svelte";
	import { Eye } from "lucide-svelte";

	// Type assertion to include mediaMap
	interface PreviewPageData {
		page: any;
		mediaMap: Record<string, any>;
		isPreview: boolean;
	}

	let { data } = $props<{ data: PreviewPageData }>();

	// Format date helper
	const formatDate = (dateString: string) => {
		return new Date(dateString).toLocaleDateString("en-US", {
			year: "numeric",
			month: "long",
			day: "numeric",
		});
	};

	// Status badge colors matching your image
	const statusColors = {
		draft: {
			bg: "bg-yellow-100",
			text: "text-yellow-800",
			banner: "bg-yellow-500 text-yellow-900",
		},
		published: {
			bg: "bg-green-100",
			text: "text-green-800",
			banner: "bg-green-500 text-green-900",
		},
	};

	const currentStatus = data.page.status || "draft";
	const colors =
		statusColors[currentStatus as keyof typeof statusColors] ||
		statusColors.draft;
</script>

<svelte:head>
	<title>Preview: {data.page.title}</title>
	<meta name="robots" content="noindex, nofollow" />
</svelte:head>

<!-- Preview Banner with dynamic colors -->
{#if data.isPreview}
  {	console.log("data: ", data)}
	<div class="{colors.banner} px-4 py-3 text-center font-medium sticky top-0 z-50 shadow-md">
		<div class="flex items-center justify-center gap-3">
			<Eye size={16}/>
			<span>Preview Mode</span>
			<span class="{colors.bg} {colors.text} px-3 py-1 rounded-full text-sm font-semibold uppercase">
				{data.page.status}
			</span>
		</div>
	</div>
{/if}

<!-- Page Content -->
<div class="min-h-screen bg-white">
	<!-- ✨ Pass mediaMap to BlockRenderer -->
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
							<h3 class="text-sm font-medium text-gray-500 mb-2">Live Site</h3>
							<a href={data.page.project_data.project_url} target="_blank" rel="noopener noreferrer" class="text-lg font-semibold text-blue-600 hover:text-blue-700">
								View Project →
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
								<span class="px-3 py-1 bg-blue-100 text-blue-700 rounded-full text-sm font-medium">
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
		<section class="py-8 border-t border-gray-200">
			<div class="container mx-auto px-4 max-w-4xl">
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
