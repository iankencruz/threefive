<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { Editor } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import Underline from '@tiptap/extension-underline';

	let element = $state<HTMLElement>();
	let editor = $state<Editor>();
	let ui = $state(0); // bumps to force a re-render

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
			extensions: [StarterKit.configure({ heading: { levels: [1, 2, 3] } }), Underline],
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
			{#each BLOCK_OPTIONS as opt}
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
