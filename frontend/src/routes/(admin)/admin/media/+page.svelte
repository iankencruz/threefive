<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchMedia, deleteMedia, updateMedia } from '$lib/api/media';
	import UploadModal from '$lib/components/Media/UploadModal.svelte';
	import MediaGrid from '$lib/components/Media/MediaGrid.svelte';
	import MediaList from '$lib/components/Media/MediaList.svelte';
	import { LayoutGrid, LayoutList } from '@lucide/svelte';
	import Pagination from '$lib/components/Navigation/Pagination.svelte';
	import { toast } from 'svelte-sonner';

	let view = $state('grid'); // 'list' or 'grid'
	let modalOpen = $state(false);

	let page = $state(1);
	let limit = $state(10);

	let media = $state([]);
	let totalPages = $state(1);
	let totalCount = $state(1);

	// onMount(async () => {
	// 	await loadMedia();
	// });

	async function loadMedia(p: number) {
		try {
			const result = await fetchMedia(p, limit);
			media = result.items ?? [];
			totalPages = result.total_pages ?? 1;
			totalCount = result.total_count ?? 0;
		} catch (err) {
			console.error(err);
			toast.error('Failed to load media');
		}
	}

	// media = res.items ?? [];
	// totalPages = res.total_pages ?? 1;

	$effect(() => {
		loadMedia(page);
	});

	function handlePageChange(newPage: number) {
		console.log('Changing to page', newPage);
		page = newPage;
	}

	function refreshCurrentPage() {
		loadMedia(page); // ✅ preserve current page
	}
</script>

<svelte:head>
	<title>Media | ThreeFive</title>
	<meta name="description" content="Media Admin Page" />
</svelte:head>

<section class="relative">
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

			<span class="isolate inline-flex rounded-md shadow-xs">
				<button
					onclick={() => (view = 'grid')}
					type="button"
					disabled={view === 'grid'}
					class="relative inline-flex cursor-pointer items-center rounded-l-md bg-white px-2 py-2 text-sm font-semibold text-gray-900 ring-1 ring-gray-300 ring-inset hover:bg-gray-50 focus:z-10 disabled:bg-gray-200"
				>
					<LayoutGrid size={18} />
				</button>
				<button
					onclick={() => (view = 'list')}
					type="button"
					disabled={view === 'list'}
					class="relative -ml-px inline-flex cursor-pointer items-center rounded-r-md bg-white px-2 py-2 text-sm font-semibold text-gray-900 ring-1 ring-gray-300 ring-inset hover:bg-gray-50 focus:z-10 disabled:bg-gray-200"
				>
					<LayoutList size={18} />
				</button>
			</span>
		</div>
		{#if media.length > 0}
			<div class="mt-12">
				{#if view === 'grid'}
					<MediaGrid {media} refresh={refreshCurrentPage} />
				{:else}
					<MediaList {media} refresh={refreshCurrentPage} />
				{/if}
			</div>
		{:else}
			<div class="my-8">
				<button
					onclick={() => (modalOpen = true)}
					type="button"
					class="relative block w-full rounded-lg border-2 border-dashed border-gray-300 p-12 text-center hover:border-gray-400 focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 focus:outline-hidden"
				>
					<svg
						class="mx-auto size-12 text-gray-400"
						stroke="currentColor"
						fill="none"
						viewBox="0 0 48 48"
						aria-hidden="true"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M8 14v20c0 4.418 7.163 8 16 8 1.381 0 2.721-.087 4-.252M8 14c0 4.418 7.163 8 16 8s16-3.582 16-8M8 14c0-4.418 7.163-8 16-8s16 3.582 16 8m0 0v14m0-4c0 4.418-7.163 8-16 8S8 28.418 8 24m32 10v6m0 0v6m0-6h6m-6 0h-6"
						/>
					</svg>
					<span class="mt-2 block text-sm font-semibold text-gray-900">Upload Media</span>
				</button>
			</div>
		{/if}
	</div>

	<div class="absolute w-full">
		<Pagination
			totalMedia={totalCount}
			pageSize={limit}
			{page}
			{totalPages}
			onchange={handlePageChange}
		/>
	</div>

	<UploadModal
		open={modalOpen}
		onclose={() => (modalOpen = false)}
		onuploaded={() => loadMedia(page)}
		maxWidth="max-w-2xl"
	/>
</section>
