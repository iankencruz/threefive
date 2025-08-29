<script lang="ts">
	import { page } from '$app/state';
	import type { Project, MediaItem } from '$lib/types';

	import LinkedMediaGrid from '$lib/components/Media/LinkedMediaGrid.svelte';
	import LinkMediaModal from '$lib/components/Media/LinkMediaModal.svelte';
	import { slugify, formatDate } from '$lib/utils/utilities';
	import { toast } from 'svelte-sonner';
	import { goto } from '$app/navigation';
	import {
		getGalleryBySlug,
		linkMediaToGallery,
		unlinkMediaFromGallery,
		updateGalleryBySlug
	} from '$lib/api/galleries';

	let loading = $state(true);
	let gallery = $state<Project | null>(null);
	let galleryMedia = $state<MediaItem[]>([]);

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
		if (!gallery) return;
		try {
			const refreshed = await getGalleryBySlug(gallery.slug);
			gallery = refreshed;
			galleryMedia = [...(refreshed.media ?? [])];
			linkedMediaIds = galleryMedia.map((m) => m.id);
		} catch (err) {
			console.error('❌ Failed to refresh media:', err);
			toast.error('Could not refresh media list');
		}
	}

	// Called when a new media item is linked via modal
	function handleMediaLinked(item: MediaItem) {
		if (!gallery) return;
		gallery = { ...gallery, media: [...(gallery.media ?? []), item] };
		galleryMedia = [...(gallery.media ?? [])];
		linkedMediaIds = galleryMedia.map((m) => m.id);
	}

	$effect(() => {
		(async () => {
			loading = true;
			try {
				const slugParam = page.params.slug;
				const fetched = await getGalleryBySlug(slugParam);
				gallery = fetched;
				galleryMedia = [...(fetched.media ?? [])];
				linkedMediaIds = galleryMedia.map((m) => m.id);
			} catch (err) {
				console.error('Failed to load project', err);
				error = 'Project not found or failed to load';
				gallery = null;
			} finally {
				loading = false;
			}
		})();
	});

	$effect(() => {
		if (gallery) {
			title = gallery.title;
			slug = gallery.slug;
			description = gallery.description ?? '';
			if (gallery.cover_media_id && gallery.media?.some((m) => m.id === gallery?.cover_media_id)) {
				coverMediaId = gallery.cover_media_id;
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
			const updated = await updateGalleryBySlug(page.params.slug, {
				title,
				slug,
				description
			});
			toast.success('Project updated');
			editing = false;

			if (updated.slug !== page.params.slug) {
				await goto(`/admin/galleries/${updated.slug}`);
			} else {
				gallery = await getProjectBySlug(updated.slug);
				galleryMedia = [...(gallery?.media ?? [])];
				linkedMediaIds = galleryMedia.map((m) => m.id);
			}
		} catch (err) {
			console.error('Save failed', err);
			toast.error('Failed to save project');
		} finally {
			saving = false;
		}
	}

	async function handleLinkMedia(media: MediaItem) {
		if (!gallery) return;

		try {
			await linkMediaToGallery(gallery.slug, media.id); // extract ID here
			toast.success('Media linked');

			// refresh full project
			const refreshed = await getGalleryBySlug(gallery.slug);
			gallery = refreshed;
			galleryMedia = [...(refreshed.media ?? [])];
			linkedMediaIds = galleryMedia.map((m) => m.id);
		} catch (err) {
			console.error('Link failed', err);
			toast.error('Failed to link media');
		}
	}

	async function handleUnlinkMedia(mediaId: string) {
		if (!gallery) return;
		try {
			await unlinkMediaFromGallery(gallery.slug, mediaId);
			toast.success('Media unlinked');

			if (gallery.cover_media_id === mediaId) {
				coverMediaId = null;
			}

			gallery = await getProjectBySlug(gallery.slug);
			galleryMedia = [...(gallery?.media ?? [])];
			linkedMediaIds = galleryMedia.map((m) => m.id);
			await saveChanges();
		} catch (err) {
			console.error('Unlink failed', err);
			toast.error('Failed to unlink media');
		}
	}

	async function handleSortMedia(updated: MediaItem[]) {
		if (!gallery) return;
		try {
			const ids = updated.map((m) => m.id);
			await updateProjectMediaOrder(gallery.slug, ids);
			gallery = { ...gallery, media: updated };
			galleryMedia = [...updated];
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
	{:else if gallery}
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
						<h2 class="text-lg font-medium text-gray-900">Title: {gallery.title}</h2>
						<p class="mt-1 text-sm text-gray-500">Slug: {gallery.slug}</p>
					</div>
					<div class="justify-self-end text-right">
						<p>Created: {formatDate(gallery.created_at, 'relative')}</p>
						<p>Updated: {formatDate(gallery.updated_at, 'relative')}</p>
					</div>
				</div>

				<p class="prose prose-sm text-gray-700">{gallery.description}</p>

				<div class="mt-4">
					<p class="mb-1 text-sm text-gray-500">Cover Image:</p>
					{#if gallery.cover_media_id}
						{#if galleryMedia.some((m) => m.id === gallery?.cover_media_id)}
							{#each galleryMedia as media (media.id)}
								{#if media.id === gallery.cover_media_id}
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
			media={galleryMedia}
			onremove={handleUnlinkMedia}
			type="galleries"
			slug={gallery.slug}
			onsort={handleSortMedia}
			onrefresh={refreshMediaGrid}
		/>

		<LinkMediaModal
			open={modalOpen}
			context={{ type: 'gallery', id: gallery.slug }}
			onclose={handleModalClose}
			onlinked={handleLinkMedia}
			{linkedMediaIds}
		/>

		<LinkMediaModal
			open={coverModalOpen}
			context={{ type: 'gallery', id: gallery.slug }}
			onclose={() => (coverModalOpen = false)}
			onlinked={handleCoverMediaSelected}
			selectOnly={true}
			mediaPool={galleryMedia}
			{linkedMediaIds}
		/>
	{/if}
</section>
