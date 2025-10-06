<!-- frontend/src/lib/components/blocks/display/RichtextBlock.svelte -->
<script lang="ts">
interface RichtextBlockData {
	content: string;
}

interface Props {
	data: RichtextBlockData;
}

let { data }: Props = $props();

// Simple markdown-like parsing (you can replace with a proper markdown library later)
const formatContent = (text: string) => {
	return text
		.replace(/\*\*(.*?)\*\*/g, "<strong>$1</strong>") // **bold**
		.replace(/\*(.*?)\*/g, "<em>$1</em>") // *italic*
		.replace(/\n\n/g, "</p><p>") // paragraphs
		.replace(/\n/g, "<br>"); // line breaks
};
</script>

<section class="py-16 md:py-20">
	<div class="container mx-auto px-4 max-w-4xl">
		<div class="prose prose-lg max-w-none">
			<div class="text-gray-800 leading-relaxed">
				{@html '<p>' + formatContent(data.content) + '</p>'}
			</div>
		</div>
	</div>
</section>

<style>
	:global(.prose) {
		color: #374151;
		line-height: 1.75;
	}

	:global(.prose p) {
		margin-bottom: 1.25rem;
	}

	:global(.prose strong) {
		font-weight: 600;
		color: #1f2937;
	}

	:global(.prose em) {
		font-style: italic;
	}
</style>
