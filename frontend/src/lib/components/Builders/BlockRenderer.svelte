<script lang="ts">
	import HeadingBlock from '../Blocks/HeadingBlock.svelte';
	import type { Block } from '$lib/types';
	import ImageBlock from '../Blocks/ImageBlock.svelte';
	import RichTextBlock from '../Blocks/RichTextBlock.svelte';

	let {
		block,
		onupdate
	}: {
		block: Block;
		onupdate: (updated: Block) => void;
	} = $props();

	function handleUpdate(updatedBlock: Block) {
		onupdate(updatedBlock); // pass it directly
	}
</script>

{#if block?.type === 'heading'}
	<HeadingBlock {block} onupdate={handleUpdate} />
{:else if block?.type === 'image'}
	<ImageBlock {block} onupdate={handleUpdate} />
{:else if block.type === 'richtext'}
	<RichTextBlock {block} onupdate={handleUpdate} />
{:else}
	<div class="text-red-500">⚠ Unknown block type: {block?.type as Block}</div>
{/if}
