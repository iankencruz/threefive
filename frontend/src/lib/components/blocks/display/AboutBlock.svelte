<script lang="ts">
	import type { AboutBlockData } from '../forms/AboutBlockForm.svelte';

	interface Props {
		data: AboutBlockData;
	}

	let { data }: Props = $props();
</script>

<section class="bg-background py-12">
	<h1>{data?.heading}</h1>
	<h2>{data?.subheading}</h2>
</section>
