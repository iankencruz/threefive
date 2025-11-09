<!-- routes/(public)/+layout.svelte -->
<script lang="ts">
	import "$src/app.css";
	import { setContext } from "svelte";
	import Navbar from "$lib/components/ui/Navbar.svelte";
	import { Toaster } from "svelte-sonner";

	const { children } = $props();

	// Create reactive state for navbar variant
	let navbarVariant = $state<"transparent" | "opaque">("opaque");

	// Share navbar control via context
	setContext("navbar", {
		get variant() {
			return navbarVariant;
		},
		setVariant: (v: "transparent" | "opaque") => {
			navbarVariant = v;
		},
	});
</script>

<!-- Navbar with reactive variant -->
<Navbar variant={navbarVariant} />

<!-- Page content -->
{@render children()}


