<!-- frontend/src/lib/components/blocks/BlockEditor.svelte -->
<script lang="ts">
import HeaderBlockForm, { type HeaderBlockData } from "./forms/HeaderBlockForm.svelte";
import HeroBlockForm, { type HeroBlockData } from "./forms/HeroBlockForm.svelte";
import RichTextBlockForm, { type RichtextBlockData } from "./forms/RichTextBlockForm.svelte";

interface BlockTypeMap {
	hero: HeroBlockData;
	richtext: RichtextBlockData;
	header: HeaderBlockData;
	// gallery: GalleryBlockData; // ✅ Just add new ones here
}

type BlockType = keyof BlockTypeMap;
type BlockData = BlockTypeMap[BlockType];

interface Block {
	id?: string;
	type: BlockType;
	data: BlockData;
}

interface Props {
	blocks?: Block[];
	onUpdate?: (blocks: Block[]) => void;
}

let { blocks = $bindable([]), onUpdate }: Props = $props();

// Track which separator is showing the menu
let activeSeparator = $state<number | null>(null);

const addBlockAt = (type: BlockType, position: number) => {
	const newBlock: Block = {
		type,
		data: getDefaultBlockData(type),
	};

	// Insert block at specific position
	const newBlocks = [...blocks];
	newBlocks.splice(position, 0, newBlock);
	blocks = newBlocks;

	// Close the menu
	activeSeparator = null;
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
		[newBlocks[index], newBlocks[newIndex]] = [newBlocks[newIndex], newBlocks[index]];
		blocks = newBlocks;
		onUpdate?.(blocks);
	}
};

const updateBlockData = (index: number, data: Record<string, any>) => {
	blocks[index].data = data;
	onUpdate?.(blocks);
};

