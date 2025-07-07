<script lang="ts">
	import { goto } from '$app/navigation';
	import { pages } from '$src/lib/store/pages.svelte';
	import { toast } from 'svelte-sonner';

	let { children } = $props();
	let sentinel: HTMLDivElement;
	let isSticky = $state(false);

	$effect(() => {
		async function loadPages() {
			try {
				const res = await fetch('/api/v1/admin/pages');
				const json = await res.json();
				pages.setPages(json.data);
			} catch (err) {
				toast.error('Failed to load pages');
			}
		}
		loadPages();
	});

	function createNew(): void {
		goto('/admin/pages');
	}
</script>

<!-- PAGE LAYOUT -->
<div class="grid min-h-screen w-full grid-cols-[16rem_1fr]">
	<!-- Sidebar -->
	<div class="border-r bg-white pr-4">
		<div class="sticky top-0 max-h-screen overflow-y-auto pt-6">
			<h2 class="mb-2 text-sm font-semibold tracking-wide text-gray-500 uppercase">Pages</h2>

			<button
				onclick={createNew}
				class="mb-3 w-full rounded bg-black px-3 py-2 text-sm font-medium text-white hover:bg-gray-800"
			>
				+ New Page
			</button>

			<ul class="space-y-1 pb-12">
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
	</div>

	<!-- Main content -->
	<div class="min-h-screen px-4 py-8">
		{@render children()}
	</div>
</div>
