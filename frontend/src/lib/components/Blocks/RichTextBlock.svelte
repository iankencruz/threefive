<script lang="ts">
	import { Tipex, type TipexEditor } from '@friendofsvelte/tipex';
	import '@friendofsvelte/tipex/styles/index.css';

	import type { RichTextBlock } from '$lib/types';

	let {
		block,
		onupdate
	}: {
		block: RichTextBlock;
		onupdate: (updated: RichTextBlock) => void;
	} = $props();

	let body = $state(block.props?.html ?? '');
	let editor = $state<TipexEditor>();

	$effect(() => {
		if (!editor) return;

		editor.on('update', () => {
			const html = editor?.getHTML() ?? '';
			body = html;
			onupdate({ ...block, props: { html } });
		});
	});
</script>

<Tipex
	{body}
	bind:tipex={editor}
	floating
	focal
	class="mt-2 mb-0 h-[35vh] min-h-[15vh] border border-neutral-200"
/>
