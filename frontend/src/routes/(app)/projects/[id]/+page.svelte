<script lang="ts">
	import { page } from '$app/stores';
	import { getProjectById } from '$lib/api/project';

	let loading = $state(true);
	let project: any = $state(null);

	$effect(() => {
		// kick off async logic manually
		(async () => {
			loading = true;
			const id = $page.params.id;

			try {
				project = await getProjectById(id);
			} catch (err) {
				console.error('Failed to load project', err);
			} finally {
				loading = false;
			}
		})();
	});
</script>

<section class="p-6">
	{#if loading}
		<p class="text-gray-500">Loading project...</p>
	{:else if !project}
		<p class="text-red-500">Project not found.</p>
	{:else}
		<h1 class="text-2xl font-semibold text-gray-900">{project.title}</h1>
		<p class="text-sm text-gray-600">Slug: {project.slug}</p>
		<div class="prose prose-sm mt-4 text-gray-700">
			<p>{project.description}</p>
		</div>
	{/if}
</section>
