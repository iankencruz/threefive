<script lang="ts">
	import type { Page } from '$lib/types';
	import PageBuilder from '$lib/components/Builders/PageBuilder.svelte';
	import { toast } from 'svelte-sonner';
	import { page } from '$app/state';

	let {
		content,
		onsubmit,
		ondelete
	}: {
		content: Page;
		onsubmit: (data: Page) => void;
		ondelete?: (data: Page) => void;
	} = $props();

	function updateSlug(): void {
		content.slug = content.title
			.toLowerCase()
			.replace(/[^\w\s-]/g, '')
			.replace(/\s+/g, '-')
			.replace(/-+/g, '-')
			.trim();
	}

	// function handleSubmit() {
	// 	onsubmit(content); // emit back to parent
	// }
	//
	// function handleDelete() {
	// 	ondelete(content); // emit back to parent
	// }

	$effect(() => {
		if (!content.is_draft && !content.is_published) {
			content.is_draft = true;
		}
	});

	// Effect to handle mutual exclusivity
	$effect(() => {
		if (content.is_published) content.is_draft = false;
	});

	$effect(() => {
		if (content.is_draft) content.is_published = false;
	});
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
		<label for="banner-image" class="block text-sm font-medium text-gray-700">Banner Image</label>
		<input
			name="banner-image"
			type="text"
			class="mt-1 w-full rounded border px-3 py-2"
			bind:value={content.banner_image_id}
			placeholder="Enter media ID manually"
		/>
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
		<h3 class="text-sm font-semibold text-gray-800">Content Blocks</h3>
		<PageBuilder bind:content={content.content} />
	</div>

	<!-- Save Controls -->
	<div class="flex items-center justify-between">
		<div class="flex items-center gap-3">
			<label><input type="checkbox" bind:checked={content.is_draft} /> Draft</label>
			<label><input type="checkbox" bind:checked={content.is_published} /> Published</label>
		</div>
		<div class="flex gap-x-4">
			<button
				onclick={() => ondelete(content)}
				class="rounded bg-red-600 px-4 py-2 text-white hover:bg-gray-800"
			>
				Delete
			</button>
			<button
				onclick={() => onsubmit(content)}
				class="rounded bg-black px-4 py-2 text-white hover:bg-gray-800"
			>
				Save Page
			</button>
		</div>
	</div>
</div>
