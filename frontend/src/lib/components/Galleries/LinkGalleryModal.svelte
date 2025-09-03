<script lang="ts">
	import type { Gallery } from '$lib/types';

	let {
		open,
		pageSlug,
		onclose,
		onlinked,
		linkedGalleryIds = []
	} = $props<{
		open: boolean;
		pageSlug: string;
		onclose: () => void;
		onlinked: (gallery: Gallery) => void;
		linkedGalleryIds?: string[];
	}>();

	let galleries = $state<Gallery[]>([]);
	let loading = $state(false);

	async function loadGalleries() {
		try {
			loading = true;
			const res = await fetch(`/api/v1/admin/galleries`);
			if (!res.ok) throw new Error('Failed to fetch galleries');
			const json = await res.json();

			const all = json.data as Gallery[];
			galleries = all.filter((g) => !linkedGalleryIds.includes(g.id));
		} catch (err) {
			console.error('loadGalleries error:', err);
			galleries = [];
		} finally {
			loading = false;
		}
	}

	// ✅ Refresh when modal opens or linkedGalleryIds changes
	$effect(() => {
		if (open) {
			loadGalleries();
		}
	});

	function handleSelect(g: Gallery) {
		onlinked(g);
		onclose();
	}
</script>

{#if open}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
		<div class="w-full max-w-lg rounded bg-white p-4 shadow">
			<h2 class="mb-2 text-lg font-semibold">Link a Gallery</h2>

			{#if loading}
				<p>Loading galleries…</p>
			{:else if galleries.length > 0}
				<ul class="divide-y divide-gray-200">
					{#each galleries as g}
						<li class="flex items-center justify-between py-2">
							<div>
								<p class="font-medium">{g.title}</p>
								<p class="text-sm text-gray-500">{g.slug}</p>
							</div>
							<button
								class="rounded bg-indigo-600 px-2 py-1 text-sm text-white"
								onclick={() => handleSelect(g)}
							>
								Link
							</button>
						</li>
					{/each}
				</ul>
			{:else}
				<p class="mt-2 text-sm text-gray-500">All galleries already linked 🎉</p>
			{/if}

			<div class="mt-4 text-right">
				<button class="rounded border px-3 py-1 text-sm" onclick={onclose}> Close </button>
			</div>
		</div>
	</div>
{/if}
