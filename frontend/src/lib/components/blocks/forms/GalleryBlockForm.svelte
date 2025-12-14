<script lang="ts">
	import { Image, X, Plus } from 'lucide-svelte';
	import MediaPicker from '$lib/components/ui/MediaPicker.svelte';
	import type { Media } from '$api/media';

	export interface GalleryBlockData {
		title?: string;
		media_ids?: string[];
		media?: Media[];
	}

	interface Props {
		data: GalleryBlockData;
		onchange: (data: GalleryBlockData) => void;
	}

	let { data, onchange }: Props = $props();

	let title = $state(data?.title || '');
	let selectedMedia = $state<Media[]>(data?.media || []);
	let showMediaPicker = $state(false);

	// Sync initial media only once
	$effect(() => {
		if (data?.media && data.media.length > 0 && selectedMedia.length === 0) {
			selectedMedia = [...data.media];
		}
	});

	// Only notify on actual data changes
	function notifyChange() {
		onchange({
			title: title || undefined,
			media_ids: selectedMedia.map((m) => m.id),
			media: selectedMedia
		});
	}

	function handleMediaSelect(mediaId: string, media: Media) {
		const exists = selectedMedia.some((m) => m.id === mediaId);
		if (!exists) {
			selectedMedia = [...selectedMedia, media];
			notifyChange(); // Explicitly notify
		}
		showMediaPicker = false;
	}

	function removeMedia(mediaId: string) {
		selectedMedia = selectedMedia.filter((m) => m.id !== mediaId);
		notifyChange(); // Explicitly notify
	}

	function moveMediaUp(index: number) {
		if (index > 0) {
			const newMedia = [...selectedMedia];
			[newMedia[index - 1], newMedia[index]] = [newMedia[index], newMedia[index - 1]];
			selectedMedia = newMedia;
			notifyChange(); // Explicitly notify
		}
	}

	function moveMediaDown(index: number) {
		if (index < selectedMedia.length - 1) {
			const newMedia = [...selectedMedia];
			[newMedia[index], newMedia[index + 1]] = [newMedia[index + 1], newMedia[index]];
			selectedMedia = newMedia;
			notifyChange(); // Explicitly notify
		}
	}
</script>

<div class="space-y-4">
	<!-- Title Input -->
	<div>
		<label for="gallery-title" class="mb-2 block text-sm font-medium">
			Gallery Title (Optional)
		</label>
		<input
			id="gallery-title"
			type="text"
			bind:value={title}
			onblur={notifyChange}
			placeholder="e.g., Summer Vacation Photos"
			class="focus:ring--accent w-full rounded-lg border border-gray-600 bg-surface px-4 py-2 focus:border-transparent focus:ring-2"
		/>
	</div>

	<!-- Media Selection -->
	<div>
		<div class="mb-3 flex items-center justify-between">
			<label for="add-image-button" class="block text-sm font-medium">
				Images ({selectedMedia.length})
			</label>
			<button
				name="add-image-button"
				type="button"
				onclick={() => (showMediaPicker = true)}
				class="hover:accent flex items-center gap-2 rounded-lg bg-primary px-3 py-1.5 text-sm font-medium transition-colors"
			>
				<Plus class="h-4 w-4" />
				Add Image
			</button>
		</div>

		{#if selectedMedia.length === 0}
			<div class="rounded-lg border-2 border-dashed border-gray-600 p-8 text-center">
				<div
					class="mx-auto mb-3 flex h-16 w-16 items-center justify-center rounded-full bg-gray-700"
				>
					<Image class="h-8 w-8 text-gray-400" />
				</div>
				<p class="mb-2 text-gray-400">No images added yet</p>
				<p class="text-sm text-gray-500">Click "Add Image" to select images</p>
			</div>
		{:else}
			<div class="grid grid-cols-2 gap-3 md:grid-cols-5">
				{#each selectedMedia as media, index (media.id)}
					<div
						class="group relative overflow-hidden rounded-lg border border-gray-700 bg-gray-800 transition-all hover:border-accent"
					>
						<div class="relative aspect-square bg-gray-700">
							<img
								src={media.thumbnail_url || media.medium_url || media.url}
								alt={media.original_filename}
								class="h-full w-full object-cover"
							/>

							<!-- Controls overlay -->
							<div
								class="group-hover:bg-opacity-50 absolute top-2 right-2 flex flex-col items-center justify-center gap-1 transition-all"
							>
								<button
									type="button"
									onclick={() => removeMedia(media.id)}
									class="flex h-7 w-7 items-center justify-center rounded bg-red-600 text-white opacity-0 transition-opacity group-hover:opacity-100 hover:bg-red-700"
									title="Remove"
								>
									<X class="h-4 w-4" />
								</button>
								<button
									type="button"
									onclick={() => moveMediaUp(index)}
									disabled={index === 0}
									class="flex h-7 w-7 items-center justify-center rounded bg-white text-sm font-bold text-gray-900 opacity-0 transition-opacity group-hover:opacity-100 hover:bg-gray-100 disabled:hidden disabled:cursor-not-allowed"
									title="Move up"
								>
									↑
								</button>
								<button
									type="button"
									onclick={() => moveMediaDown(index)}
									disabled={index === selectedMedia.length - 1}
									class="flex h-7 w-7 items-center justify-center rounded bg-white text-sm font-bold text-gray-900 opacity-0 transition-opacity group-hover:opacity-100 hover:bg-gray-100 disabled:hidden disabled:cursor-not-allowed"
									title="Move down"
								>
									↓
								</button>
							</div>

							<!-- Position badge -->
							<div
								class="absolute top-2 left-2 flex h-6 w-6 items-center justify-center rounded-full bg-primary text-xs font-bold text-white shadow-lg"
							>
								{index + 1}
							</div>
						</div>
						<div class="p-2">
							<p class="truncate text-xs text-gray-400">{media.original_filename}</p>
						</div>
					</div>
				{/each}
			</div>
		{/if}

		<!-- Validation warning -->
		{#if selectedMedia.length === 0}
			<p class="mt-2 flex items-center gap-2 text-sm text-amber-500">
				<span>!</span>
				Gallery must have at least one image
			</p>
		{/if}
	</div>
</div>

<!-- Media Picker Modal -->
{#if showMediaPicker}
	<MediaPicker
		show={showMediaPicker}
		onselect={handleMediaSelect}
		onclose={() => (showMediaPicker = false)}
	/>
{/if}
