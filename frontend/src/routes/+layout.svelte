<script>
	import { Toaster } from 'svelte-sonner';

	let { data, children } = $props();
</script>

<Toaster closeButton richColors position="top-right" expand={true} />

<div class="font-inconsolata">
	{@render children?.()}
</div>
