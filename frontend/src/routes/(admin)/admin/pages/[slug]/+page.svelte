<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';
	import type { Gallery, MediaItem, Page } from '$lib/types';
	import { deletePage, updatePage } from '$lib/api/pages';
	import { getMediaById } from '$lib/api/media';
	import { slugify } from '$lib/utils/utilities';
	import TipTap from '$lib/components/Builders/TipTap.svelte';
	import LinkGalleryModal from '$lib/components/Galleries/LinkGalleryModal.svelte';

	let data = $state<Page | null>(null);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let saving = $state(false);

	let coverMedia = $state<MediaItem | null>(null);
	let bannerModalOpen = $state(false);

	// ✅ state for galleries
	let linkedGalleries = $state<Gallery[]>([]);
	let linkedGalleryIds = $state<string[]>([]);
	let galleryModalOpen = $state(false);

	let body = $state('');

	// --- Derived slug (only thing our effects should watch) ---
	const slug = $derived(page.params.slug);

	// --- Fetch functions ---
	async function fetchPage(slug: string, signal: AbortSignal) {
		try {
			loading = true;
			error = null;
			data = null;

			const res = await fetch(`/api/v1/admin/pages/${slug}`, { signal });
			if (!res.ok) throw new Error(`Failed to load page (${res.status})`);
			const json = await res.json();
			data = json.data as Page;
			body = data.content ?? '';
		} catch (e) {
			if ((e as any).name !== 'AbortError') {
				console.error(e);
				error = (e as Error).message ?? 'Failed to load page';
			}
		} finally {
			loading = false;
		}
	}

	async function fetchLinkedGalleries(slug: string, signal?: AbortSignal) {
		try {
			const res = await fetch(`/api/v1/admin/pages/${slug}/galleries`, { signal });
			if (!res.ok) {
				console.error('fetchLinkedGalleries failed:', res.status);
				return;
			}
			const json = await res.json();
			linkedGalleries = Array.isArray(json.data) ? [...(json.data as Gallery[])] : [];
			linkedGalleryIds = linkedGalleries.map((g) => g.id);
		} catch (err) {
			if ((err as any).name === 'AbortError') return;
			console.error('fetchLinkedGalleries error:', err);
			toast.error('Failed to fetch galleries');
			linkedGalleries = [];
			linkedGalleryIds = [];
		}
	}

	// --- Effects ---

	// ✅ Effect 1: fetch page + galleries together on slug change
	$effect(() => {
		if (!slug) return;

		const ctrl = new AbortController();

		(async () => {
			try {
				await fetchPage(slug, ctrl.signal);
				await fetchLinkedGalleries(slug, ctrl.signal);
			} catch (e) {
				if ((e as any).name !== 'AbortError') {
					console.error('Failed to fetch page or galleries:', e);
				}
			}
		})();

		return () => ctrl.abort();
	});

	// ✅ Effect 2: fetch cover media once when cover_image_id changes
	$effect(() => {
		const coverId = data?.cover_image_id;
		if (!coverId) return;

		(async () => {
			try {
				const media = await getMediaById(coverId);
				if (!coverMedia || coverMedia.id !== media.id) {
					coverMedia = media;
				}
			} catch (e) {
				console.error('Failed to fetch banner image:', e);
				coverMedia = null;
				toast.error('Failed to load banner image. Please try again.');
			}
		})();
	});

	// --- Handlers ---
	async function handleUpdate(next: Page): Promise<void> {
		try {
			saving = true;
			next.content = body;
			await updatePage(next, slug);
			toast.success('Page updated');
			await goto('/admin/pages');
		} catch (e) {
			console.error(e);
			toast.error('Update failed');
		} finally {
			saving = false;
		}
	}

	async function handleDelete(next: Page): Promise<void> {
		if (!confirm('Are you sure you want to delete this page?')) return;

		try {
			await deletePage(next.slug);
			toast.success('Page deleted');
			await goto('/admin/pages');
		} catch (err) {
			console.error(err);
			toast.error('Delete failed');
		}
	}

	function updateSlug(): void {
		if (!data) return;
		data.slug = slugify(data.title ?? '');
	}

	function handleBannerSelected(item: MediaItem) {
		if (!data) return;
		coverMedia = item;
		data.cover_image_id = item.id;
		bannerModalOpen = false;
	}

	// ✅ auto-refresh galleries after linking
	async function handleLinkGallery(gallery: Gallery) {
		try {
			const res = await fetch(`/api/v1/admin/pages/${slug}/galleries`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ gallery_id: gallery.id })
			});
			if (!res.ok) throw new Error('Failed to link gallery');
			toast.success('Gallery linked');
			await fetchLinkedGalleries(slug);
		} catch (err) {
			console.error(err);
			toast.error('Link failed');
		}
	}

	async function handleUnlinkGallery(galleryId: string) {
		try {
			const res = await fetch(`/api/v1/admin/pages/${slug}/galleries/${galleryId}`, {
				method: 'DELETE'
			});
			if (!res.ok) throw new Error('Failed to unlink gallery');
			toast.success('Gallery unlinked');
			await fetchLinkedGalleries(slug);
		} catch (err) {
			console.error(err);
			toast.error('Unlink failed');
		}
	}
