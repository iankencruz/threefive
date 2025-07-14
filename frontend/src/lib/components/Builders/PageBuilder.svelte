<script lang="ts">
	import AddBlockMenu from '$lib/components/Builders/AddBlockMenu.svelte';
	import BlockRenderer from '$lib/components/Builders/BlockRenderer.svelte';
	import { getDefaultProps } from '$lib/components/Builders/blockDefaults';
	import type { Block } from '$src/lib/types';
	import { v4 as uuidv4 } from 'uuid';

	let { content = $bindable() }: { content: Block[] } = $props();

	function insertBlock(index: number, type: string): void {
		const newBlock: Block = {
			id: uuidv4(),
			type,
			props: getDefaultProps(type)
		};
		content.splice(index, 0, newBlock);
	}

	function removeBlock(index: number): void {
		content.splice(index, 1);
	}

	function updateBlock(index: number, updated: Block) {
		content[index] = updated;
		content = content; // 🧠 ← This is what syncs it back to PageForm
	}
</script>

<div class="space-y-3 rounded-md border bg-gray-50 p-4">
	{#each content as block, index (block.id)}
		<div class="group relative rounded border bg-white p-4 shadow-sm">
			<BlockRenderer {block} onupdate={(updatedBlock) => updateBlock(index, updatedBlock)} />

			<!-- Controls -->
			<div class="absolute -top-2 -right-2 opacity-0 transition group-hover:opacity-100">
				<button
					onclick={() => removeBlock(index)}
					class="rounded-full bg-red-600 px-2 py-1 text-xs text-white"
				>
					✕
				</button>
			</div>
		</div>

		<!-- Insert Block Button -->
		<div class="flex justify-center">
			<AddBlockMenu onselect={(type: string) => insertBlock(index + 1, type)} />
		</div>
	{/each}

	<!-- Initial Block if empty -->
	{#if content.length === 0}
		<div class="text-center text-sm text-gray-500">No blocks yet. Add one below.</div>
		<div class="mt-4 flex justify-center">
			<AddBlockMenu onselect={(type: string) => insertBlock(0, type)} />
		</div>
	{/if}
</div>
