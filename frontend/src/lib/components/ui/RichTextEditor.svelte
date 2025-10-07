<!-- frontend/src/lib/components/ui/RichTextEditor.svelte -->
<script lang="ts">
import { onMount, onDestroy } from "svelte";
import { Editor } from "@tiptap/core";
import StarterKit from "@tiptap/starter-kit";
import Underline from "@tiptap/extension-underline";
import Image from "@tiptap/extension-image";
import { mediaApi, type Media, getMediaUrl } from "$api/media";

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

// Media popup state
let showMediaPopup = $state(false);
let insertMode = $state<"url" | "library">("library");
let mediaUrl = $state("");
let mediaAlt = $state("");

// Media library state
let media = $state<Media[]>([]);
let selectedMediaItem = $state<Media | null>(null);
let loadingMedia = $state(false);
let searchQuery = $state("");
let uploadFile = $state<File | null>(null);
let uploading = $state(false);
let viewMode = $state<"grid" | "list">("grid");
let currentPage = $state(1);
let totalPages = $state(1);
let limit = $state(12);

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
					"prose prose-sm sm:prose-base max-w-none p-4 focus:outline-none min-h-[200px]",
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

async function openMediaPopup() {
	showMediaPopup = true;
	insertMode = "library";
	mediaUrl = "";
	mediaAlt = "";
	selectedMediaItem = null;
	currentPage = 1;
	await loadMedia();
}

function closeMediaPopup() {
	showMediaPopup = false;
	mediaUrl = "";
	mediaAlt = "";
	selectedMediaItem = null;
	uploadFile = null;
}

async function loadMedia() {
	loadingMedia = true;
	try {
		const response = await mediaApi.listMedia(currentPage, limit);
		media = response.data || [];

		if (response.pagination) {
			totalPages = response.pagination.total_pages || 1;
		}
	} catch (err) {
		console.error("Failed to load media:", err);
	} finally {
		loadingMedia = false;
	}
}

async function changePage(newPage: number) {
	if (newPage < 1 || newPage > totalPages) return;
	currentPage = newPage;
	await loadMedia();
}

async function handleFileSelect(e: Event) {
	const input = e.target as HTMLInputElement;
	if (input.files && input.files[0]) {
		uploadFile = input.files[0];
		await handleUpload();
	}
}

async function handleUpload() {
	if (!uploadFile) return;

	uploading = true;
	try {
		const uploaded = await mediaApi.uploadMedia(uploadFile);
		media = [uploaded, ...media];
		selectedMediaItem = uploaded;
		uploadFile = null;
	} catch (err) {
		console.error("Upload failed:", err);
		alert("Failed to upload file");
	} finally {
		uploading = false;
	}
}

function selectMediaFromLibrary(m: Media) {
	selectedMediaItem = m;
	mediaUrl = getMediaUrl(m);
	mediaAlt = m.original_filename;
}

function insertImage() {
	if (editor && mediaUrl.trim()) {
		editor
			.chain()
			.focus()
			.setImage({
				src: mediaUrl.trim(),
				alt: mediaAlt.trim() || undefined,
			})
			.run();
		closeMediaPopup();
	}
}

