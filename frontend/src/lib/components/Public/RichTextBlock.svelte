<script lang="ts">
	let { block }: { block: { props: { html: string } } } = $props();
</script>

<div class="prose prose-neutral max-w-none">{@html block.props.html}</div>
