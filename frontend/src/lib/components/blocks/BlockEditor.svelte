<!-- frontend/src/lib/components/blocks/BlockEditor.svelte -->
<script lang="ts">
	import GalleryBlockForm, { type GalleryBlockData } from './forms/GalleryBlockForm.svelte';
	import HeaderBlockForm, { type HeaderBlockData } from './forms/HeaderBlockForm.svelte';
	import HeroBlockForm, { type HeroBlockData } from './forms/HeroBlockForm.svelte';
	import RichTextBlockForm, { type RichtextBlockData } from './forms/RichTextBlockForm.svelte';
	import AboutBlockForm, { type AboutBlockData } from './forms/AboutBlockForm.svelte';

	interface BlockTypeMap {
		hero: HeroBlockData;
		richtext: RichtextBlockData;
		header: HeaderBlockData;
		gallery: GalleryBlockData;
		about: AboutBlockData;
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
			data: getDefaultBlockData(type)
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

	const moveBlock = (index: number, direction: 'up' | 'down') => {
		const newBlocks = [...blocks];
		const newIndex = direction === 'up' ? index - 1 : index + 1;

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
			case 'hero':
				return { title: '', subtitle: '', cta_text: '', cta_url: '' };
			case 'richtext':
				return { content: '' };
			case 'header':
				return { heading: '', subheading: '', level: 'h2' };
			case 'about':
				return { title: '', description: '', heading: '', subheading: '' };
			default:
				return {};
		}
	};

	const getBlockTypeColor = (type: BlockType) => {
		switch (type) {
			case 'hero':
				return 'bg-blue-100 text-blue-800';
			case 'richtext':
				return 'bg-purple-100 text-purple-800';
			case 'header':
				return 'bg-green-100 text-green-800';
			case 'gallery':
				return 'bg-orange-100 text-orange-800';
			case 'about':
				return 'bg-teal-100 text-teal-800';
			default:
				return 'bg-gray-100 text-gray-800';
		}
	};

	// Close separator menu when clicking outside
	$effect(() => {
		const handleClick = (e: MouseEvent) => {
			const target = e.target as HTMLElement;
			if (!target.closest('.separator-menu') && !target.closest('.separator-trigger')) {
				activeSeparator = null;
			}
		};

		if (activeSeparator !== null) {
			document.addEventListener('click', handleClick);
			return () => document.removeEventListener('click', handleClick);
		}
	});
</script>

