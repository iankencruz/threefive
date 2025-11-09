<script lang="ts">
	import "$src/app.css";
	import favicon from "$lib/assets/favicon.svg";
	import { Toaster } from "svelte-sonner";
	import Navbar from "$components/ui/Navbar.svelte";

	const { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<Navbar/>
{@render children?.()}



<Toaster richColors  expand={false} position="top-right" closeButton />
