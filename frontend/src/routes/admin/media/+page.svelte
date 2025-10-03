<script lang="ts">
import { onMount } from "svelte";
import { mediaApi, type Media, type ErrorResponse } from "$api/media";
import MediaCard from "$components/media/MediaCard.svelte";
import MediaUploadModal from "$components/media/MediaUploadModal.svelte";

let media = $state<Media[]>([]);
let selectedIds = $state<Set<string>>(new Set());
let loading = $state(true);
let error = $state("");
let showUploadModal = $state(false);
let searchQuery = $state("");

onMount(() => {
	loadMedia();
});

async function loadMedia() {
	loading = true;
	error = "";

	try {
		const response = await mediaApi.listMedia(1, 50);
		media = response.data || [];
	} catch (err) {
		const errorResponse = err as ErrorResponse;
		error = errorResponse.message || "Failed to load media";
	} finally {
		loading = false;
	}
}

function toggleSelect(id: string) {
	const newSet = new Set(selectedIds);
	if (newSet.has(id)) {
		newSet.delete(id);
	} else {
		newSet.add(id);
	}
	selectedIds = newSet;
}

async function deleteSelected() {
	if (selectedIds.size === 0) return;
	if (
		!confirm(
			`Delete ${selectedIds.size} item${selectedIds.size > 1 ? "s" : ""}?`,
		)
	)
		return;

	try {
		await Promise.all(
			Array.from(selectedIds).map((id) => mediaApi.deleteMedia(id)),
		);
		selectedIds = new Set();
		await loadMedia();
	} catch (err) {
		const errorResponse = err as ErrorResponse;
		alert(errorResponse.message || "Failed to delete media");
	}
}

async function handleUploadComplete() {
	showUploadModal = false;
	await loadMedia();
}

const filteredMedia = $derived(
	media.filter((m) => {
		if (!searchQuery) return true;
		const query = searchQuery.toLowerCase();
		return (
			m.filename.toLowerCase().includes(query) ||
			m.original_filename.toLowerCase().includes(query)
		);
	}),
);
</script>

