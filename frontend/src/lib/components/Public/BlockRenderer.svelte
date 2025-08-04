<script lang="ts">
	import HeadingBlock from './HeadingBlock.svelte';
	import ImageBlock from './ImageBlock.svelte';
	import RichTextBlock from './RichTextBlock.svelte';
	import type { Block } from '$lib/types';

	let { blocks = [] }: { blocks: Block[] } = $props();
</script>

{#each blocks as block (block.id)}
	{console.log('type: ', block.type)}
	{#if block.type === 'heading'}
		<HeadingBlock {block} />
	{:else if block.type === 'image'}
		<ImageBlock {block} />
	{:else if block.type === 'richtext'}
		<RichTextBlock {block} />
	{:else}
		<p class="text-red-600">Unsupported block type: {(block as any).type}</p>
	{/if}
{/each}
