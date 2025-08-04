<script lang="ts">
	import { goto } from '$app/navigation';
	import { pages } from '$lib/store/pages.svelte';
	import { toast } from 'svelte-sonner';

	let { children } = $props();

	function createNew(): void {
		goto('/admin/pages');
	}
</script>

<!-- PAGE LAYOUT -->
<div class="grid min-h-screen w-full grid-cols-1">
	<!-- Main content -->
	<div class="mx-auto min-h-screen w-full px-4 py-8">
		{@render children()}
	</div>
</div>
