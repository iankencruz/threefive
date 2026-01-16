<script lang="ts">
	import type { Project } from '$types/projects.js';
	import { goto } from '$app/navigation';

	let { data } = $props();
	let projects = $derived(data.projects);
	let pagination = $derived(data.pagination);
</script>

<svelte:head>
	<title>Projects</title>
</svelte:head>

<div class="mx-auto w-full py-8">
	<h1 class="mb-8 text-3xl font-bold">Projects</h1>

	<!-- Projects Grid -->
	<div class="grid grid-cols-1 gap-2 md:grid-cols-2 lg:grid-cols-3">
		{#each projects as project}
			<a
				href="/projects/{project.slug}"
				class="group relative block h-full w-auto overflow-hidden border transition-transform hover:scale-[1.02]"
			>
				{#if project.featured_image}
					<img
						src={project.featured_image.medium_url}
						alt={project.title}
						class="h-full w-full object-cover transition-all duration-300 group-hover:brightness-75"
					/>

					<!-- Title Overlay on Hover -->
					<div
						class="absolute inset-0 flex items-center justify-center bg-black/50 opacity-0 transition-opacity duration-300 group-hover:opacity-100"
					>
						<div class="px-6 text-center">
							<h2 class="mb-2 text-2xl font-bold text-white">
								{project.title}
							</h2>
							{#if project.description}
								<p class="line-clamp-2 text-sm text-white/90">
									{project.description}
								</p>
							{/if}
							{#if project.technologies && project.technologies.length > 0}
								<div class="mt-3 flex flex-wrap justify-center gap-2">
									{#each project.technologies.slice(0, 3) as tech}
										<span
											class="rounded-full bg-white/20 px-2 py-1 text-xs text-white backdrop-blur-sm"
										>
											{tech}
										</span>
									{/each}
								</div>
							{/if}
						</div>
					</div>
				{:else}
					<!-- Fallback if no image -->
					<div
						class="flex h-64 items-center justify-center bg-gradient-to-br from-primary/10 to-primary/5"
					>
						<span class="text-4xl font-bold text-primary/30">{project.title.charAt(0)}</span>
					</div>
				{/if}
			</a>
		{/each}
	</div>

	<!-- Pagination Controls -->
	{#if pagination && pagination.total_pages > 1}
		<div class="mt-8 flex items-center justify-center gap-2">
			<button
				onclick={() => goto(`/projects?page=${pagination.page - 1}`)}
				disabled={pagination.page === 1}
				class="rounded bg-gray-200 px-4 py-2 disabled:opacity-50"
			>
				Previous
			</button>

			<span class="px-4">
				Page {pagination.page} of {pagination.total_pages}
			</span>

			<button
				onclick={() => goto(`/projects?page=${pagination.page + 1}`)}
				disabled={pagination.page === pagination.total_pages}
				class="rounded bg-gray-200 px-4 py-2 disabled:opacity-50"
			>
				Next
			</button>
		</div>

		<p class="mt-4 text-center text-sm text-gray-600">
			Showing {projects.length} of {pagination.total_count} projects
		</p>
	{/if}
</div>

<style>
	.line-clamp-2 {
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
</style>
