<script lang="ts">
	import { page } from '$app/stores';
	import { getProjectById, updateProjectById } from '$src/lib/api/projects';
	import ProjectMediaGrid from '$src/lib/components/Media/ProjectMediaGrid.svelte';
	import Breadcrumb from '$src/lib/components/Navigation/Breadcrumb.svelte';
	import type { Project } from '$src/lib/types';
	import { slugify, formatDate } from '$src/lib/utils/utilities';
	import { toast } from 'svelte-sonner';

	let loading: boolean = $state(true);
	let project: Project | null = $state(null);

	let editing: boolean = $state(false);
	let saving: boolean = $state(false);

	// editable fields (copy from loaded project)
	let title = $state('');
	let slug = $state('');
	let description = $state('');
	let coverMediaId = $state<string | null>(null);

	let error = $state<string | null>(null);

	$effect(() => {
		// kick off async logic manually
		(async () => {
			loading = true;
			const id = $page.params.id;

			try {
				project = await getProjectById(id);
			} catch (err) {
				console.error('Failed to load project', err);
				error = 'Project not found or failed to load';
				toast.error(error);
				project = null;
			} finally {
				loading = false;
			}
		})();
	});

	$effect(() => {
		if (project) {
			title = project.title;
			slug = project.slug;
			description = project.description ?? '';
			coverMediaId = project.cover_media_id;
		}
	});

	$effect(() => {
		slug = slugify(title);
	});

	async function saveChanges() {
		saving = true;
		try {
			project = await updateProjectById($page.params.id, {
				title,
				slug,
				description,
				cover_media_id: coverMediaId
			});
			toast.success('Project updated successfully');
			editing = false;
		} catch (err) {
			console.error('❌ Save failed:', err);
			toast.error('Failed to save changes. Please try again.');
		} finally {
			saving = false;
		}
	}
</script>

<section class="mt-12">
	{#if loading}
		<p class="text-gray-500">Loading project...</p>
	{:else if error}
		<p class="text-red-500">{error}</p>
	{:else if project}
		<div class="mb-4 flex items-start justify-between">
			<h1 class="text-2xl font-bold">Project Details</h1>

			<button onclick={() => (editing = !editing)} class="text-sm text-indigo-600 hover:underline">
				{editing ? 'Cancel' : 'Edit'}
			</button>
		</div>

		{#if editing}
			<div class="mt-4 space-y-4 border-t border-gray-200 py-4">
				<label for="title" class="text-sm text-gray-600"> Title:</label>
				<input
					type="text"
					name="title"
					bind:value={title}
					class="w-full rounded border border-gray-300 px-3 py-2 text-lg"
					placeholder="Project title"
				/>

				<label for="slug" class="text-sm text-gray-600">Slug (auto-generated):</label>
				<input
					type="text"
					name="slug"
					bind:value={slug}
					disabled
					class="w-full rounded border border-gray-300 px-3 py-2 text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 disabled:cursor-not-allowed disabled:bg-gray-50 disabled:text-gray-500 disabled:outline-gray-200 sm:text-sm/6"
				/>

				<label for="description" class="text-sm text-gray-600">Description</label>
				<textarea
					bind:value={description}
					name="description"
					class="w-full rounded border border-gray-300 px-3 py-2 text-sm"
					rows="5"
					placeholder="Project description"
				>
				</textarea>

				<!-- Media linking placeholder -->
				<div class="rounded border bg-gray-50 p-4 text-sm">
					<p class="text-gray-700">Cover Media ID: {coverMediaId || 'None'}</p>
					<button class="mt-2 text-xs text-indigo-600 hover:underline">
						Link/Change Cover Media
					</button>
				</div>

				<button
					disabled={saving}
					onclick={saveChanges}
					class="mt-4 rounded bg-indigo-600 px-4 py-2 text-white hover:bg-indigo-700 disabled:opacity-50"
				>
					{saving ? 'Saving...' : 'Save Changes'}
				</button>
			</div>
		{:else}
			<div class="mt-4 w-full space-y-4 border-t border-gray-200 py-4">
				<div class="grid grid-cols-2 justify-between space-y-2 text-sm text-gray-400">
					<div>
						<h2 id="applicant-information-title" class="text-lg/6 font-medium text-gray-900">
							{project.title}
						</h2>
						<p class="text-sm text-gray-500">Slug: {project.slug}</p>
					</div>
					<div class="justify-self-end text-right">
						<p>Created: {formatDate(project.created_at, 'relative')}</p>
						<p>Last Updated: {formatDate(project.updated_at, 'relative')}</p>
					</div>
				</div>

				<div class="prose prose-sm text-gray-700">
					<p>{project.description}</p>
				</div>
			</div>
		{/if}
	{/if}

	{#if project && project.media && project.media.length > 0}
		<ProjectMediaGrid media={project.media} />
	{/if}
</section>
