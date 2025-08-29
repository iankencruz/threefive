<script lang="ts">
	import { page } from '$app/state';
	import {
		getProjectBySlug,
		updateProjectBySlug,
		linkMediaToProject,
		unlinkMediaFromProject,
		updateProjectMediaOrder
	} from '$lib/api/projects';
	import type { Project, MediaItem } from '$lib/types';
	import { slugify, formatDate } from '$lib/utils/utilities';
	import { toast } from 'svelte-sonner';
	import { goto } from '$app/navigation';

	import LinkedMediaGrid from '$lib/components/Media/LinkedMediaGrid.svelte';
	import LinkMediaModal from '$lib/components/Media/LinkMediaModal.svelte';

	let loading = $state(true);
	let project = $state<Project | null>(null);
	let projectMedia = $state<MediaItem[]>([]);

	let editing = $state(false);
	let saving = $state(false);

	let title = $state('');
	let slug = $state('');
	let description = $state('');
	let coverMediaId = $state<string | null>(null);

	let error = $state<string | null>(null);
	let modalOpen = $state(false);
	let coverModalOpen = $state(false);
	let linkedMediaIds = $state<string[]>([]);

	function handleCoverClick() {
		coverModalOpen = true;
	}

	function handleCoverMediaSelected(item: MediaItem) {
		coverMediaId = item.id;
		coverModalOpen = false;
	}

	function handleLinkClick() {
		modalOpen = true;
	}

	function handleModalClose() {
		modalOpen = false;
	}

	async function refreshMediaGrid() {
		if (!project) return;
		try {
			const refreshed = await getProjectBySlug(project.slug);
			project = refreshed;
			projectMedia = [...(refreshed.media ?? [])];
			linkedMediaIds = projectMedia.map((m) => m.id);
		} catch (err) {
			console.error('❌ Failed to refresh media:', err);
			toast.error('Could not refresh media list');
		}
	}

	// Called when a new media item is linked via modal
	function handleMediaLinked(item: MediaItem) {
		if (!project) return;
		project = { ...project, media: [...(project.media ?? []), item] };
		projectMedia = [...(project.media ?? [])];
		linkedMediaIds = projectMedia.map((m) => m.id);
	}

	$effect(() => {
		(async () => {
			loading = true;
			try {
				const slugParam = page.params.slug;
				const fetched = await getProjectBySlug(slugParam);
				project = fetched;
				projectMedia = [...(fetched.media ?? [])];
				linkedMediaIds = projectMedia.map((m) => m.id);
			} catch (err) {
				console.error('Failed to load project', err);
				error = 'Project not found or failed to load';
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
			if (project.cover_media_id && project.media?.some((m) => m.id === project?.cover_media_id)) {
				coverMediaId = project.cover_media_id;
			} else {
				coverMediaId = null;
			}
		}
	});

	$effect(() => {
		slug = slugify(title);
	});

	async function saveChanges() {
		saving = true;
		try {
			const updated = await updateProjectBySlug(page.params.slug, {
				title,
				slug,
				description,
				cover_media_id: coverMediaId
			});
			toast.success('Project updated');
			editing = false;

			if (updated.slug !== page.params.slug) {
				await goto(`/admin/projects/${updated.slug}`);
			} else {
				project = await getProjectBySlug(updated.slug);
				projectMedia = [...(project?.media ?? [])];
				linkedMediaIds = projectMedia.map((m) => m.id);
			}
		} catch (err) {
			console.error('Save failed', err);
			toast.error('Failed to save project');
		} finally {
			saving = false;
		}
	}

	async function handleLinkMedia(media: MediaItem) {
		if (!project) return;

		try {
			await linkMediaToProject(project.slug, media.id); // extract ID here
			toast.success('Media linked');

			// refresh full project
			const refreshed = await getProjectBySlug(project.slug);
			project = refreshed;
			projectMedia = [...(refreshed.media ?? [])];
			linkedMediaIds = projectMedia.map((m) => m.id);
		} catch (err) {
			console.error('Link failed', err);
			toast.error('Failed to link media');
		}
	}

	async function handleUnlinkMedia(mediaId: string) {
		if (!project) return;
		try {
			await unlinkMediaFromProject(project.slug, mediaId);
			toast.success('Media unlinked');

			if (project.cover_media_id === mediaId) {
				coverMediaId = null;
			}

			project = await getProjectBySlug(project.slug);
			projectMedia = [...(project?.media ?? [])];
			linkedMediaIds = projectMedia.map((m) => m.id);
			await saveChanges();
		} catch (err) {
			console.error('Unlink failed', err);
			toast.error('Failed to unlink media');
		}
	}

	async function handleSortMedia(updated: MediaItem[]) {
		if (!project) return;
		try {
			const ids = updated.map((m) => m.id);
			await updateProjectMediaOrder(project.slug, ids);
			project = { ...project, media: updated };
			projectMedia = [...updated];
		} catch (err) {
			console.error('Sort failed', err);
			toast.error('Failed to update media order');
		}
	}
</script>

<section class="mt-8">
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
				<label class="text-sm text-gray-600">Title</label>
				<input bind:value={title} class="w-full rounded border border-gray-300 px-3 py-2 text-lg" />

				<label class="text-sm text-gray-600">Slug</label>
				<input
					bind:value={slug}
					disabled
					class="w-full rounded border border-gray-300 px-3 py-2 text-gray-900"
				/>

				<label class="text-sm text-gray-600">Description</label>
				<textarea
					bind:value={description}
					class="w-full rounded border border-gray-300 px-3 py-2 text-sm"
					rows="5"
				></textarea>

				<div class="rounded border bg-gray-50 p-4 text-sm">
					<p class="text-gray-700">Cover Media ID: {coverMediaId || 'None'}</p>
					<button class="mt-2 text-xs text-indigo-600 hover:underline" onclick={handleCoverClick}>
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
				<div class="grid grid-cols-2 justify-between text-sm text-gray-400">
					<div>
						<h2 class="text-lg font-medium text-gray-900">Title: {project.title}</h2>
						<p class="mt-1 text-sm text-gray-500">Slug: {project.slug}</p>
					</div>
					<div class="justify-self-end text-right">
						<p>Created: {formatDate(project.created_at, 'relative')}</p>
						<p>Updated: {formatDate(project.updated_at, 'relative')}</p>
					</div>
				</div>

				<p class="prose prose-sm text-gray-700">{project.description}</p>

				<div class="mt-4">
					<p class="mb-1 text-sm text-gray-500">Cover Image:</p>
					{#if project.cover_media_id}
						{#if projectMedia.some((m) => m.id === project?.cover_media_id)}
							{#each projectMedia as media (media.id)}
								{#if media.id === project.cover_media_id}
									<img
										src={media.medium_url || media.url}
										alt={media.title}
										class="aspect-video h-96 w-full rounded object-cover ring-1 ring-gray-200"
									/>
								{/if}
							{/each}
						{:else}
							<div
								class="flex h-96 items-center justify-center rounded bg-gray-100 text-gray-400 ring-1 ring-gray-200"
							>
								Selected cover image is not linked
							</div>
						{/if}
					{:else}
						<div
							class="flex h-96 items-center justify-center rounded bg-gray-100 text-gray-400 ring-1 ring-gray-200"
						>
							No cover image selected
						</div>
					{/if}
				</div>
			</div>
		{/if}

		<LinkedMediaGrid
			media={projectMedia}
			onremove={handleUnlinkMedia}
			slug={project.slug}
			type="projects"
			onsort={handleSortMedia}
			onrefresh={refreshMediaGrid}
		/>

		<LinkMediaModal
			open={modalOpen}
			context={{ type: 'project', id: project.slug }}
			onclose={handleModalClose}
			onlinked={handleLinkMedia}
			{linkedMediaIds}
		/>

		<LinkMediaModal
			open={coverModalOpen}
			context={{ type: 'project', id: project.slug }}
			onclose={() => (coverModalOpen = false)}
			onlinked={handleCoverMediaSelected}
			selectOnly={true}
			mediaPool={projectMedia}
			{linkedMediaIds}
		/>
	{/if}
</section>
