<script lang="ts">
	import type { Page, MediaItem } from '$lib/types';
	import { toast } from 'svelte-sonner';
	import { getMediaById } from '$lib/api/media';
	import LinkMediaModal from '../Media/LinkMediaModal.svelte';
	import TipTap from '../Builders/TipTap.svelte';
	import { Editor } from '@tiptap/core';

	let {
		content = $bindable(),
		onsubmit,
		ondelete
	}: {
		content: Page;
		onsubmit: (data: Page) => void;
		ondelete?: (data: Page) => void;
	} = $props();

	let localContent = $state({ ...content });

	// Initialise from existing content
	let body = $state(localContent.content ?? '');

	// Always mirror editor HTML into the payload
	$effect(() => {
		localContent.content = body;
	});

	function updateSlug(): void {
		localContent.slug = localContent.title
			.toLowerCase()
			.replace(/[^\w\s-]/g, '')
			.replace(/\s+/g, '-')
			.replace(/-+/g, '-')
			.trim();
	}

	let coverMedia = $state<MediaItem | null>(null);
	let bannerModalOpen = $state(false);

	$effect(() => {
		if (localContent.cover_image_id && !coverMedia) {
			(async () => {
				try {
					coverMedia = await getMediaById(localContent.cover_image_id!);
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
		localContent.cover_image_id = item.id;
		bannerModalOpen = false;
	}
</script>

<div class="w-full space-y-6 px-4">
	<div class="flex w-full flex-row justify-between gap-8">
		<div class="grid grow grid-cols-1">
			<!-- Title -->
			<div>
				<label for="title" class="block text-sm font-medium text-gray-700">Title</label>
				<input
					name="title"
					type="text"
					class="mt-1 w-full rounded border px-3 py-2"
					bind:value={localContent.title}
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
					bind:value={localContent.slug}
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
							localContent.cover_image_id = null;
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

	<!-- SEO Fields -->
	<div>
		<label for="seo-title" class="block text-sm font-medium text-gray-700">SEO Title</label>
		<input
			name="seo-title"
			type="text"
			class="mt-1 w-full rounded border px-3 py-2"
			bind:value={localContent.seo_title}
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
			bind:value={localContent.seo_description}
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
			bind:value={localContent.seo_canonical}
		/>
	</div>

	<!-- Save Controls -->
	<div class="flex items-center justify-between">
		<div class="flex gap-x-4">
			{#if ondelete}
				<button
					onclick={() => ondelete?.(localContent)}
					class="rounded bg-red-600 px-4 py-2 text-white hover:bg-gray-800"
				>
					Delete
				</button>
			{/if}
			<button
				onclick={() => onsubmit(localContent)}
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
	context={{ type: 'page', id: localContent.slug }}
	selectOnly={true}
	linkedMediaIds={localContent.cover_image_id ? [localContent.cover_image_id] : []}
/>
