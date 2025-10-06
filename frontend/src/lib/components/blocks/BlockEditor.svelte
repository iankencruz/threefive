<!-- frontend/src/lib/components/blocks/forms/BlockEditor.svelte -->
<script lang="ts">
import HeaderBlockForm from "./forms/HeaderBlockForm.svelte";
import HeroBlockForm from "./forms/HeroBlockForm.svelte";
import RichTextBlockForm from "./forms/RichTextBlockForm.svelte";

type BlockType = "hero" | "richtext" | "header";

interface Block {
	id?: string;
	type: BlockType;
	data: Record<string, any>;
}

interface Props {
	blocks?: Block[];
	onUpdate?: (blocks: Block[]) => void;
}

let { blocks = $bindable([]), onUpdate }: Props = $props();

const addBlock = (type: BlockType) => {
	const newBlock: Block = {
		type,
		data: getDefaultBlockData(type),
	};
	blocks = [...blocks, newBlock];
	onUpdate?.(blocks);
};

const removeBlock = (index: number) => {
	blocks = blocks.filter((_, i) => i !== index);
	onUpdate?.(blocks);
};

const moveBlock = (index: number, direction: "up" | "down") => {
	const newBlocks = [...blocks];
	const newIndex = direction === "up" ? index - 1 : index + 1;

	if (newIndex >= 0 && newIndex < newBlocks.length) {
		[newBlocks[index], newBlocks[newIndex]] = [
			newBlocks[newIndex],
			newBlocks[index],
		];
		blocks = newBlocks;
		onUpdate?.(blocks);
	}
};

const updateBlockData = (index: number, data: Record<string, any>) => {
	blocks[index].data = data;
	onUpdate?.(blocks);
};

const getDefaultBlockData = (type: BlockType): Record<string, any> => {
	switch (type) {
		case "hero":
			return { title: "", subtitle: "", cta_text: "", cta_url: "" };
		case "richtext":
			return { content: "" };
		case "header":
			return { heading: "", subheading: "", level: "h2" };
		default:
			return {};
	}
};

const getBlockTypeColor = (type: BlockType) => {
	switch (type) {
		case "hero":
			return "bg-blue-100 text-blue-800";
		case "richtext":
			return "bg-purple-100 text-purple-800";
		case "header":
			return "bg-green-100 text-green-800";
		default:
			return "bg-gray-100 text-gray-800";
	}
};
</script>

<div class="block-editor">
	<div class="flex justify-between items-center mb-6">
		<h2 class="text-xl font-semibold text-gray-900">Content Blocks</h2>
		<div class="flex gap-2">
			<button
				onclick={() => addBlock('hero')}
				type="button"
				class="px-3 py-1.5 bg-blue-50 text-blue-700 rounded-lg hover:bg-blue-100 text-sm font-medium transition-colors"
			>
				+ Hero
			</button>
			<button
				onclick={() => addBlock('richtext')}
				type="button"
				class="px-3 py-1.5 bg-purple-50 text-purple-700 rounded-lg hover:bg-purple-100 text-sm font-medium transition-colors"
			>
				+ Rich Text
			</button>
			<button
				onclick={() => addBlock('header')}
				type="button"
				class="px-3 py-1.5 bg-green-50 text-green-700 rounded-lg hover:bg-green-100 text-sm font-medium transition-colors"
			>
				+ Header
			</button>
		</div>
	</div>

	{#if blocks.length === 0}
		<div class="text-center py-12 bg-gray-50 rounded-lg border-2 border-dashed border-gray-300">
			<svg class="w-12 h-12 mx-auto text-gray-400 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
			</svg>
			<p class="text-gray-600">No blocks added yet. Click a button above to add content blocks.</p>
		</div>
	{:else}
		<div class="space-y-4">
			{#each blocks as block, index (index)}
				<div class="border border-gray-200 rounded-lg p-4 bg-white shadow-sm hover:shadow-md transition-shadow">
					<div class="flex justify-between items-start mb-4">
						<span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium capitalize {getBlockTypeColor(block.type)}">
							{block.type}
						</span>
						<div class="flex gap-2">
							<button
								onclick={() => moveBlock(index, 'up')}
								disabled={index === 0}
								type="button"
								class="p-1 text-gray-600 hover:text-gray-900 disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
								aria-label="Move block up"
							>
								<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
								</svg>
							</button>
							<button
								onclick={() => moveBlock(index, 'down')}
								disabled={index === blocks.length - 1}
								type="button"
								class="p-1 text-gray-600 hover:text-gray-900 disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
								aria-label="Move block down"
							>
								<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
								</svg>
							</button>
							<button
								onclick={() => removeBlock(index)}
								type="button"
								class="p-1 text-red-600 hover:text-red-900 transition-colors"
								aria-label="Remove block"
							>
								<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
								</svg>
							</button>
						</div>
					</div>

					{#if block.type === 'hero'}
						<HeroBlockForm bind:data={block.data} onchange={(data) => updateBlockData(index, data)} />
					{:else if block.type === 'richtext'}
						<RichTextBlockForm bind:data={block.data} onchange={(data) => updateBlockData(index, data)} />
					{:else if block.type === 'header'}
						<HeaderBlockForm bind:data={block.data} onchange={(data) => updateBlockData(index, data)} />
					{/if}
				</div>
			{/each}
		</div>
	{/if}
</div>
