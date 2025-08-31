<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';
	import PageForm from '$lib/components/Forms/PageForm.svelte';
	import type { MediaItem, Page } from '$lib/types';
	import { deletePage, updatePage } from '$lib/api/pages';
	import { getMediaById } from '$lib/api/media';
	import { slugify } from '$lib/utils/utilities';
	import LinkMediaModal from '$lib/components/Media/LinkMediaModal.svelte';
	import TipTap from '$lib/components/Builders/TipTap.svelte';

	let data = $state<Page | null>(null);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let saving = $state(false);

	let coverMedia = $state<MediaItem | null>(null);
	let bannerModalOpen = $state(false);

	async function fetchPage(slug: string, signal: AbortSignal) {
		loading = true;
		error = null;
		data = null;

		const res = await fetch(`/api/v1/admin/pages/${slug}`, { signal });
		if (!res.ok) throw new Error(`Failed to load page (${res.status})`);
		const json = await res.json();
		data = json.data as Page;
		loading = false;
	}

	// Initialise from existing content
	let body = $state('');

	$effect(() => {
		if (data) {
			body = data.content ?? '';
		}
	});

	$effect(() => {
		if (data?.cover_image_id && !coverMedia) {
			(async () => {
				try {
					coverMedia = await getMediaById(data?.cover_image_id!);
				} catch (e) {
					console.error('Failed to fetch banner image:', e);
					coverMedia = null;
					toast.error('Failed to load banner image. Please try again.');
				}
			})();
		}
	});

	// React only to slug changes
	$effect(() => {
		const slug = page.params.slug;
		if (!slug) {
			error = 'Missing slug';
			loading = false;
			return;
		}

		const ctrl = new AbortController();
		fetchPage(slug, ctrl.signal).catch((e) => {
			if (e.name !== 'AbortError') {
				console.error(e);
				error = e.message ?? 'Failed to load page';
				loading = false;
			}
		});

		// cancel in-flight request if slug changes / component reinitialises
		return () => ctrl.abort();
	});

	async function handleUpdate(next: Page): Promise<void> {
		try {
			saving = true;
			next.content = body;
			await updatePage(next, page.params.slug); // IMPORTANT: await this
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
		if (!confirm('Are you sure you want to delete this page? This action cannot be undone.')) {
			return;
		}

		try {
			await deletePage(next.slug); // your DELETE /api/v1/admin/pages/:slug
			toast.success('Page deleted');
			await goto('/admin/pages');
		} catch (err) {
			console.error(err);
			toast.error('Delete failed');
		}
	}

	function updateSlug(): void {
		if (!data) return; // narrow: data is Page here
		const title = data.title ?? '';
		data.slug = slugify(title);
	}

	function handleBannerSelected(item: MediaItem) {
		if (!data) return;
		coverMedia = item;
		data.cover_image_id = item.id;
		bannerModalOpen = false;
	}
</script>

{#if loading}
	<p>Fetching…</p>
{:else if error}
	<p class="text-red-600">Something went wrong: {error}</p>
{:else if data}
	<!-- Consider wiring `saving` into PageForm to disable submit buttons while saving -->
	<div class="flex max-w-2xl flex-col space-y-6">
		<!-- <PageForm content={data} onsubmit={handleUpdate} ondelete={handleDelete} /> -->
		<div class="flex w-full flex-row justify-between gap-8">
			<div class="grid grow grid-cols-1">
				<!-- Title -->
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

				<!-- Slug -->
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
			<div class="">
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
								if (!data) return; // bail if data is null
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

		<!-- end wrapper -->
	</div>
{/if}

<LinkMediaModal
	open={bannerModalOpen}
	onclose={() => (bannerModalOpen = false)}
	onlinked={handleBannerSelected}
	context={{ type: 'page', id: data!.slug }}
	selectOnly={true}
	linkedMediaIds={data?.cover_image_id ? [data?.cover_image_id] : []}
/>
