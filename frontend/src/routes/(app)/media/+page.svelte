<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchMedia, deleteMedia, updateMedia } from '$lib/api/media';
	import MediaGrid from '$src/components/MediaGrid.svelte';
	import MediaList from '$src/components/MediaList.svelte';
	import UploadModal from '$src/components/UploadModal.svelte';

	let view = $state('grid'); // 'list' or 'grid'
	let media = $state([]);
	let page = $state(1);
	let totalPages = $state(1);
	let modalOpen = $state(false);

	onMount(async () => {
		await loadMedia();
	});

	async function loadMedia() {
		const res = await fetchMedia(page);
		media = res.data?.items ?? [];
		totalPages = res.total_pages ?? 1;
		console.log('res.data.items: ', res.data.items);
		console.log('res total pages', totalPages);
		console.log($state.snapshot(media));
	}
</script>

<svelte:head>
	<title>Media | ThreeFive</title>
	<meta name="description" content="Media Admin Page" />
</svelte:head>

<section>
	<div class="md:flex md:items-center md:justify-between">
		<div class="min-w-0 flex-1">
			<h2 class="text-2xl/7 font-bold text-gray-900 sm:truncate sm:text-3xl sm:tracking-tight">
				Media Page
			</h2>
		</div>
	</div>

	<div>
		<div class="mt-6 flex justify-end gap-x-4 md:mt-0 md:ml-4">
			<!-- Open Upload Modal -->
			<button
				onclick={() => (modalOpen = true)}
				class="inline-flex items-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500"
			>
				Upload
			</button>

			<button
				onclick={() => (view = 'grid')}
				type="button"
				disabled={view === 'grid'}
				class="cursor:pointer inline-flex items-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-xs ring-1 ring-gray-300 ring-inset hover:bg-gray-50 disabled:bg-gray-200"
			>
				Grid
			</button>
			<button
				onclick={() => (view = 'list')}
				type="button"
				disabled={view === 'list'}
				class="cursor:pointer inline-flex items-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-xs ring-1 ring-gray-300 ring-inset hover:bg-gray-50 disabled:bg-gray-200"
			>
				List
			</button>
		</div>

		<div class="mt-4">
			{#if view === 'grid'}
				<MediaGrid {media} refresh={loadMedia} />
			{:else}
				<MediaList {media} refresh={loadMedia} />
			{/if}
		</div>
	</div>

	<UploadModal
		open={modalOpen}
		onclose={() => (modalOpen = false)}
		onuploaded={() => loadMedia()}
	/>
</section>
