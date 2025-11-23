<!-- frontend/src/routes/admin/projects/[id]/preview/+page.svelte -->
<script lang="ts">
	import { goto } from '$app/navigation';
	import BlockRenderer from '$components/blocks/BlockRenderer.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	const getStatusBadge = (status: string) => {
		switch (status) {
			case 'published':
				return 'bg-green-100 text-green-800';
			case 'draft':
				return 'bg-yellow-100 text-yellow-800';
			case 'archived':
				return 'bg-gray-100 text-gray-800';
			default:
				return 'bg-gray-100 text-gray-800';
		}
	};
</script>

<svelte:head>
	<title>Preview: {data.project?.title}</title>
</svelte:head>

<!-- Preview Banner -->
<div class="sticky top-0 z-50 border-b border-yellow-400 bg-yellow-50 px-4 py-3 shadow-sm">
	<div class="mx-auto flex max-w-7xl items-center justify-between">
		<div class="flex items-center gap-4">
			<span class="flex items-center gap-2 text-sm font-medium text-yellow-800">
				<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
					/>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
					/>
				</svg>
				Preview Mode
			</span>
			<span
				class="inline-flex items-center rounded-full px-3 py-1 text-xs font-medium capitalize {getStatusBadge(
					data.project.status
				)}"
			>
				{data.project.status}
			</span>
		</div>
		<div class="flex items-center gap-2">
			<button
				onclick={() => goto(`/admin/projects/${data.project.id}/edit`)}
				class="rounded-lg bg-white px-4 py-2 text-sm font-medium text-gray-700 shadow-sm transition-colors hover:bg-gray-50"
			>
				Edit
			</button>
			<button
				onclick={() => goto('/admin/projects')}
				class="rounded-lg bg-gray-800 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-gray-700"
			>
				Close Preview
			</button>
		</div>
	</div>
</div>

<!-- Page Content -->
<div class="mx-auto max-w-7xl px-4 py-8">
	<article>
		<!-- Page Header -->
		<header class="mb-12">
			<h1 class="mb-4 text-4xl font-bold">{data.project.title}</h1>
			{#if data.project.seo?.meta_description}
				<p class="text-xl text-gray-600">{data.project.seo.meta_description}</p>
			{/if}
		</header>

		<!-- Blocks -->
		{#if data.project.blocks && data.project.blocks.length > 0}
			<BlockRenderer blocks={data.project.blocks} />
		{:else}
			<div class="rounded-lg border-2 border-dashed border-gray-300 bg-gray-50 px-8 py-16 text-center">
				<p class="text-gray-500">No content blocks yet</p>
			</div>
		{/if}
	</article>
</div>