function handleMediaPopupKeydown(e: KeyboardEvent) {
	if (e.key === "Escape") {
		closeMediaPopup();
	} else if (e.key === "Enter" && e.ctrlKey && mediaUrl.trim()) {
		insertImage();
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

function formatFileSize(bytes: number): string {
	if (bytes < 1024) return bytes + " B";
	if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB";
	return (bytes / (1024 * 1024)).toFixed(1) + " MB";
}

function formatDate(dateString: string): string {
	return new Date(dateString).toLocaleDateString("en-US", {
		year: "numeric",
		month: "short",
		day: "numeric",
	});
}

const filteredMedia = $derived(
	media.filter((m) => {
		if (!searchQuery) return true;
		const query = searchQuery.toLowerCase();
		return (
			m.filename.toLowerCase().includes(query) ||
			m.original_filename.toLowerCase().includes(query)
		);
	}),
);
</script>

<div class="border rounded-lg overflow-hidden bg-white {className}">
	{#if editor}
		<div class="flex flex-wrap items-center gap-2 border-b bg-gray-50 p-2" data-ui={ui}>
			<select
				class="h-8 rounded border bg-white px-2 text-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500"
				onchange={(e) => applyBlock((e.target as HTMLSelectElement).value)}
				value={currentBlockValue()}
				title="Block format"
			>
				{#each BLOCK_OPTIONS as opt (opt.value)}
					<option value={opt.value}>{opt.label}</option>
				{/each}
			</select>

			<div class="h-6 w-px bg-gray-300" />

			<div class="inline-flex overflow-hidden rounded border bg-white">
				<button type="button" class="px-3 py-1 hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed transition-colors" class:bg-blue-100={editor?.isActive('bold')} class:text-blue-600={editor?.isActive('bold')} title="Bold (Ctrl+B)" onclick={() => run(() => editor?.chain().focus().toggleBold().run())} disabled={!editor?.can().chain().focus().toggleBold().run()}><strong>B</strong></button>
				<button type="button" class="px-3 py-1 hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed transition-colors border-l" class:bg-blue-100={editor?.isActive('italic')} class:text-blue-600={editor?.isActive('italic')} title="Italic (Ctrl+I)" onclick={() => run(() => editor?.chain().focus().toggleItalic().run())} disabled={!editor?.can().chain().focus().toggleItalic().run()}><em>I</em></button>
				<button type="button" class="px-3 py-1 hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed transition-colors border-l" class:bg-blue-100={editor?.isActive('underline')} class:text-blue-600={editor?.isActive('underline')} title="Underline (Ctrl+U)" onclick={() => run(() => editor?.chain().focus().toggleUnderline().run())} disabled={!editor?.can().chain().focus().toggleUnderline().run()}><u>U</u></button>
				<button type="button" class="px-3 py-1 hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed transition-colors border-l" class:bg-blue-100={editor?.isActive('strike')} class:text-blue-600={editor?.isActive('strike')} title="Strikethrough" onclick={() => run(() => editor?.chain().focus().toggleStrike().run())} disabled={!editor?.can().chain().focus().toggleStrike().run()}><s>S</s></button>
			</div>

			<div class="inline-flex overflow-hidden rounded border bg-white">
				<button type="button" class="px-3 py-1 hover:bg-gray-100 transition-colors" class:bg-blue-100={editor?.isActive('bulletList')} class:text-blue-600={editor?.isActive('bulletList')} title="Bullet list" onclick={() => run(() => editor?.chain().focus().toggleBulletList().run())}><span class="text-sm">• List</span></button>
				<button type="button" class="px-3 py-1 hover:bg-gray-100 transition-colors border-l" class:bg-blue-100={editor?.isActive('orderedList')} class:text-blue-600={editor?.isActive('orderedList')} title="Numbered list" onclick={() => run(() => editor?.chain().focus().toggleOrderedList().run())}><span class="text-sm">1. List</span></button>
			</div>

			<div class="inline-flex overflow-hidden rounded border bg-white">
				<button type="button" class="px-3 py-1 hover:bg-gray-100 transition-colors" class:bg-blue-100={editor?.isActive('blockquote')} class:text-blue-600={editor?.isActive('blockquote')} title="Blockquote" onclick={() => run(() => editor?.chain().focus().toggleBlockquote().run())}><span class="text-sm">&ldquo;Quote&rdquo;</span></button>
				<button type="button" class="px-3 py-1 hover:bg-gray-100 transition-colors border-l" class:bg-blue-100={editor?.isActive('codeBlock')} class:text-blue-600={editor?.isActive('codeBlock')} title="Code block" onclick={() => run(() => editor?.chain().focus().toggleCodeBlock().run())}><span class="text-sm">[code]</span></button>
			</div>

			<div class="inline-flex overflow-hidden rounded border bg-white">
				<button type="button" class="px-3 py-1 hover:bg-gray-100 transition-colors" title="Horizontal rule" onclick={() => run(() => editor?.chain().focus().setHorizontalRule().run())}><span class="text-sm">—</span></button>
				<button type="button" class="px-3 py-1 hover:bg-gray-100 transition-colors border-l" class:bg-blue-100={editor?.isActive('image')} class:text-blue-600={editor?.isActive('image')} title="Insert image" onclick={openMediaPopup}><span class="text-sm">📷 Image</span></button>
			</div>

			<div class="inline-flex overflow-hidden rounded border bg-white">
				<button type="button" class="px-3 py-1 hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed transition-colors" title="Undo (Ctrl+Z)" onclick={() => run(() => editor?.chain().focus().undo().run())} disabled={!editor?.can().chain().focus().undo().run()}><span class="text-sm">⤺ Undo</span></button>
				<button type="button" class="px-3 py-1 hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed transition-colors border-l" title="Redo (Ctrl+Shift+Z)" onclick={() => run(() => editor?.chain().focus().redo().run())} disabled={!editor?.can().chain().focus().redo().run()}><span class="text-sm">⤻ Redo</span></button>
			</div>
		</div>
	{/if}

	<div class="bg-white" bind:this={element}></div>
</div>

{#if showMediaPopup}
	<div class="fixed inset-0 z-50 overflow-y-auto" onkeydown={handleMediaPopupKeydown} role="dialog" aria-modal="true" tabindex="-1">
		<div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
			<div class="fixed inset-0 transition-opacity bg-gray-500 bg-opacity-75" onclick={closeMediaPopup}></div>

			<div class="inline-block w-full max-w-5xl my-8 overflow-hidden text-left align-middle transition-all transform bg-white rounded-lg shadow-xl">
				<div class="flex items-center justify-between px-6 py-4 border-b border-gray-200">
					<h3 class="text-lg font-semibold text-gray-900">Insert Image</h3>
					<button type="button" onclick={closeMediaPopup} class="text-2xl text-gray-400 hover:text-gray-600 transition-colors" title="Close (Esc)">×</button>
				</div>

				<div class="px-6 py-3 bg-gray-50 border-b border-gray-200">
					<div class="flex gap-2">
						<button type="button" onclick={() => insertMode = 'library'} class="px-4 py-2 rounded-lg font-medium transition-colors {insertMode === 'library' ? 'bg-blue-600 text-white' : 'bg-white text-gray-700 hover:bg-gray-100'}">Media Library</button>
						<button type="button" onclick={() => insertMode = 'url'} class="px-4 py-2 rounded-lg font-medium transition-colors {insertMode === 'url' ? 'bg-blue-600 text-white' : 'bg-white text-gray-700 hover:bg-gray-100'}">Insert by URL</button>
					</div>
				</div>

				{#if insertMode === 'library'}
					<div class="px-6 py-4 bg-gray-50 border-b border-gray-200">
						<div class="flex gap-4 items-center">
							<div class="relative flex-1">
								<svg class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
								<input type="text" bind:value={searchQuery} placeholder="Search media..." class="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
							</div>

							<div class="flex border border-gray-300 rounded-lg overflow-hidden">
								<button type="button" onclick={() => viewMode = 'grid'} class="px-3 py-2 transition-colors {viewMode === 'grid' ? 'bg-blue-600 text-white' : 'bg-white text-gray-700 hover:bg-gray-50'}"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" /></svg></button>
								<button type="button" onclick={() => viewMode = 'list'} class="px-3 py-2 transition-colors {viewMode === 'list' ? 'bg-blue-600 text-white' : 'bg-white text-gray-700 hover:bg-gray-50'}"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" /></svg></button>
							</div>

							<label class="px-4 py-2 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700 cursor-pointer flex items-center gap-2 transition-colors">
								{#if uploading}
									<svg class="animate-spin h-5 w-5" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
									Uploading...
								{:else}
									<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" /></svg>
									Upload
								{/if}
								<input type="file" accept="image/*" onchange={handleFileSelect} class="hidden" disabled={uploading} />
							</label>
						</div>
					</div>

					<div class="px-6 py-6 max-h-[500px] overflow-y-auto">
						{#if loadingMedia}
							<div class="flex items-center justify-center py-12">
								<svg class="animate-spin h-8 w-8 text-gray-400" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
							</div>
						{:else if filteredMedia.length === 0}
							<div class="text-center py-12 text-gray-500">
								<svg class="mx-auto h-12 w-12 text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>
								<p>No images found</p>
								<p class="text-sm mt-2">Upload an image to get started</p>
							</div>
						{:else if viewMode === 'grid'}
							<div class="grid grid-cols-4 gap-4">
								{#each filteredMedia as m (m.id)}
									<button type="button" onclick={() => selectMediaFromLibrary(m)} class="group relative aspect-square rounded-lg overflow-hidden border-2 transition-all hover:border-blue-500 {selectedMediaItem?.id === m.id ? 'border-blue-600 ring-2 ring-blue-600' : 'border-gray-200'}">
										<img src={getMediaUrl(m)} alt={m.original_filename} class="w-full h-full object-cover group-hover:scale-105 transition-transform" />
										{#if selectedMediaItem?.id === m.id}
											<div class="absolute inset-0 bg-blue-600 bg-opacity-20 flex items-center justify-center">
												<svg class="w-8 h-8 text-white" fill="currentColor" viewBox="0 0 24 24"><path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z" /></svg>
											</div>
										{/if}
									</button>
								{/each}
							</div>
						{:else}
							<div class="bg-white rounded-lg shadow overflow-hidden">
								<table class="min-w-full divide-y divide-gray-200">
									<thead class="bg-gray-50">
										<tr>
											<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Preview</th>
											<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Filename</th>
											<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Size</th>
											<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Uploaded</th>
											<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Action</th>
										</tr>
									</thead>
									<tbody class="bg-white divide-y divide-gray-200">
										{#each filteredMedia as m (m.id)}
											<tr class="hover:bg-gray-50 transition-colors {selectedMediaItem?.id === m.id ? 'bg-blue-50' : ''}">
												<td class="px-6 py-4"><img src={getMediaUrl(m)} alt={m.original_filename} class="h-12 w-12 object-cover rounded" /></td>
												<td class="px-6 py-4"><div class="text-sm font-medium text-gray-900 truncate max-w-xs">{m.original_filename}</div></td>
												<td class="px-6 py-4 text-sm text-gray-500">{formatFileSize(m.size_bytes)}</td>
												<td class="px-6 py-4 text-sm text-gray-500">{formatDate(m.created_at)}</td>
												<td class="px-6 py-4 text-sm font-medium"><button type="button" onclick={() => selectMediaFromLibrary(m)} class="text-blue-600 hover:text-blue-900">{selectedMediaItem?.id === m.id ? 'Selected ✓' : 'Select'}</button></td>
											</tr>
										{/each}
									</tbody>
								</table>
							</div>
						{/if}
					</div>

					{#if totalPages > 1}
						<div class="px-6 py-4 border-t border-gray-200 flex items-center justify-between">
							<button type="button" onclick={() => changePage(currentPage - 1)} disabled={currentPage === 1} class="px-4 py-2 border border-gray-300 rounded-lg bg-white text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors">Previous</button>
							<span class="text-sm text-gray-600">Page {currentPage} of {totalPages}</span>
							<button type="button" onclick={() => changePage(currentPage + 1)} disabled={currentPage === totalPages} class="px-4 py-2 border border-gray-300 rounded-lg bg-white text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors">Next</button>
						</div>
					{/if}
				{:else}
					<div class="px-6 py-6 space-y-4">
						<div>
							<label for="media-url" class="mb-1 block text-sm font-medium text-gray-700">Image URL <span class="text-red-500">*</span></label>
							<input id="media-url" type="url" bind:value={mediaUrl} placeholder="https://example.com/image.jpg" class="w-full rounded-md border border-gray-300 px-3 py-2 focus:border-transparent focus:ring-2 focus:ring-blue-500 focus:outline-none" autofocus />
						</div>

						<div>
							<label for="media-alt" class="mb-1 block text-sm font-medium text-gray-700">Alt Text (optional)</label>
							<input id="media-alt" type="text" bind:value={mediaAlt} placeholder="Description of the image" class="w-full rounded-md border border-gray-300 px-3 py-2 focus:border-transparent focus:ring-2 focus:ring-blue-500 focus:outline-none" />
							<p class="mt-1 text-xs text-gray-500">Helps with accessibility and SEO</p>
						</div>

						{#if mediaUrl.trim()}
							<div>
								<label class="mb-2 block text-sm font-medium text-gray-700">Preview:</label>
								<div class="rounded-md border border-gray-200 bg-gray-50 p-2">
									<img src={mediaUrl} alt={mediaAlt || 'Image preview'} class="mx-auto h-auto max-h-32 max-w-full rounded" onerror={(e) => { const img = e.target as HTMLImageElement; img.style.display = 'none'; if (!img.nextElementSibling) { const errorMsg = document.createElement('div'); errorMsg.textContent = '! Invalid image URL'; errorMsg.className = 'text-red-500 text-sm text-center py-4'; img.parentNode?.appendChild(errorMsg); } }} />
								</div>
							</div>
						{/if}
					</div>
				{/if}

				<div class="px-6 py-4 border-t border-gray-200 flex justify-end gap-3">
					<button type="button" onclick={closeMediaPopup} class="rounded-md border border-gray-300 px-4 py-2 text-gray-600 hover:bg-gray-50 focus:ring-2 focus:ring-gray-500 focus:outline-none transition-colors">Cancel</button>
					<button type="button" onclick={insertImage} disabled={!mediaUrl.trim()} class="rounded-md bg-blue-600 px-4 py-2 text-white hover:bg-blue-700 focus:ring-2 focus:ring-blue-500 focus:outline-none disabled:cursor-not-allowed disabled:bg-gray-300 transition-colors">Insert Image</button>
				</div>
			</div>
		</div>
	</div>
{/if}
