<!-- frontend/src/lib/components/blocks/forms/RichTextBlockForm.svelte -->
<script lang="ts">
	import RichTextEditor from '$components/ui/RichTextEditor.svelte';

	export interface RichtextBlockData {
		content: string;
	}

	interface Props {
		data: RichtextBlockData;
		onchange?: (data: RichtextBlockData) => void;
	}

	let { data = $bindable({ content: '' }), onchange }: Props = $props();

	// Watch for content changes and notify parent
	$effect(() => {
		if (onchange) {
			onchange(data);
		}
	});
</script>

<div class="space-y-2">
	<label for="content-editor" class="block text-sm font-medium">
		Content
		<span class="text-red-500">*</span>
	</label>
	<RichTextEditor id="content_editor" bind:value={data.content} />
	<p class="text-xs text-gray-500">
		Use the rich text editor to format your content with headings, lists, images, and more.
	</p>
</div>