<div class="p-8 max-w-7xl mx-auto">
	<!-- Header -->
	<div class="flex justify-between items-start mb-8">
		<div>
			<h1 class="text-3xl font-bold text-gray-900">Media Library</h1>
			<p class="text-gray-600 mt-1">Manage your images, videos, and documents</p>
		</div>
		
		<button
			onclick={() => showUploadModal = true}
			class="flex items-center gap-2 px-4 py-2.5 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
		>
			<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
			</svg>
			Upload Media
		</button>
	</div>

	<!-- Search & Actions Bar -->
	<div class="flex gap-4 mb-6 items-center">
		<!-- Search -->
		<div class="relative flex-1 max-w-md">
			<svg class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
			</svg>
			<input
				type="text"
				placeholder="Search media..."
				bind:value={searchQuery}
				class="w-full pl-10 pr-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
			/>
		</div>

		<!-- Delete Button (shows when items selected) -->
		{#if selectedIds.size > 0}
			<button
				onclick={deleteSelected}
				class="flex items-center gap-2 px-4 py-2.5 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
			>
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
				</svg>
				Delete ({selectedIds.size})
			</button>
		{/if}
	</div>

	<!-- Content States -->
	{#if loading}
		<!-- Loading State -->
		<div class="text-center py-12">
			<div class="inline-block h-12 w-12 animate-spin rounded-full border-4 border-solid border-blue-600 border-r-transparent align-[-0.125em] motion-reduce:animate-[spin_1.5s_linear_infinite]" role="status">
				<span class="!absolute !-m-px !h-px !w-px !overflow-hidden !whitespace-nowrap !border-0 !p-0 ![clip:rect(0,0,0,0)]">Loading...</span>
			</div>
			<p class="mt-4 text-sm text-gray-600">Loading media...</p>
		</div>
	{:else if error}
		<!-- Error State -->
		<div class="text-center py-12">
			<svg class="mx-auto h-12 w-12 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
			</svg>
			<h3 class="mt-2 text-sm font-semibold text-gray-900">Error loading media</h3>
			<p class="mt-1 text-sm text-gray-500">{error}</p>
			<div class="mt-6">
				<button
					onclick={loadMedia}
					type="button"
					class="inline-flex items-center rounded-md bg-blue-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
				>
					<svg class="-ml-0.5 mr-1.5 h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
						<path fill-rule="evenodd" d="M15.312 11.424a5.5 5.5 0 01-9.201 2.466l-.312-.311h2.433a.75.75 0 000-1.5H3.989a.75.75 0 00-.75.75v4.242a.75.75 0 001.5 0v-2.43l.31.31a7 7 0 0011.712-3.138.75.75 0 00-1.449-.39zm1.23-3.723a.75.75 0 00.219-.53V2.929a.75.75 0 00-1.5 0V5.36l-.31-.31A7 7 0 003.239 8.188a.75.75 0 101.448.389A5.5 5.5 0 0113.89 6.11l.311.31h-2.432a.75.75 0 000 1.5h4.243a.75.75 0 00.53-.219z" clip-rule="evenodd" />
					</svg>
					Try again
				</button>
			</div>
		</div>
	{:else if media.length === 0}
		<!-- Empty State - No Media Uploaded -->
		<div class="text-center py-12 border border-gray-300 rounded-lg">
			<svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
			</svg>
			<h3 class="mt-2 text-sm font-semibold text-gray-900">No media files</h3>
			<p class="mt-1 text-sm text-gray-500">Get started by uploading your first image, video, or document.</p>
			<div class="mt-6">
				<button
					onclick={() => showUploadModal = true}
					type="button"
					class="inline-flex items-center rounded-md bg-blue-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
				>
					<svg class="-ml-0.5 mr-1.5 h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
						<path d="M10.75 4.75a.75.75 0 00-1.5 0v4.5h-4.5a.75.75 0 000 1.5h4.5v4.5a.75.75 0 001.5 0v-4.5h4.5a.75.75 0 000-1.5h-4.5v-4.5z" />
					</svg>
					Upload Media
				</button>
			</div>
		</div>
	{:else if filteredMedia.length === 0}
		<!-- Empty State - No Search Results -->
		<div class="text-center py-12">
			<svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
			</svg>
			<h3 class="mt-2 text-sm font-semibold text-gray-900">No results found</h3>
			<p class="mt-1 text-sm text-gray-500">No media files match "{searchQuery}". Try a different search term.</p>
			<div class="mt-6">
				<button
					onclick={() => searchQuery = ''}
					type="button"
					class="inline-flex items-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50"
				>
					<svg class="-ml-0.5 mr-1.5 h-5 w-5 text-gray-400" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
						<path fill-rule="evenodd" d="M15.312 11.424a5.5 5.5 0 01-9.201 2.466l-.312-.311h2.433a.75.75 0 000-1.5H3.989a.75.75 0 00-.75.75v4.242a.75.75 0 001.5 0v-2.43l.31.31a7 7 0 0011.712-3.138.75.75 0 00-1.449-.39zm1.23-3.723a.75.75 0 00.219-.53V2.929a.75.75 0 00-1.5 0V5.36l-.31-.31A7 7 0 003.239 8.188a.75.75 0 101.448.389A5.5 5.5 0 0113.89 6.11l.311.31h-2.432a.75.75 0 000 1.5h4.243a.75.75 0 00.53-.219z" clip-rule="evenodd" />
					</svg>
					Clear search
				</button>
			</div>
		</div>
	{:else}
		<!-- Media Grid -->
		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
			{#each filteredMedia as mediaItem (mediaItem.id)}
				<MediaCard
					item={mediaItem}
					selected={selectedIds.has(mediaItem.id)}
					onToggleSelect={() => toggleSelect(mediaItem.id)}
				/>
			{/each}
		</div>
	{/if}
</div>

<!-- Upload Modal -->
{#if showUploadModal}
	<MediaUploadModal
		onClose={() => showUploadModal = false}
		onComplete={handleUploadComplete}
	/>
{/if}
