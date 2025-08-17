<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { Editor } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import Underline from '@tiptap/extension-underline';

	let element = $state<HTMLElement>();
	let editor = $state<Editor>();

	onMount(() => {
		editor = new Editor({
			element,
			extensions: [
				StarterKit.configure({
					heading: { levels: [1, 2, 3] }
				}),
				Underline
			],
			editorProps: {
				attributes: {
					class: 'prose prose-sm sm:prose-base lg:prose-lg m-5 focus:outline-none '
				}
			},
			content: '<p>Hello World! 🌍️</p>',
			onTransaction: () => {
				// trigger reactive update so isActive() updates in Svelte
				editor = editor;
			}
		});
	});

	onDestroy(() => editor?.destroy());

	function toggle(cmd: () => void) {
		cmd();
		editor = editor; // force refresh state
	}
</script>

{#if editor}
	<div class="flex flex-wrap gap-1 border border-b bg-gray-50 p-2">
		<!-- Bold -->
		<button
			onclick={() => toggle(() => editor?.chain().focus().toggleBold().run())}
			class="tiptap"
			class:active={editor.isActive('bold')}><strong>B</strong></button
		>

		<!-- Italic -->
		<button
			onclick={() => toggle(() => editor?.chain().focus().toggleItalic().run())}
			class="tiptap"
			class:active={editor.isActive('italic')}><em>I</em></button
		>

		<!-- Underline -->
		<button
			onclick={() => toggle(() => editor?.chain().focus().toggleUnderline().run())}
			class="tiptap"
			class:active={editor.isActive('underline')}><u>U</u></button
		>

		<!-- Strike -->
		<button
			onclick={() => toggle(() => editor?.chain().focus().toggleStrike().run())}
			class="tiptap"
			class:active={editor.isActive('strike')}><s>S</s></button
		>

		<!-- Headings -->
		<button
			onclick={() => toggle(() => editor?.chain().focus().toggleHeading({ level: 1 }).run())}
			class="tiptap"
			class:active={editor.isActive('heading', { level: 1 })}>H1</button
		>
		<button
			onclick={() => toggle(() => editor?.chain().focus().toggleHeading({ level: 2 }).run())}
			class="tiptap"
			class:active={editor.isActive('heading', { level: 2 })}>H2</button
		>
		<button
			onclick={() => toggle(() => editor?.chain().focus().toggleHeading({ level: 3 }).run())}
			class="tiptap"
			class:active={editor.isActive('heading', { level: 3 })}>H3</button
		>

		<!-- Paragraph -->
		<button
			onclick={() => toggle(() => editor?.chain().focus().setParagraph().run())}
			class="tiptap"
			class:active={editor.isActive('paragraph')}>P</button
		>

		<!-- Lists -->
		<button
			onclick={() => toggle(() => editor?.chain().focus().toggleBulletList().run())}
			class="tiptap"
			class:active={editor.isActive('bulletList')}>• List</button
		>
		<button
			onclick={() => toggle(() => editor?.chain().focus().toggleOrderedList().run())}
			class="tiptap"
			class:active={editor.isActive('orderedList')}>1. List</button
		>

		<!-- Blockquote -->
		<button
			onclick={() => toggle(() => editor?.chain().focus().toggleBlockquote().run())}
			class="tiptap"
			class:active={editor.isActive('blockquote')}>&ldquo;Quote&rdquo;</button
		>

		<!-- Code -->
		<button
			onclick={() => toggle(() => editor?.chain().focus().toggleCode().run())}
			class="tiptap"
			class:active={editor.isActive('code')}>&lt;/&gt;</button
		>

		<!-- Code Block -->
		<button
			onclick={() => toggle(() => editor?.chain().focus().toggleCodeBlock().run())}
			class="tiptap"
			class:active={editor.isActive('codeBlock')}>[code]</button
		>

		<!-- Horizontal Rule -->
		<button onclick={() => editor?.chain().focus().setHorizontalRule().run()}>―</button>

		<!-- Undo / Redo -->
		<button onclick={() => editor?.chain().focus().undo().run()}>⤺ Undo</button>
		<button onclick={() => editor?.chain().focus().redo().run()}>⤻ Redo</button>
	</div>
{/if}

<div bind:this={element} />
