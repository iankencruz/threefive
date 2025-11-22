<!-- frontend/src/lib/components/blocks/BlockRenderer.svelte -->
<script lang="ts">
	import HeroBlock from "./display/HeroBlock.svelte";
	import RichTextBlock from "./display/RichTextBlock.svelte";
	import HeaderBlock from "./display/HeaderBlock.svelte";
	import GalleryBlock from "./display/GalleryBlock.svelte";	

	interface Block {
		id: string;
		type: string;
		data: Record<string, any>;
		media?: any[]; // Array of linked media
	}

	interface Props {
		blocks: Block[];
	}

	let { blocks }: Props = $props();

	// Helper to get first media for a block (for hero backgrounds)
	const getBlockMedia = (block: Block) => {
		return block.media && block.media.length > 0 ? block.media[0] : null;
	};
</script>

<div class="blocks-container">
	{#each blocks as block (block.id)}
		{#if block.type === 'hero'}
			<HeroBlock 
				data={block.data as any} 
				media={getBlockMedia(block)} 
			/>
		{:else if block.type === 'richtext'}
			<RichTextBlock data={block.data as any} />
		{:else if block.type === 'header'}
			<HeaderBlock data={block.data as any} />
		{:else if block.type === 'gallery'}
			<GalleryBlock 
				data={block.data as any}
				media={block.media || []}
			/>
		{:else}
			<div class="unknown-block">
				<p>Unknown block type: {block.type}</p>
			</div>
		{/if}
	{/each}
</div>
