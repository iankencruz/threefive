<script lang="ts">
	import Navbar from '$components/ui/Navbar.svelte';

	let { children } = $props();
</script>

<div class="min-h-screen bg-white">
	<Navbar variant="standard" />
	{@render children?.()}
</div>