const toggleSeparator = (index: number) => {
	activeSeparator = activeSeparator === index ? null : index;
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

// Close separator menu when clicking outside
$effect(() => {
	const handleClick = (e: MouseEvent) => {
		const target = e.target as HTMLElement;
		if (!target.closest(".separator-menu") && !target.closest(".separator-trigger")) {
			activeSeparator = null;
		}
	};

	if (activeSeparator !== null) {
		document.addEventListener("click", handleClick);
		return () => document.removeEventListener("click", handleClick);
	}
});
</script>

<div class="block-editor">
	<div class="flex justify-between items-center mb-6">
		<h2 class="text-xl font-semibold text-gray-900">Content Blocks</h2>
	</div>

	{#if blocks.length === 0}
		<div class="text-center py-12 bg-gray-50 rounded-lg border-2 border-dashed border-gray-300">
			<svg class="w-12 h-12 mx-auto text-gray-400 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
			</svg>
			<p class="text-gray-600 mb-4">No blocks added yet. Click below to add your first block.</p>
			
			<!-- Initial Add Block Menu -->
			<div class="flex justify-center gap-2">
				<button
					onclick={() => addBlockAt('hero', 0)}
					type="button"
					class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 text-sm font-medium transition-colors"
				>
					+ Hero Block
				</button>
				<button
					onclick={() => addBlockAt('richtext', 0)}
					type="button"
					class="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 text-sm font-medium transition-colors"
				>
					+ Rich Text
				</button>
				<button
					onclick={() => addBlockAt('header', 0)}
					type="button"
					class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 text-sm font-medium transition-colors"
				>
					+ Header
				</button>
			</div>
		</div>
	{:else}
		<div class="space-y-1">
			<!-- Insert separator at the top -->
<!-- Insert separator at the top -->
<div class="relative my-2">
	<button
		type="button"
		class="separator-trigger w-full py-2 flex items-center group rounded transition-colors"
		onclick={(e) => { e.stopPropagation(); toggleSeparator(0); }}
	>
		<div aria-hidden="true" class="w-full border-t border-gray-300 group-hover:border-gray-400 transition-colors"></div>
		<div class="relative flex justify-center px-3">
      <span class=" p-px rounded-full group-hover:bg-gray-900 flex items-center gap-1.5 text-gray-500 group-hover:text-white transition-colors">
				<svg viewBox="0 0 20 20" fill="currentColor" aria-hidden="true" class="w-5 h-5">
					<path d="M10.75 4.75a.75.75 0 0 0-1.5 0v4.5h-4.5a.75.75 0 0 0 0 1.5h4.5v4.5a.75.75 0 0 0 1.5 0v-4.5h4.5a.75.75 0 0 0 0-1.5h-4.5v-4.5Z" />
				</svg>
			</span>
		</div>
		<div aria-hidden="true" class="w-full border-t border-gray-300 group-hover:border-gray-400 transition-colors"></div>
	</button>
	{#if activeSeparator === 0}
		<div class="separator-menu absolute top-full left-1/2 -translate-x-1/2 z-10 mt-1 bg-white rounded-lg shadow-lg border border-gray-200 p-2 flex gap-2">
			<button
				onclick={() => addBlockAt('hero', 0)}
				type="button"
				class="px-4 py-2 bg-blue-50 text-blue-700 rounded-lg hover:bg-blue-100 text-sm font-medium transition-colors whitespace-nowrap"
			>
				Hero Block
			</button>
			<button
				onclick={() => addBlockAt('richtext', 0)}
				type="button"
				class="px-4 py-2 bg-purple-50 text-purple-700 rounded-lg hover:bg-purple-100 text-sm font-medium transition-colors whitespace-nowrap"
			>
				Rich Text
			</button>
			<button
				onclick={() => addBlockAt('header', 0)}
				type="button"
				class="px-4 py-2 bg-green-50 text-green-700 rounded-lg hover:bg-green-100 text-sm font-medium transition-colors whitespace-nowrap"
			>
				Header
			</button>
		</div>
	{/if}
</div>
			{#each blocks as block, index (index)}
				<!-- Block Content -->
				<div class="border border-gray-200 rounded-lg p-8 bg-white shadow-sm hover:shadow-md transition-shadow">
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
						<HeroBlockForm data={block.data} onchange={(data) => updateBlockData(index, data)} />
					{:else if block.type === 'richtext'}
						<RichTextBlockForm data={block.data} onchange={(data) => updateBlockData(index, data)} />
					{:else if block.type === 'header'}
						<HeaderBlockForm data={block.data} onchange={(data) => updateBlockData(index, data)} />
					{/if}
				</div>

				<!-- Insert separator after each block -->
        <div class="relative my-2">
          <button
            type="button"
            class="separator-trigger w-full py-2 flex items-center group  rounded transition-colors"
            onclick={(e) => { e.stopPropagation(); toggleSeparator(index + 1); }}
          >
            <div aria-hidden="true" class="w-full border-t border-gray-300 group-hover:border-gray-400 transition-colors"></div>
            <div class="relative flex justify-center px-3">
              <span class=" p-px rounded-full group-hover:bg-gray-900 flex items-center gap-1.5 text-gray-500 group-hover:text-white transition-colors">
                <svg viewBox="0 0 20 20" fill="currentColor" aria-hidden="true" class="w-5 h-5">
                  <path d="M10.75 4.75a.75.75 0 0 0-1.5 0v4.5h-4.5a.75.75 0 0 0 0 1.5h4.5v4.5a.75.75 0 0 0 1.5 0v-4.5h4.5a.75.75 0 0 0 0-1.5h-4.5v-4.5Z" />
                </svg>
              </span>
            </div>
            <div aria-hidden="true" class="w-full border-t border-gray-300 group-hover:border-gray-400 transition-colors"></div>
          </button>
          {#if activeSeparator === index + 1}
            <div class="separator-menu absolute top-full left-1/2 -translate-x-1/2 z-10 mt-1 bg-white rounded-lg shadow-lg border border-gray-200 p-2 flex gap-2">
              <button
                onclick={() => addBlockAt('hero', index + 1)}
                type="button"
                class="px-4 py-2 bg-blue-50 text-blue-700 rounded-lg hover:bg-blue-100 text-sm font-medium transition-colors whitespace-nowrap"
              >
                Hero Block
              </button>
              <button
                onclick={() => addBlockAt('richtext', index + 1)}
                type="button"
                class="px-4 py-2 bg-purple-50 text-purple-700 rounded-lg hover:bg-purple-100 text-sm font-medium transition-colors whitespace-nowrap"
              >
                Rich Text
              </button>
              <button
                onclick={() => addBlockAt('header', index + 1)}
                type="button"
                class="px-4 py-2 bg-green-50 text-green-700 rounded-lg hover:bg-green-100 text-sm font-medium transition-colors whitespace-nowrap"
              >
                Header
              </button>
            </div>
          {/if}
        </div>
			{/each}
		</div>
	{/if}
</div>
