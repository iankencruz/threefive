<!-- frontend/src/lib/components/ui/RichTextEditor.svelte -->
<script lang="ts">
	import { onMount, onDestroy } from "svelte";
	import { Editor } from "@tiptap/core";
	import StarterKit from "@tiptap/starter-kit";
	import Underline from "@tiptap/extension-underline";
	import Image from "@tiptap/extension-image";
	import MediaPicker from "./MediaPicker.svelte";
	import { PUBLIC_API_URL } from "$env/static/public";
	import type { Media } from "$api/media";
	import {
		CodeXml,
		ImageIcon,
		List,
		ListOrdered,
		Minus,
		Quote,
		Undo2,
		Redo2,
	} from "lucide-svelte";

	interface Props {
		value?: string;
		placeholder?: string;
		class?: string;
	}

	let {
		value = $bindable(""),
		placeholder = "Start typing...",
		class: className = "",
	}: Props = $props();

	let element = $state<HTMLElement>();
	let editor = $state<Editor>();
	let ui = $state(0);

	// Media picker state
	let showMediaPicker = $state(false);

	let lastValue: string | null = null;

	const refreshUI = () => {
		ui++;
	};

	onMount(() => {
		if (!element) return;

		editor = new Editor({
			element,
			extensions: [
				StarterKit.configure({
					heading: { levels: [1, 2, 3] },
				}),
				Underline,
				Image,
			],
			editorProps: {
				attributes: {
					class:
						"prose prose-neutral prose-sm sm:prose-base prose-headings:text-white prose-ol:text-white prose-ul:text-white  vmax-w-none p-4 focus:outline-none min-h-[200px]",
				},
			},
			content: value || "",
			onTransaction: () => {
				editor = editor;
			},
		});

		editor.on("selectionUpdate", refreshUI);
		editor.on("transaction", refreshUI);

		editor.on("update", () => {
			const html = editor?.getHTML();
			if (html !== value) {
				value = html || "";
				lastValue = html || "";
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
		if (value !== lastValue) {
			editor.commands.setContent(value || "", false);
			lastValue = value || "";
		}
	});

	function run(cmd: () => boolean | undefined) {
		cmd?.();
		refreshUI();
	}

	const BLOCK_OPTIONS = [
		{ value: "paragraph", label: "Paragraph" },
		{ value: "h1", label: "Heading 1" },
		{ value: "h2", label: "Heading 2" },
		{ value: "h3", label: "Heading 3" },
	] as const;

	function currentBlockValue(): string {
		if (!editor) return "paragraph";
		if (editor.isActive("heading", { level: 1 })) return "h1";
		if (editor.isActive("heading", { level: 2 })) return "h2";
		if (editor.isActive("heading", { level: 3 })) return "h3";
		return "paragraph";
	}

	function openMediaPicker() {
		showMediaPicker = true;
	}

	function closeMediaPicker() {
		showMediaPicker = false;
	}

	// Handle media selection from MediaPicker
	async function handleMediaSelect(mediaId: string, media: Media) {
		if (!editor) return;

		// Get the appropriate URL - prefer medium_url for content, fallback to large_url or url
		const imageUrl = media.medium_url || media.large_url || media.url;

		// Insert image into editor
		if (imageUrl) {
			editor
				.chain()
				.focus()
				.setImage({
					src: imageUrl,
					alt: media.original_filename || "",
				})
				.run();
		}
	}

	function applyBlock(val: string) {
		if (!editor) return;
		const chain = editor.chain().focus();
		switch (val) {
			case "paragraph":
				run(() => chain.setParagraph().run());
				break;
			case "h1":
				run(() => chain.toggleHeading({ level: 1 }).run());
				break;
			case "h2":
				run(() => chain.toggleHeading({ level: 2 }).run());
				break;
			case "h3":
				run(() => chain.toggleHeading({ level: 3 }).run());
				break;
		}
	}
</script>

<div class="border border-gray-700  rounded-lg overflow-hidden bg-surface {className}">
	{#if editor}
		<div class="flex flex-wrap items-center gap-2 border-b border-gray-700 bg-surface p-2" data-ui={ui}>
			<select
				class="h-8 rounded border  border-input-border min-w-44 bg-surface px-2 text-sm hover:bg-black/75 focus:outline-none focus:ring-2 focus:ring-blue-500"
				onchange={(e) => applyBlock((e.target as HTMLSelectElement).value)}
				value={currentBlockValue()}
				title="Block format"
			>
				{#each BLOCK_OPTIONS as opt (opt.value)}
					<option value={opt.value}>{opt.label}</option>
				{/each}
			</select>

			<div class="h-6 w-px bg-gray-300"></div>

			<div class="inline-flex overflow-hidden rounded border bg-surface border-input-border">
				<button 
          type="button" 
          class="px-3 py-1 hover:bg-black/75 disabled:opacity-50 disabled:cursor-not-allowed transition-colors" 
          class:bg-blue-100={editor?.isActive('bold')} 
          class:text-blue-600={editor?.isActive('bold')} 
          title="Bold (Ctrl+B)"
          onclick={() => run(() => editor?.chain().focus().toggleBold().run())} 
          disabled={!editor?.can().chain().focus().toggleBold().run()} 
        >
          <strong>B</strong>
        </button>
				<button 
          type="button" 
          class="px-3 py-1 hover:bg-black/75 disabled:opacity-50 disabled:cursor-not-allowed transition-colors border-l" 
          class:bg-blue-100={editor?.isActive('italic')} 
          class:text-blue-600={editor?.isActive('italic')} 
          title="Italic (Ctrl+I)" 
          onclick={() => run(() => editor?.chain().focus().toggleItalic().run())} disabled={!editor?.can().chain().focus().toggleItalic().run()}
        >
          <em>I</em>
        </button>
				<button 
          type="button" 
          class="px-3 py-1 hover:bg-black/75 disabled:opacity-50 disabled:cursor-not-allowed transition-colors border-l" 
          class:bg-blue-100={editor?.isActive('underline')} 
          class:text-blue-600={editor?.isActive('underline')} 
          title="Underline (Ctrl+U)" 
          onclick={() => run(() => editor?.chain().focus().toggleUnderline().run())} 
          disabled={!editor?.can().chain().focus().toggleUnderline().run()} 
        >
          <u>U</u>
        </button>
				<button 
          type="button" 
          class="px-3 py-1 hover:bg-black/75 disabled:opacity-50 disabled:cursor-not-allowed transition-colors border-l"
          class:bg-blue-100={editor?.isActive('strike')} 
          class:text-blue-600={editor?.isActive('strike')} 
          title="Strikethrough" 
          onclick={() => run(() => editor?.chain().focus().toggleStrike().run())} 
          disabled={!editor?.can().chain().focus().toggleStrike().run()} 
        >
          <s>S</s>
        </button>
			</div>

			<div 
        class="inline-flex overflow-hidden rounded border border-input-border bg-surface"
      >
				<button 
          type="button" 
          class="px-3 py-1 hover:bg-black/75 transition-colors" 
          class:bg-blue-100={editor?.isActive('bulletList')} 
          class:text-blue-600={editor?.isActive('bulletList')} 
          title="Bullet list" 
          onclick={() => run(() => editor?.chain().focus().toggleBulletList().run())} 
        >
          <span class="text-sm"><List size={18} /></span>
        </button>
				<button 
          type="button" 
          class="px-3 py-1 hover:bg-black/75 transition-colors border-l" 
          class:bg-blue-100={editor?.isActive('orderedList')} 
          class:text-blue-600={editor?.isActive('orderedList')} 
          title="Numbered list" 
          onclick={() => run(() => editor?.chain().focus().toggleOrderedList().run())} 
        >
         <span class="text-sm"><ListOrdered size={18} /></span>
        </button>
			</div>

			<div class="inline-flex overflow-hidden rounded border border-input-border bg-surface">
				<button 
          type="button" 
          class="px-3 py-1 hover:bg-black/75 transition-colors" 
          class:bg-blue-100={editor?.isActive('blockquote')} 
          class:text-blue-600={editor?.isActive('blockquote')} 
          title="Blockquote" 
          onclick={() => run(() => editor?.chain().focus().toggleBlockquote().run())} 
        >
          <span class="text-sm"><Quote size={20}/></span>
        </button>
				<button 
          type="button" 
          class="px-3 py-1 hover:bg-black/75 transition-colors border-l" 
          class:bg-blue-100={editor?.isActive('codeBlock')} 
          class:text-blue-600={editor?.isActive('codeBlock')} 
          title="Code block" 
          onclick={() => run(() => editor?.chain().focus().toggleCodeBlock().run())} 
        >
          <span class="text-sm"><CodeXml size={20}/></span>
        </button>
			</div>

			<div class="inline-flex overflow-hidden rounded border border-input-border bg-surface">
				<button 
          type="button" 
          class="px-3 py-1 hover:bg-black/75 transition-colors" 
          title="Horizontal rule" 
          onclick={() => run(() => editor?.chain().focus().setHorizontalRule().run())} 
        >
          <span class="text-sm"><Minus size={18} /></span>
        </button>
				<button 
          type="button" 
          class="px-3 py-1 hover:bg-black/75 transition-colors border-l" 
          class:bg-blue-100={editor?.isActive('image')} 
          class:text-blue-600={editor?.isActive('image')} 
          title="Insert image" 
          onclick={() => (
          showMediaPicker = true
          )} 
        >
          <span class="text-sm"><ImageIcon size={18} /></span>
        </button>
			</div>

			<div class="inline-flex overflow-hidden rounded border border-input-border bg-surface">
				<button 
          type="button" 
          class="px-3 py-1 hover:bg-black/75 disabled:opacity-50 disabled:cursor-not-allowed transition-colors" 
          title="Undo (Ctrl+Z)" 
          onclick={() => run(() => editor?.chain().focus().undo().run())} 
          disabled={!editor?.can().chain().focus().undo().run()} 
        >
          <span class="text-sm"><Undo2 size={18} /></span>
        </button>
				<button 
          type="button" 
          class="px-3 py-1 hover:bg-black/75 disabled:opacity-50 disabled:cursor-not-allowed transition-colors border-l" 
          title="Redo (Ctrl+Shift+Z)" 
          onclick={() => run(() => editor?.chain().focus().redo().run())} 
          disabled={!editor?.can().chain().focus().redo().run()} 
        >
          <span class="text-sm"><Redo2 size={18}/></span>
        </button>
			</div>
		</div>
	{/if}

	<div class="bg-surface " bind:this={element}></div>
</div>

<!-- Media Picker Modal - Modal Only Mode -->
<MediaPicker 
	show={showMediaPicker}
	onselect={handleMediaSelect}
	onclose={closeMediaPicker}
/>
