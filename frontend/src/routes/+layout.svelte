<script lang="ts">
	import './layout.css';
	import type { Snippet } from 'svelte';
	import Navigation from '$lib/components/Navigation.svelte';
	import { Plus } from '@lucide/svelte';

	let { children }: { children: Snippet } = $props();
</script>

<div class="flex min-h-screen bg-background font-sans text-foreground">
	<Navigation />

	<div class="flex w-full flex-1 flex-col lg:pl-64">
		<header
			class="sticky top-0 z-10 flex h-16 items-center justify-between border-b border-border bg-background/80 px-4 backdrop-blur-md md:px-6"
		>
			<div class="ml-12 flex items-center gap-3 lg:ml-0">
				<h1 class="text-lg font-semibold tracking-tight">Admin</h1>
			</div>
		</header>

		<main class="w-full max-w-7xl p-6 lg:p-10">
			{@render children()}
		</main>
	</div>
</div>
