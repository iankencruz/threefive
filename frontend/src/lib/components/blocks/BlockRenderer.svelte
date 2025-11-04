<!-- frontend/src/lib/components/blocks/BlockRenderer.svelte -->
<script lang="ts">
	import HeroBlock from "./display/HeroBlock.svelte";
	import RichTextBlock from "./display/RichTextBlock.svelte";
	import HeaderBlock from "./display/HeaderBlock.svelte";
	import GalleryBlock from "./display/GalleryBlock.svelte";

	export interface Block {
		id?: string;
		type: string;
		data: Record<string, any>;
		sort_order?: number;
	}

	interface Props {
		blocks: Block[];
		mediaMap?: Record<string, any>; // ✨ NEW: Pre-loaded media
	}

	let { blocks = [], mediaMap = {} }: Props = $props();

	// Sort blocks by sort_order if available
	const sortedBlocks = $derived(
		[...blocks].sort((a, b) => (a.sort_order || 0) - (b.sort_order || 0)),
	);
</script>

{#if sortedBlocks.length === 0}
	<div class="min-h-screen flex items-center justify-center bg-gray-50">
		<div class="text-center">
			<svg class="w-16 h-16 mx-auto text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
			</svg>
			<p class="text-gray-600 text-lg">No content blocks available</p>
		</div>
	</div>
{:else}
	{#each sortedBlocks as block (block.id || block.type)}
		{#if block.type === 'hero'}
			<!-- ✨ Pass pre-loaded media to HeroBlock -->
			<HeroBlock data={block.data} media={mediaMap[block.data?.image_id]} />
		{:else if block.type === 'richtext'}
			<RichTextBlock data={block.data} />
		{:else if block.type === 'header'}
			<HeaderBlock data={block.data} />
		{:else if block.type === 'gallery'}
			<GalleryBlock data={block.data} mediaMap={mediaMap} />
		{:else}
			<!-- Unknown block type -->
			<div class="py-8 bg-yellow-50 border-l-4 border-yellow-400">
				<div class="container mx-auto px-4 max-w-4xl">
					<p class="text-yellow-700">
						Unknown block type: <strong>{block.type}</strong>
					</p>
				</div>
			</div>
		{/if}
	{/each}
{/if}