</script>

{#if loading}
	<p>Fetching…</p>
{:else if error}
	<p class="text-red-600">Something went wrong: {error}</p>
{:else if data}
	<div class="flex max-w-2xl flex-col space-y-6">
		<!-- Title & Slug -->
		<div class="flex w-full flex-row justify-between gap-8">
			<div class="grid grow grid-cols-1">
				<div>
					<label for="title" class="block text-sm font-medium text-gray-700">Title</label>
					<input
						name="title"
						type="text"
						class="mt-1 w-full rounded border px-3 py-2"
						bind:value={data.title}
						oninput={updateSlug}
					/>
				</div>
				<div>
					<label for="slug" class="block text-sm font-medium text-gray-700">Slug</label>
					<input
						name="slug"
						type="text"
						class="mt-1 w-full rounded border bg-gray-100 px-3 py-2"
						bind:value={data.slug}
						readonly
					/>
				</div>
			</div>
			<!-- Banner Image -->
			<div>
				<p class="text-sm font-medium text-gray-700">Banner Image</p>
				{#if coverMedia}
					<div class="relative w-max">
						<img
							src={coverMedia.thumbnail_url || coverMedia.url}
							alt={coverMedia.alt_text || coverMedia.title || 'Banner Image'}
							class="max-h-40 rounded-md ring-1 ring-gray-200"
						/>
						<button
							class="absolute top-1 right-1 rounded-full bg-white p-1 px-2 text-xs text-gray-600 shadow hover:bg-gray-100"
							onclick={() => {
								if (!data) return;
								data.cover_image_id = null;
								coverMedia = null;
							}}
						>
							clear
						</button>
					</div>
				{:else}
					<p class="text-sm text-gray-400">No banner image selected</p>
					<button
						type="button"
						onclick={() => (bannerModalOpen = true)}
						class="mt-2 rounded border border-gray-300 bg-white px-3 py-1.5 text-sm shadow-sm hover:bg-gray-50"
					>
						{coverMedia ? 'Change' : 'Select'} Banner Image
					</button>
				{/if}
			</div>
		</div>

		<!-- Content -->
		<div>
			<label for="content" class="block text-sm font-medium text-gray-700">Content</label>
			<TipTap bind:body />
		</div>

		<!-- Attach Galleries -->
		<div class="w-full rounded-md border p-2">
			<div class="flex flex-row items-center justify-between">
				<p>Galleries</p>
				<button
					type="button"
					onclick={() => (galleryModalOpen = true)}
					class="rounded border border-gray-300 bg-white px-3 py-1.5 text-sm shadow-sm hover:bg-gray-50"
				>
					Link Gallery
				</button>
			</div>
			{#if linkedGalleries.length > 0}
				<ul class="mt-2 divide-y divide-gray-200">
					{#each linkedGalleries as g}
						<li class="flex items-center justify-between py-2">
							<div>
								<p class="font-medium">{g.title}</p>
								<p class="text-sm text-gray-500">{g.slug}</p>
							</div>
							<button
								class="rounded bg-red-50 px-2 py-1 text-sm text-red-700 hover:bg-red-100"
								onclick={() => handleUnlinkGallery(g.id)}
							>
								Unlink
							</button>
						</li>
					{/each}
				</ul>
			{:else}
				<p class="mt-2 text-sm text-gray-500">No galleries linked</p>
			{/if}
		</div>

		<!-- Save Controls -->
		<div class="flex items-center justify-between">
			<div class="flex gap-x-4">
				<button
					onclick={() => handleDelete(data!)}
					class="rounded bg-red-600 px-4 py-2 text-white hover:bg-gray-800"
				>
					Delete
				</button>
				<button
					onclick={() => handleUpdate(data!)}
					class="rounded bg-black px-4 py-2 text-white hover:bg-gray-800"
				>
					Save Page
				</button>
			</div>
		</div>
	</div>
{/if}

<LinkGalleryModal
	open={galleryModalOpen}
	onclose={() => (galleryModalOpen = false)}
	pageSlug={data!.slug}
	{linkedGalleryIds}
	onlinked={(g) => handleLinkGallery(g)}
/>
