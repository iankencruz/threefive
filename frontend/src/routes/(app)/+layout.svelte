<script lang="ts">
	import "$src/app.css";
	import favicon from "$lib/assets/favicon.svg";
	import { Toaster } from "svelte-sonner";

	const { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="bg-background min-h-screen">{@render children?.()}</div>



<Toaster richColors  expand={false} position="top-right" closeButton />
