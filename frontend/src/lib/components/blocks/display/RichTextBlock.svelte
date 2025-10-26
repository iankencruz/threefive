<!-- frontend/src/lib/components/blocks/display/RichTextBlock.svelte -->
<script lang="ts">
	interface RichtextBlockData {
		content: string;
	}

	interface Props {
		data: RichtextBlockData;
	}

	let { data }: Props = $props();
</script>

<section class="py-12 bg-background">
	<div class="container mx-auto px-4 max-w-7xl">
		<article class="prose prose-lg max-w-none">
			{@html data.content}
		</article>
	</div>
</section>

