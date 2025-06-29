<script lang="ts">
	import Breadcrumb from '$src/lib/components/Navigation/Breadcrumb.svelte';

	let { children } = $props();
</script>

<div class="px-6">
	<Breadcrumb />
	{@render children()}
</div>