{#snippet AddBlockMenu(position: number)}
	<button
		onclick={() => addBlockAt('hero', position)}
		type="button"
		class="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-nowrap text-white transition-colors hover:bg-blue-700"
	>
		+ Hero Block
	</button>
	<button
		onclick={() => addBlockAt('richtext', position)}
		type="button"
		class="rounded-lg bg-purple-600 px-4 py-2 text-sm font-medium text-nowrap text-white transition-colors hover:bg-purple-700"
	>
		+ Rich Text
	</button>
	<button
		onclick={() => addBlockAt('header', position)}
		type="button"
		class="rounded-lg bg-green-600 px-4 py-2 text-sm font-medium text-nowrap text-white transition-colors hover:bg-green-700"
	>
		+ Header
	</button>
	<button
		onclick={() => addBlockAt('gallery', position)}
		type="button"
		class="rounded-lg bg-orange-600 px-4 py-2 text-sm font-medium text-nowrap text-white transition-colors hover:bg-orange-700"
	>
		+ Gallery
	</button>
	<button
		onclick={() => addBlockAt('about', position)}
		type="button"
		class="rounded-lg bg-teal-600 px-4 py-2 text-sm font-medium text-nowrap text-white transition-colors hover:bg-teal-700"
	>
		+ About
	</button>
{/snippet}

<div class="block-editor">
	<div class="mb-6 flex items-center justify-between">
		<h2 class="text-2xl">Content Blocks</h2>
	</div>

	{#if blocks.length === 0}
		<div
			class="rounded-lg border-2 border-dashed border-foreground-muted bg-surface py-12 text-center"
		>
			<svg
				class="mx-auto mb-3 h-12 w-12 text-gray-400"
				fill="none"
				stroke="currentColor"
				viewBox="0 0 24 24"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
				/>
			</svg>
			<p class="mb-4 text-foreground-muted">
				No blocks added yet. Click below to add your first block.
			</p>

			<!-- Initial Add Block Menu -->
			<div class="flex justify-center gap-2">
				{@render AddBlockMenu(0)}
			</div>
		</div>
	{:else}
		<div class="space-y-1">
			<!-- Insert separator at the top -->
			<!-- Insert separator at the top -->
			<div class="relative my-2">
				<button
					type="button"
					class="separator-trigger group flex w-full items-center rounded py-2 transition-colors"
					onclick={(e) => {
						e.stopPropagation();
						toggleSeparator(0);
					}}
				>
					<div
						aria-hidden="true"
						class="w-full border-t border-gray-300 transition-colors group-hover:border-gray-400"
					></div>
					<div class="relative flex justify-center px-3">
						<span
							class=" flex items-center gap-1.5 rounded-full p-px text-gray-500 transition-colors group-hover:bg-gray-900 group-hover:text-white"
						>
							<svg viewBox="0 0 20 20" fill="currentColor" aria-hidden="true" class="h-5 w-5">
								<path
									d="M10.75 4.75a.75.75 0 0 0-1.5 0v4.5h-4.5a.75.75 0 0 0 0 1.5h4.5v4.5a.75.75 0 0 0 1.5 0v-4.5h4.5a.75.75 0 0 0 0-1.5h-4.5v-4.5Z"
								/>
							</svg>
						</span>
					</div>
					<div
						aria-hidden="true"
						class="w-full border-t border-gray-300 transition-colors group-hover:border-gray-400"
					></div>
				</button>
				{#if activeSeparator === 0}
					<div
						class="separator-menu absolute top-full left-1/2 z-10 mt-1 flex -translate-x-1/2 gap-2 rounded-lg border border-gray-200 bg-white p-2 shadow-lg"
					>
						{@render AddBlockMenu(0)}
					</div>
				{/if}
			</div>
			{#each blocks as block, index (index)}
				<!-- Block Content -->
				<div
					class="rounded-lg border border-gray-700 bg-surface p-8 shadow-sm transition-shadow hover:shadow-md"
				>
					<div class="mb-4 flex items-start justify-between">
						<span
							class="inline-flex items-center rounded-full px-3 py-1 text-xs font-medium capitalize {getBlockTypeColor(
								block.type
							)}"
						>
							{block.type}
						</span>
						<div class="flex gap-2">
							<button
								onclick={() => moveBlock(index, 'up')}
								disabled={index === 0}
								type="button"
								class="p-1 text-gray-600 transition-colors hover:text-gray-900 disabled:cursor-not-allowed disabled:opacity-30"
								aria-label="Move block up"
							>
								<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M5 15l7-7 7 7"
									/>
								</svg>
							</button>
							<button
								onclick={() => moveBlock(index, 'down')}
								disabled={index === blocks.length - 1}
								type="button"
								class="p-1 text-gray-600 transition-colors hover:text-gray-900 disabled:cursor-not-allowed disabled:opacity-30"
								aria-label="Move block down"
							>
								<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M19 9l-7 7-7-7"
									/>
								</svg>
							</button>
							<button
								onclick={() => removeBlock(index)}
								type="button"
								class="p-1 text-red-600 transition-colors hover:text-red-900"
								aria-label="Remove block"
							>
								<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M6 18L18 6M6 6l12 12"
									/>
								</svg>
							</button>
						</div>
					</div>

					{#if block.type === 'hero'}
						<HeroBlockForm
							data={block.data as HeroBlockData}
							onchange={(data) => updateBlockData(index, data)}
						/>
					{:else if block.type === 'richtext'}
						<RichTextBlockForm
							data={block.data as RichtextBlockData}
							onchange={(data) => updateBlockData(index, data)}
						/>
					{:else if block.type === 'header'}
						<HeaderBlockForm
							data={block.data as HeaderBlockData}
							onchange={(data) => updateBlockData(index, data)}
						/>
					{:else if block.type === 'gallery'}
						<GalleryBlockForm
							data={block.data as GalleryBlockData}
							onchange={(data) => updateBlockData(index, data)}
						/>
					{:else if block.type === 'about'}
						<AboutBlockForm
							data={block.data as AboutBlockData}
							onchange={(data) => updateBlockData(index, data)}
						/>
					{/if}
				</div>

				<!-- Insert separator after each block -->
				<div class="relative my-2">
					<button
						type="button"
						class="separator-trigger group flex w-full items-center rounded py-2 transition-colors"
						onclick={(e) => {
							e.stopPropagation();
							toggleSeparator(index + 1);
						}}
					>
						<div
							aria-hidden="true"
							class="w-full border-t border-gray-300 transition-colors group-hover:border-gray-400"
						></div>
						<div class="relative flex justify-center px-3">
							<span
								class=" flex items-center gap-1.5 rounded-full p-px text-gray-500 transition-colors group-hover:bg-gray-900 group-hover:text-white"
							>
								<svg viewBox="0 0 20 20" fill="currentColor" aria-hidden="true" class="h-5 w-5">
									<path
										d="M10.75 4.75a.75.75 0 0 0-1.5 0v4.5h-4.5a.75.75 0 0 0 0 1.5h4.5v4.5a.75.75 0 0 0 1.5 0v-4.5h4.5a.75.75 0 0 0 0-1.5h-4.5v-4.5Z"
									/>
								</svg>
							</span>
						</div>
						<div
							aria-hidden="true"
							class="w-full border-t border-gray-300 transition-colors group-hover:border-gray-400"
						></div>
					</button>
					{#if activeSeparator === index + 1}
						<div
							class="separator-menu absolute top-full left-1/2 z-10 mt-1 flex -translate-x-1/2 gap-2 rounded-lg border border-gray-200 bg-white p-2 shadow-lg"
						>
							{@render AddBlockMenu(index + 1)}
						</div>
					{/if}
				</div>
			{/each}
		</div>
	{/if}
</div>
