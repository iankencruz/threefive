<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { Editor } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import Underline from '@tiptap/extension-underline';
	import Image from '@tiptap/extension-image';
	import { toast } from 'svelte-sonner';
	import { X } from '@lucide/svelte';
	// Import your media API and types
	// import { fetchMedia } from '$lib/api/media';
	// import type { MediaItem } from '$lib/types';
	// import Pagination from '../Navigation/Pagination.svelte';
	// import UploadMediaForm from '../Forms/UploadMediaForm.svelte';

	// Mock types for demo - replace with your actual imports
	type MediaItem = {
		id: string;
		title: string;
		url: string;
		thumbnail_url?: string;
	};

	let element = $state<HTMLElement>();
	let editor = $state<Editor>();
	let ui = $state(0); // bumps to force a re-render

	// Media popup state
	let showMediaPopup = $state(false);
	let tab = $state<'link' | 'upload'>('link');
	let loading = $state(false);
	let page = $state(1);
	let totalPages = $state(1);
	let pageSize = 10;
	let totalMedia = $state(0);
	let allUnlinkedMedia: MediaItem[] = [];
	let filteredMedia = $state<MediaItem[]>([]);
	let hasLoaded = false;

	let { body = $bindable(), context = { type: 'gallery', id: 'default' } } = $props<{
		body?: string;
		context?: { type: 'project' | 'gallery' | 'block' | 'page'; id: string };
	}>();

	let lastBody: string | null = null; // track last applied external body

	const refreshUI = () => {
		ui++;
	};

	// Mock media data - replace with your actual fetchMedia function
	const mockMedia: MediaItem[] = [
		{
			id: '1',
			title: 'Sample Image 1',
			url: 'https://picsum.photos/400/300?random=1',
			thumbnail_url: 'https://picsum.photos/200/150?random=1'
		},
		{
			id: '2',
			title: 'Sample Image 2',
			url: 'https://picsum.photos/400/300?random=2',
			thumbnail_url: 'https://picsum.photos/200/150?random=2'
		},
		{
			id: '3',
			title: 'Sample Image 3',
			url: 'https://picsum.photos/400/300?random=3',
			thumbnail_url: 'https://picsum.photos/200/150?random=3'
		},
		{
			id: '4',
			title: 'Sample Image 4',
			url: 'https://picsum.photos/400/300?random=4',
			thumbnail_url: 'https://picsum.photos/200/150?random=4'
		},
		{
			id: '5',
			title: 'Sample Image 5',
			url: 'https://picsum.photos/400/300?random=5',
			thumbnail_url: 'https://picsum.photos/200/150?random=5'
		},
		{
			id: '6',
			title: 'Sample Image 6',
			url: 'https://picsum.photos/400/300?random=6',
			thumbnail_url: 'https://picsum.photos/200/150?random=6'
		}
	];

	function paginateMediaPool() {
		const start = (page - 1) * pageSize;
		const end = page * pageSize;

		filteredMedia = allUnlinkedMedia.slice(start, end);
		totalMedia = allUnlinkedMedia.length;
		totalPages = Math.max(1, Math.ceil(totalMedia / pageSize));
	}

	async function loadMedia() {
		loading = true;
		try {
			// Replace this with your actual fetchMedia call
			// const res = await fetchMedia(1, 1000);
			// allUnlinkedMedia = res.items;

			// Mock implementation
			await new Promise((resolve) => setTimeout(resolve, 500)); // Simulate loading
			allUnlinkedMedia = [...mockMedia];
			paginateMediaPool();
		} catch (err) {
			console.error('❌ Failed to load media:', err);
			// toast.error('Could not load media');
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		editor = new Editor({
			element,
			extensions: [StarterKit.configure({ heading: { levels: [1, 2, 3] } }), Underline, Image],
			editorProps: {
				attributes: { class: 'prose prose-sm sm:prose-base lg:prose-lg m-5 focus:outline-none' }
			},
			content: body,
			onTransaction: () => {
				editor = editor;
			}
		});

		editor.on('selectionUpdate', refreshUI);
		editor.on('transaction', refreshUI);

		editor.on('update', () => {
			const html = editor?.getHTML();
			if (html !== body) {
				body = html; // push up to parent
				lastBody = html || ''; // mark as synced
			}
			refreshUI();
		});
	});

	onDestroy(() => {
		if (editor) {
			editor.destroy();
			editor = undefined;
		}
	});

	$effect(() => {
		if (!editor) return;
		if (body !== lastBody) {
			editor.commands.setContent(body, false);
			lastBody = body;
		}
	});

	$effect(() => {
		if (showMediaPopup && !hasLoaded) {
			tab = 'link';
			page = 1;
			loadMedia();
			hasLoaded = true;
		}
	});

	$effect(() => {
		if (!showMediaPopup) hasLoaded = false;
	});

	function run(cmd: () => boolean | undefined) {
		cmd?.();
		refreshUI();
	}

	const BLOCK_OPTIONS = [
		{ value: 'paragraph', label: 'Paragraph' },
		{ value: 'h1', label: 'Heading 1' },
		{ value: 'h2', label: 'Heading 2' },
		{ value: 'h3', label: 'Heading 3' }
	] as const;

	function currentBlockValue(): string {
		if (!editor) return 'paragraph';
		if (editor.isActive('heading', { level: 1 })) return 'h1';
		if (editor.isActive('heading', { level: 2 })) return 'h2';
		if (editor.isActive('heading', { level: 3 })) return 'h3';
		return 'paragraph';
	}

	function openMediaPopup() {
		showMediaPopup = true;
	}

	function closeMediaPopup() {
		showMediaPopup = false;
	}

	function handlePageChange(newPage: number) {
		if (newPage >= 1 && newPage <= totalPages && newPage !== page) {
			page = newPage;
			paginateMediaPool();
		}
	}

	function linkMedia(item: MediaItem) {
		const chain = editor?.chain().focus();
		chain
			?.setImage({
				src: item.url,
				alt: item.title || undefined
			})
			.run();
		closeMediaPopup();
		toast.success('Image inserted');
	}

	function handleMediaPopupKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') {
			closeMediaPopup();
		}
	}

	function applyBlock(value: string) {
		const chain = editor?.chain().focus();
		switch (value) {
			case 'paragraph':
				run(() => chain?.setParagraph().run());
				break;
			case 'h1':
				run(() => chain?.toggleHeading({ level: 1 }).run());
				break;
			case 'h2':
				run(() => chain?.toggleHeading({ level: 2 }).run());
				break;
			case 'h3':
				run(() => chain?.toggleHeading({ level: 3 }).run());
				break;
		}
	}
