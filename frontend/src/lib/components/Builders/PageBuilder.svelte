<script lang="ts">
	import AddBlockMenu from '$lib/components/Builders/AddBlockMenu.svelte';
	import BlockRenderer from '$lib/components/Builders/BlockRenderer.svelte';
	import type { Block } from '$lib/types';
	import { ArrowDown, ArrowUp } from '@lucide/svelte';
	import { v4 as uuidv4 } from 'uuid';

	let {
		blocks = $bindable([]),
		onreorder
	}: {
		blocks: Block[];
		onreorder?: (blocks: Block[]) => void;
	} = $props();

	function insertBlock(index: number, type: Block['type']): void {
		if (type === 'heading') {
			const newBlock: Block = {
				id: uuidv4(),
				type: 'heading',
				sort_order: index,
				props: {
					text: '',
					level: 1
				}
			};
			blocks.splice(index, 0, newBlock);
			return;
		}

		if (type === 'richtext') {
			const newBlock: Block = {
				id: uuidv4(),
				type: 'richtext',
				sort_order: index,
				props: {
					html: ''
				}
			};
			blocks.splice(index, 0, newBlock);
			return;
		}

		if (type === 'image') {
			const newBlock: Block = {
				id: uuidv4(),
				type: 'image',
				sort_order: index,
				props: {
					media_id: '',
					alt_text: '',
					align: 'center',
					object_fit: 'cover'
				}
			};
			blocks.splice(index, 0, newBlock);
			return;
		}

		console.warn(`Unsupported block type: ${type}`);
	}

	function removeBlock(index: number): void {
		blocks.splice(index, 1);
	}

	function updateBlock(index: number, updated: Block) {
		blocks[index] = updated;
		blocks = blocks; // 🧠 ← This is what syncs it back to PageForm
	}

	function moveBlockUp(index: number) {
		if (index <= 0) return;

		// Swap sort_order
		[blocks[index - 1].sort_order, blocks[index].sort_order] = [
			blocks[index].sort_order,
			blocks[index - 1].sort_order
		];

		// Reorder in array
		[blocks[index - 1], blocks[index]] = [blocks[index], blocks[index - 1]];
		blocks = blocks; // Trigger bind

		if (onreorder) onreorder(blocks);
	}

	function moveBlockDown(index: number) {
		if (index === blocks.length - 1) return;

		[blocks[index + 1].sort_order, blocks[index].sort_order] = [
			blocks[index].sort_order,
			blocks[index + 1].sort_order
		];

		[blocks[index + 1], blocks[index]] = [blocks[index], blocks[index + 1]];
		blocks = blocks; // Trigger bind

		if (onreorder) onreorder(blocks);
	}
</script>

<div class="space-y-3 rounded-md border bg-gray-50 p-4">
	{#each blocks as block, index (block.id)}
		<div class="group relative rounded border bg-white p-4 shadow-sm">
			<BlockRenderer {block} onupdate={(updatedBlock) => updateBlock(index, updatedBlock)} />

			<!-- Sort Buttons -->
			<div
				class="absolute top-0 -left-13 flex flex-col gap-2 opacity-0 transition group-hover:opacity-100"
			>
				<button
					type="button"
					class=" rounded-md bg-white p-1 text-xs outline
          {index === 0 ? '  text-gray-300' : 'cursor-pointer text-gray-500 hover:text-black'}"
					onclick={() => moveBlockUp(index)}
					disabled={index === 0}
					title="Move up"
				>
					<ArrowUp />
				</button>

				<button
					type="button"
					class=" rounded-md bg-white p-1 text-xs outline
          {index === blocks.length - 1
						? 'cursor-default text-gray-300'
						: 'cursor-pointer text-gray-500 hover:text-black'}"
					onclick={() => moveBlockDown(index)}
					disabled={index === blocks.length - 1}
					title="Move down"
				>
					<ArrowDown />
				</button>
			</div>

			<span class="text-xs text-gray-400">#{index + 1}</span>

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
			<AddBlockMenu onselect={(type: string) => insertBlock(index + 1, type as Block['type'])} />
		</div>
	{/each}

	<!-- Initial Block if empty -->
	{#if blocks.length === 0}
		<div class="text-center text-sm text-gray-500">No blocks yet. Add one below.</div>
		<div class="mt-4 flex justify-center">
			<AddBlockMenu onselect={(type: string) => insertBlock(0, type as Block['type'])} />
		</div>
	{/if}
</div>
