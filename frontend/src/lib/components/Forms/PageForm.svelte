<script lang="ts">
	import type { Block, Page } from '$lib/types';
	import PageBuilder from '$lib/components/Builders/PageBuilder.svelte';
	import { toast } from 'svelte-sonner';
	import type { MediaItem } from '$lib/types';
	import { getMediaById } from '$lib/api/media';
	import LinkMediaModal from '../Media/LinkMediaModal.svelte';
	import BlockRenderer from '../Builders/BlockRenderer.svelte';
	import { sortBlocks } from '$lib/api/pages';

	let {
		content,
		blocks,
		onsubmit,
		ondelete
	}: {
		content: Page;
		blocks: Block[];
		onsubmit: (data: { page: Page; blocks: Block[] }) => void;
		ondelete?: (data: Page) => void;
	} = $props();

	// Ensure `blocks` is never undefined for bind
	if (!blocks) {
		blocks = [];
	}

	let localBlocks = $state<Block[]>(blocks);

	function updateSlug(): void {
		content.slug = content.title
			.toLowerCase()
			.replace(/[^\w\s-]/g, '')
			.replace(/\s+/g, '-')
			.replace(/-+/g, '-')
			.trim();
	}

	let previewMode = $state(false);

	function togglePreview() {
		previewMode = !previewMode;
	}

	let coverMedia = $state<MediaItem | null>(null);
	let bannerModalOpen = $state(false);

	$effect(() => {
		if (content.cover_image_id && !coverMedia) {
			(async () => {
				try {
					coverMedia = await getMediaById(content.cover_image_id!);
				} catch (e) {
					console.error('Failed to fetch banner image:', e);
					coverMedia = null;
					toast.error('Failed to load banner image. Please try again.');
				}
			})();
		}
	});

	function handleBannerSelected(item: MediaItem) {
		coverMedia = item;
		content.cover_image_id = item.id;
		bannerModalOpen = false;
	}

	$effect(() => {
		if (!content.is_draft && !content.is_published) {
			content.is_draft = true;
		}
	});

	$effect(() => {
		if (content.is_published) content.is_draft = false;
		if (content.is_draft) content.is_published = false;
	});

	async function handleSavePage() {
		try {
			await sortBlocks(content.slug, localBlocks); // 👈 persist block order
			onsubmit({ page: content, blocks: localBlocks }); // 👈 then continue with save
			toast.success('Page saved');
		} catch (e) {
			toast.error('Failed to save page');
			console.error(e);
		}
	}

	function handleReorder(updatedBlocks: Block[]) {
		// Optionally sort & reindex
		const reordered = updatedBlocks.map((b, i) => ({
			...b,
			sort_order: i
		}));

		onsubmit({ page: content, blocks: reordered });
	}
</script>

<div class="space-y-6 px-4">
	<!-- Title -->
	<div>
		<label for="title" class="block text-sm font-medium text-gray-700">Title</label>
		<input
			name="title"
			type="text"
			class="mt-1 w-full rounded border px-3 py-2"
			bind:value={content.title}
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
			bind:value={content.slug}
			readonly
		/>
	</div>

	<!-- Banner Image -->
	<div>
		<p class="block text-sm font-medium text-gray-700">Banner Image</p>

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
						content.cover_image_id = null;
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

	<!-- SEO Fields -->
	<div>
		<label for="seo-title" class="block text-sm font-medium text-gray-700">SEO Title</label>
		<input
			name="seo-title"
			type="text"
			class="mt-1 w-full rounded border px-3 py-2"
			bind:value={content.seo_title}
		/>
	</div>

	<div>
		<label for="seo-description" class="block text-sm font-medium text-gray-700"
			>SEO Description</label
		>
		<textarea
			name="seo-description"
			class="mt-1 w-full rounded border px-3 py-2"
			rows="3"
			bind:value={content.seo_description}
		></textarea>
	</div>

	<div>
		<label for="seo-canonical" class="block text-sm font-medium text-gray-700"
			>SEO Canonical URL</label
		>
		<input
			name="seo-canonical"
			type="text"
			class="mt-1 w-full rounded border px-3 py-2"
			bind:value={content.seo_canonical}
		/>
	</div>

	<!-- Page Builder -->
	<div>
		<!-- Render builder or preview -->
		{#if previewMode}
			{#each localBlocks as block (block.id)}
				<BlockRenderer {block} onupdate={() => {}} />
			{/each}
		{:else}
			<h3 class="text-sm font-semibold text-gray-800">Content Blocks</h3>
			<PageBuilder bind:blocks={localBlocks} onreorder={handleReorder} />
		{/if}
	</div>

	<!-- Save Controls -->
	<div class="flex items-center justify-between">
		<div class="flex items-center gap-3">
			<label><input type="checkbox" bind:checked={content.is_draft} /> Draft</label>
			<label><input type="checkbox" bind:checked={content.is_published} /> Published</label>
		</div>
		<!-- Add toggle button -->
		<div class="mb-4 flex justify-end">
			<button
				class="rounded bg-gray-200 px-4 py-1 text-sm hover:bg-gray-300"
				onclick={togglePreview}
			>
				{previewMode ? 'Switch to Edit Mode' : 'Preview Page'}
			</button>
		</div>
		<div class="flex gap-x-4">
			{#if ondelete}
				<button
					onclick={() => ondelete?.(content)}
					class="rounded bg-red-600 px-4 py-2 text-white hover:bg-gray-800"
				>
					Delete
				</button>
			{/if}
			<button
				onclick={() => onsubmit({ page: content, blocks: localBlocks })}
				class="rounded bg-black px-4 py-2 text-white hover:bg-gray-800"
			>
				Save Page
			</button>
		</div>
	</div>
</div>

<LinkMediaModal
	open={bannerModalOpen}
	onclose={() => (bannerModalOpen = false)}
	onlinked={handleBannerSelected}
	context={{ type: 'page', id: content.slug }}
	selectOnly={true}
	linkedMediaIds={content.cover_image_id ? [content.cover_image_id] : []}
/>
