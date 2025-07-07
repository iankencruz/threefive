<script lang="ts">
	import Breadcrumb from '$src/lib/components/Navigation/Breadcrumb.svelte';

	let { children } = $props();
</script>

<div class="">
	<Breadcrumb />
	{@render children()}
</div>
