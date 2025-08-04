<script lang="ts">
	import type { HeadingBlock } from '$lib/types';

	let { block }: { block: HeadingBlock } = $props();
	console.log('block heading:', block.props.title);
</script>

<h1 class="mb-4 text-4xl font-bold">{block.props.title}</h1>
<p>{block.props.description}</p>
