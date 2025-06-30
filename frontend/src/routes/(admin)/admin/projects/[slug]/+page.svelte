<script lang="ts">
	import { page } from '$app/stores';
	import {
		getProjectBySlug,
		updateProjectBySlug,
		linkMediaToProject,
		unlinkMediaFromProject,
		updateProjectMediaOrder
	} from '$src/lib/api/projects';
	import LinkMediaModal from '$lib/components/Media/LinkMediaModal.svelte';
	import ProjectMediaGrid from '$src/lib/components/Media/ProjectMediaGrid.svelte';
	import Breadcrumb from '$src/lib/components/Navigation/Breadcrumb.svelte';
	import type { Project } from '$src/lib/types';
	import { slugify, formatDate } from '$src/lib/utils/utilities';
	import { toast } from 'svelte-sonner';
	import { goto } from '$app/navigation';

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
	let modalOpen = $state(false);
	let linkedMediaIds = $state<string[]>([]);

	let coverModalOpen = $state(false);

	function handleCoverClick() {
		coverModalOpen = true;
	}

	function handleCoverMediaSelected(item: any) {
		coverMediaId = item.id;
		coverModalOpen = false;
	}

	function handleLinkClick() {
		modalOpen = true;
	}

	function handleModalClose() {
		modalOpen = false;
	}

	function handleMediaLinked(item: any) {
		if (!project) return;

		project = {
			...project,
			id: project.id, // explicitly required fields
			title: project.title,
			slug: project.slug,
			description: project.description,
			meta_description: project.meta_description,
			canonical_url: project.canonical_url,
			cover_media_id: project.cover_media_id,
			is_published: project.is_published,
			published_at: project.published_at,
			created_at: project.created_at,
			updated_at: project.updated_at,
			media: [...(project.media ?? []), item]
		};

		linkedMediaIds = (project?.media ?? []).map((m) => m.id);
	}

	$effect(() => {
		(async () => {
			loading = true;
			const slug = $page.params.slug;

			try {
				project = await getProjectBySlug(slug);
				linkedMediaIds = project?.media?.map((m) => m.id) ?? [];
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
			// only set coverMediaId if it still exists in media list
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
			const updated = await updateProjectBySlug($page.params.slug, {
				title,
				slug,
				description,
				cover_media_id: coverMediaId
			});
			toast.success('Project updated successfully');
			editing = false;

			// 🧠 navigate to the new slug if it changed

			if (updated.slug !== $page.params.slug) {
				await goto(`/admin/projects/${updated.slug}`);
			} else {
				project = await getProjectBySlug(updated.slug); // ✅ re-fetch full project with media
			}
		} catch (err) {
			console.error('❌ Save failed:', err);
			toast.error('Failed to save changes. Please try again.');
		} finally {
			saving = false;
		}
	}

	async function handleLinkMedia(mediaId: string) {
		if (!project) return;
		try {
			await linkMediaToProject(project.slug, mediaId);
			toast.success('Media linked!');
			project = await getProjectBySlug(project.slug);
			linkedMediaIds = project?.media?.map((m) => m.id) ?? [];
		} catch (err) {
			console.error('Linking failed', err);
			toast.error('Failed to link media');
		}
	}

	async function handleUnlinkMedia(mediaId: string) {
		if (!project) return;

		try {
			await unlinkMediaFromProject(project.slug, mediaId);

			// ✅ If the cover image was just unlinked, clear it locally
			if (project.cover_media_id === mediaId) {
				coverMediaId = null;
			}

			toast.success('Media unlinked!');
			if (project.cover_media_id === mediaId) {
				coverMediaId = null;
			}
			project = await getProjectBySlug(project.slug);
			linkedMediaIds = project?.media?.map((m) => m.id) ?? [];
			saveChanges();
		} catch (err) {
			console.error('Unlinking failed', err);
			toast.error('Failed to unlink media');
		}
	}

	async function handleSortMedia(updated: any[]) {
		if (!project) return;

		try {
			const ids = updated.map((m) => m.id);
			await updateProjectMediaOrder(project.slug, ids);
			project = { ...project, media: updated };
		} catch (err) {
			console.error('Sorting failed:', err);
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
					class="w-full rounded border border-gray-300 px-3 py-2 text-gray-900"
				/>

				<label for="description" class="text-sm text-gray-600">Description</label>
				<textarea
					bind:value={description}
					name="description"
					class="w-full rounded border border-gray-300 px-3 py-2 text-sm"
					rows="5"
					placeholder="Project description"
				></textarea>

				<!-- {@debug coverMediaId} -->
				<!-- Media linking placeholder -->
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
				<div class="grid grid-cols-2 justify-between space-y-2 text-sm text-gray-400">
					<div>
						<h2 id="applicant-information-title" class="text-lg font-medium text-gray-900">
							Title: {project.title}
						</h2>
						<p class="mt-2 text-sm text-gray-500">Slug: {project.slug}</p>
					</div>
					<div class="justify-self-end text-justify">
						<p>Created: {formatDate(project.created_at, 'relative')}</p>
						<p>Last Updated: {formatDate(project.updated_at, 'relative')}</p>
					</div>
				</div>

				<div class="prose prose-sm text-gray-700">
					<p>{project.description}</p>
				</div>

				<div class="mt-4">
					<p class="mb-1 text-sm text-gray-500">Cover Image:</p>
					{#if project.cover_media_id}
						{#if project.media?.some((m) => m.id === project?.cover_media_id)}
							{#each project.media as media (media.id)}
								{#if media.id === project.cover_media_id}
									<img
										src={media.medium_url || media.url}
										alt={media.title}
										class="aspect-video h-96 w-full rounded object-cover ring-1 ring-gray-200"
									/>
								{/if}
							{/each}
						{:else if project.media?.some((m) => m.id !== project?.cover_media_id)}
							<div
								class="flex aspect-video h-96 w-full items-center justify-center rounded bg-gray-100 text-gray-400 ring-1 ring-gray-200"
							>
								Selected cover image is not linked to this project
							</div>
						{/if}
					{:else}
						<div
							class="flex aspect-video h-96 w-full items-center justify-center rounded bg-gray-100 text-gray-400 ring-1 ring-gray-200"
						>
							No cover image selected
						</div>
					{/if}
				</div>
			</div>
		{/if}

		{#if project.media}
			<ProjectMediaGrid
				media={project.media}
				onremove={handleUnlinkMedia}
				onlink={handleLinkClick}
				onsort={handleSortMedia}
			/>
			<!-- Project Media Modal -->
			<LinkMediaModal
				open={modalOpen}
				projectSlug={project.slug}
				onclose={handleModalClose}
				onlinked={handleMediaLinked}
				{linkedMediaIds}
			/>
		{/if}
		<!-- Cover Media Modal -->
		{#if project?.media}
			<LinkMediaModal
				open={coverModalOpen}
				projectSlug={project.slug}
				onclose={() => (coverModalOpen = false)}
				onlinked={handleCoverMediaSelected}
				{linkedMediaIds}
				selectOnly={true}
				mediaPool={project?.media}
			/>
		{/if}
	{/if}
</section>
