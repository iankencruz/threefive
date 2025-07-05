<script lang="ts">
	import HeadingBlock from '../Blocks/HeadingBlock.svelte';

	let {
		block,
		onupdate
	}: {
		block: { type: string; props: any };
		onupdate: (updatedBlock: any) => void;
	} = $props();

	function handleUpdate(updatedProps: any) {
		onupdate({ ...block, props: updatedProps });
	}
</script>

{#if block.type === 'heading'}
	<HeadingBlock props={block.props} onupdate={handleUpdate} />
{:else}
	<div class="text-red-500">⚠ Unknown block type: {block.type}</div>
{/if}