</script>

{#if editor}
	<!-- data-ui={ui} forces Svelte to re-evaluate bindings when selection changes -->
	<div
		class="rte-toolbar flex flex-wrap items-center gap-2 rounded-t-md border bg-gray-50 p-2 shadow-sm"
		data-ui={ui}
	>
		<!-- Paragraph/H1/H2/H3 dropdown -->
		<select
			class="tiptap h-8 rounded border bg-white px-2 text-sm"
			onchange={(e) => applyBlock((e.target as HTMLSelectElement).value)}
			value={currentBlockValue()}
			title="Block format"
		>
			{#each BLOCK_OPTIONS as opt (opt.value)}
				<option value={opt.value}>{opt.label}</option>
			{/each}
		</select>

		<div class="mx-1 h-6 w-px bg-gray-300" />

		<!-- Inline styles -->
		<div class="inline-flex overflow-hidden rounded border bg-white">
			<button
				class="tiptap"
				data-active={editor?.isActive('bold') ?? false}
				title="Bold"
				onclick={() => run(() => editor?.chain().focus().toggleBold().run())}
				disabled={(editor?.can().chain().focus().toggleBold().run() ?? false) === false}
				><strong>B</strong></button
			>

			<button
				class="tiptap"
				data-active={editor?.isActive('italic') ?? false}
				title="Italic"
				onclick={() => run(() => editor?.chain().focus().toggleItalic().run())}
				disabled={(editor?.can().chain().focus().toggleItalic().run() ?? false) === false}
				><em>I</em></button
			>

			<button
				class="tiptap"
				data-active={editor?.isActive('underline') ?? false}
				title="Underline"
				onclick={() => run(() => editor?.chain().focus().toggleUnderline().run())}
				disabled={(editor?.can().chain().focus().toggleUnderline().run() ?? false) === false}
				><u>U</u></button
			>

			<button
				class="tiptap"
				data-active={editor?.isActive('strike') ?? false}
				title="Strikethrough"
				onclick={() => run(() => editor?.chain().focus().toggleStrike().run())}
				disabled={(editor?.can().chain().focus().toggleStrike().run() ?? false) === false}
				><s>S</s></button
			>
		</div>

		<!-- Lists -->
		<div class="inline-flex overflow-hidden rounded border bg-white">
			<button
				class="tiptap"
				data-active={editor?.isActive('bulletList') ?? false}
				title="Bullet list"
				onclick={() => run(() => editor?.chain().focus().toggleBulletList().run())}>• List</button
			>

			<button
				class="tiptap"
				data-active={editor?.isActive('orderedList') ?? false}
				title="Numbered list"
				onclick={() => run(() => editor?.chain().focus().toggleOrderedList().run())}>1. List</button
			>
		</div>

		<!-- Blockquote & Code block -->
		<div class="inline-flex overflow-hidden rounded border bg-white">
			<button
				class="tiptap"
				data-active={editor?.isActive('blockquote') ?? false}
				title="Blockquote"
				onclick={() => run(() => editor?.chain().focus().toggleBlockquote().run())}
				>&ldquo;Quote&rdquo;</button
			>

			<button
				class="tiptap"
				data-active={editor?.isActive('codeBlock') ?? false}
				title="Code block"
				onclick={() => run(() => editor?.chain().focus().toggleCodeBlock().run())}>[code]</button
			>
		</div>

		<!-- Insert -->
		<div class="inline-flex overflow-hidden rounded border bg-white">
			<button
				class="tiptap"
				title="Horizontal rule"
				onclick={() => run(() => editor?.chain().focus().setHorizontalRule().run())}>—</button
			>
		</div>

		<!-- Image -->
		<div class="inline-flex overflow-hidden rounded border bg-white">
			<button
				class="tiptap"
				data-active={editor?.isActive('image') ?? false}
				title="Image"
				onclick={openMediaPopup}
			>
				📷 Image</button
			>
		</div>

		<!-- Undo / Redo -->
		<div class="inline-flex overflow-hidden rounded border bg-white">
			<button
				class="tiptap"
				title="Undo"
				onclick={() => run(() => editor?.chain().focus().undo().run())}
				disabled={(editor?.can().chain().focus().undo().run() ?? false) === false}>⤺ Undo</button
			>

			<button
				class="tiptap"
				title="Redo"
				onclick={() => run(() => editor?.chain().focus().redo().run())}
				disabled={(editor?.can().chain().focus().redo().run() ?? false) === false}>⤻ Redo</button
			>
		</div>
	</div>
{/if}

<div class="rounded-b-md border border-t-0 p-2" bind:this={element} />

<!-- Media Popup Modal - Similar to LinkedMediaModal -->
{#if showMediaPopup}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
		onkeydown={handleMediaPopupKeydown}
		tabindex="-1"
	>
		<div class="relative w-full max-w-4xl rounded-lg bg-white shadow-xl">
			<!-- Header -->
			<div class="flex items-center justify-between border-b p-4">
				<h3 class="text-lg font-semibold text-gray-900">Select Media</h3>
				<button
					onclick={closeMediaPopup}
					class="rounded-md p-2 text-gray-400 hover:bg-gray-100 hover:text-gray-600"
					title="Close"
				>
					<X size={20} />
				</button>
			</div>

			<div class="p-4">
				<!-- Tabs -->
				<div class="mt-4 flex gap-2 border-b text-sm font-medium">
					<button
						onclick={() => (tab = 'link')}
						class="rounded-t px-4 py-2"
						class:font-bold={tab === 'link'}
						class:border-b-2={tab === 'link'}
						class:border-blue-500={tab === 'link'}
					>
						Link Media
					</button>
					<button
						onclick={() => (tab = 'upload')}
						class="rounded-t px-4 py-2"
						class:font-bold={tab === 'upload'}
						class:border-b-2={tab === 'upload'}
						class:border-blue-500={tab === 'upload'}
					>
						Upload Media
					</button>
				</div>

				<!-- Tab Content -->
				{#if tab === 'link'}
					{#if loading}
						<div class="flex items-center justify-center p-8">
							<div class="h-8 w-8 animate-spin rounded-full border-b-2 border-blue-600"></div>
							<span class="ml-2 text-gray-500">Loading media...</span>
						</div>
					{:else if filteredMedia.length === 0}
						<p class="p-8 text-center text-gray-500">No media available to link.</p>
					{:else}
						<div class="mt-4">
							<ul
								class="grid grid-cols-2 gap-x-4 gap-y-8 sm:grid-cols-3 md:grid-cols-4 xl:grid-cols-6 xl:gap-x-8"
							>
								{#each filteredMedia as item (item.id)}
									<li
										class="group relative block overflow-hidden rounded-lg bg-gray-100 ring-1 ring-gray-200 transition-all hover:cursor-pointer hover:ring-2 hover:ring-blue-500"
									>
										<button
											class="h-auto w-full text-left"
											onclick={() => linkMedia(item)}
											title={`Insert ${item.title || 'Untitled'}`}
										>
											<img
												src={item.thumbnail_url || item.url}
												alt={item.title}
												class="aspect-video w-full object-cover transition-opacity group-hover:opacity-75"
											/>
											<div class="p-2">
												<p class="truncate text-sm font-medium text-gray-900">
													{item.title || 'Untitled'}
												</p>
											</div>
										</button>
									</li>
								{/each}
							</ul>

							<!-- Simple Pagination -->
							{#if totalPages > 1}
								<div class="mt-6 flex items-center justify-between">
									<div class="text-sm text-gray-700">
										Showing {Math.min((page - 1) * pageSize + 1, totalMedia)} to {Math.min(
											page * pageSize,
											totalMedia
										)} of {totalMedia} results
									</div>
									<div class="flex gap-2">
										<button
											onclick={() => handlePageChange(page - 1)}
											disabled={page <= 1}
											class="rounded border px-3 py-1 text-sm hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50"
										>
											Previous
										</button>
										<span class="px-3 py-1 text-sm font-medium">
											Page {page} of {totalPages}
										</span>
										<button
											onclick={() => handlePageChange(page + 1)}
											disabled={page >= totalPages}
											class="rounded border px-3 py-1 text-sm hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50"
										>
											Next
										</button>
									</div>
								</div>
							{/if}
						</div>
					{/if}
				{:else if tab === 'upload'}
					<div class="mt-4 rounded-lg border-2 border-dashed border-gray-300 p-8 text-center">
						<p class="text-gray-500">Upload functionality would go here</p>
						<p class="mt-2 text-sm text-gray-400">
							Replace this with your UploadMediaForm component
						</p>
					</div>
				{/if}
			</div>
		</div>
	</div>
{/if}
