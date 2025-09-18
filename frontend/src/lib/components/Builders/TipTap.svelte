<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { Editor } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import Underline from '@tiptap/extension-underline';
	import Image from '@tiptap/extension-image';

	let element = $state<HTMLElement>();
	let editor = $state<Editor>();
	let ui = $state(0); // bumps to force a re-render

	// Media popup state
	let showMediaPopup = $state(false);
	let mediaUrl = $state('');
	let mediaAlt = $state('');

	let { body = $bindable() } = $props();

	let lastBody: string | null = null; // track last applied external body

	const refreshUI = () => {
		ui++;
	};

	// flag to prevent infinite update loop
	let isUpdating = false;

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

	// ✅ Cleanup editor when component unmounts
	onDestroy(() => {
		if (editor) {
			editor.destroy();
			editor = undefined;
		}
	});

	// 🔑 Only update editor if body changed externally
	$effect(() => {
		if (!editor) return;
		if (body !== lastBody) {
			editor.commands.setContent(body, false);
			lastBody = body;
		}
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
		mediaUrl = '';
		mediaAlt = '';
		showMediaPopup = true;
	}

	function closeMediaPopup() {
		showMediaPopup = false;
		mediaUrl = '';
		mediaAlt = '';
	}

	function insertImage() {
		if (mediaUrl.trim()) {
			const chain = editor?.chain().focus();
			chain
				?.setImage({
					src: mediaUrl.trim(),
					alt: mediaAlt.trim() || undefined
				})
				.run();
		}
		closeMediaPopup();
	}

	function handleMediaPopupKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') {
			closeMediaPopup();
		} else if (e.key === 'Enter' && e.ctrlKey) {
			insertImage();
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

<!-- Media Popup Modal -->
{#if showMediaPopup}
	<!-- Overlay -->
	<div
		class="bg-opacity-50 fixed inset-0 z-50 flex items-center justify-center bg-black p-4"
		onclick={closeMediaPopup}
		onkeydown={handleMediaPopupKeydown}
		tabindex="-1"
	>
		<!-- Modal Content -->
		<div
			class="w-full max-w-md rounded-lg bg-white p-6 shadow-xl"
			onclick={(e) => e.stopPropagation()}
		>
			<div class="mb-4 flex items-center justify-between">
				<h3 class="text-lg font-semibold text-gray-900">Insert Image</h3>
				<button
					onclick={closeMediaPopup}
					class="text-xl font-semibold text-gray-400 hover:text-gray-600"
					title="Close"
				>
					×
				</button>
			</div>

			<div class="space-y-4">
				<div>
					<label for="media-url" class="mb-1 block text-sm font-medium text-gray-700">
						Image URL *
					</label>
					<input
						id="media-url"
						type="url"
						bind:value={mediaUrl}
						placeholder="https://example.com/image.jpg"
						class="w-full rounded-md border border-gray-300 px-3 py-2 focus:border-transparent focus:ring-2 focus:ring-blue-500 focus:outline-none"
						autofocus
					/>
				</div>

				<div>
					<label for="media-alt" class="mb-1 block text-sm font-medium text-gray-700">
						Alt Text (optional)
					</label>
					<input
						id="media-alt"
						type="text"
						bind:value={mediaAlt}
						placeholder="Description of the image"
						class="w-full rounded-md border border-gray-300 px-3 py-2 focus:border-transparent focus:ring-2 focus:ring-blue-500 focus:outline-none"
					/>
				</div>

				{#if mediaUrl.trim()}
					<div>
						<label class="mb-2 block text-sm font-medium text-gray-700">Preview:</label>
						<div class="rounded-md border border-gray-200 bg-gray-50 p-2">
							<img
								src={mediaUrl}
								alt={mediaAlt || 'Image preview'}
								class="mx-auto h-auto max-h-32 max-w-full rounded"
								onerror={(e) => {
									(e.target as HTMLImageElement).style.display = 'none';
									(e.target as HTMLImageElement).nextElementSibling?.remove();
									const errorMsg = document.createElement('div');
									errorMsg.textContent = 'Invalid image URL';
									errorMsg.className = 'text-red-500 text-sm text-center';
									(e.target as HTMLImageElement).parentNode?.appendChild(errorMsg);
								}}
							/>
						</div>
					</div>
				{/if}
			</div>

			<div class="mt-6 flex justify-end gap-3">
				<button
					onclick={closeMediaPopup}
					class="rounded-md border border-gray-300 px-4 py-2 text-gray-600 hover:bg-gray-50 focus:ring-2 focus:ring-gray-500 focus:outline-none"
				>
					Cancel
				</button>
				<button
					onclick={insertImage}
					disabled={!mediaUrl.trim()}
					class="rounded-md bg-blue-600 px-4 py-2 text-white hover:bg-blue-700 focus:ring-2 focus:ring-blue-500 focus:outline-none disabled:cursor-not-allowed disabled:bg-gray-300"
				>
					Insert Image
				</button>
			</div>
		</div>
	</div>
{/if}
