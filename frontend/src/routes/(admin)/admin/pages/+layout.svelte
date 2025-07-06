<!-- /admin/pages/+layout.svelte -->
<script lang="ts">
	import { onMount } from 'svelte';
	import type { Page } from '$lib/types';
	import { toast } from 'svelte-sonner';
	import { goto } from '$app/navigation';
	import { pages } from '$src/lib/store/pages.svelte';

	let { children } = $props();

	$effect(() => {
		async function loadPages() {
			try {
				const res = await fetch('/api/v1/admin/pages');
				const json = await res.json();
				pages.setPages(json.data); // ✅ use safe mutation
			} catch (err) {
				toast.error('Failed to load pages');
			}
		}
		loadPages();
	});
	//
	// async function loadPages() {
	// 	try {
	// 		const res = await fetch('/api/v1/admin/pages');
	// 		const json = await res.json();
	// 		pages = json.data;
	// 	} catch (err) {
	// 		toast.error('Failed to load pages');
	// 	}
	// }

	// refresh pages list if value changes (watch pages)

	function createNew(): void {
		goto('/admin/pages'); // Reset to blank form view
	}
</script>

<div class="flex h-full gap-4">
	<!-- Sidebar -->
	<div class="w-64 border-r pr-4">
		<h2 class="mb-2 text-sm font-semibold tracking-wide text-gray-500 uppercase">Pages</h2>

		<button
			onclick={createNew}
			class="mb-3 w-full rounded bg-black px-3 py-2 text-sm font-medium text-white hover:bg-gray-800"
		>
			+ New Page
		</button>

		<ul class="space-y-1">
			{#each pages.pagesStore as page}
				<li>
					<button
						class="w-full rounded px-3 py-2 text-left hover:bg-gray-100"
						onclick={() => goto(`/admin/pages/${page.slug}`)}
					>
						{page.title}
					</button>
				</li>
			{/each}
		</ul>
	</div>

	<!-- Content area -->
	<div class="h-full flex-1 overflow-auto">
		{@render children()}
	</div>
</div>
