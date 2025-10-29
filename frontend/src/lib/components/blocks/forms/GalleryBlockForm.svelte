<script lang="ts">
	import { Image, X, Plus } from "lucide-svelte";
	import MediaPicker from "$lib/components/ui/MediaPicker.svelte";
	import type { Media } from "$api/media";

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

	let title = $state(data?.title || "");
	let selectedMedia = $state<Media[]>(data?.media || []);
	let showMediaPicker = $state(false);

	// Sync initial media
	$effect(() => {
		if (data?.media && data.media.length !== selectedMedia.length) {
			selectedMedia = [...data.media];
		}
	});

	// Emit changes
	$effect(() => {
		onchange({
			title: title || undefined,
			media_ids: selectedMedia.map((m) => m.id),
		});
	});

	function handleMediaSelect(mediaId: string, media: Media) {
		const exists = selectedMedia.some((m) => m.id === mediaId);
		if (!exists) {
			selectedMedia = [...selectedMedia, media];
		}
		showMediaPicker = false;
	}

	function removeMedia(mediaId: string) {
		selectedMedia = selectedMedia.filter((m) => m.id !== mediaId);
	}

	function moveMediaUp(index: number) {
		if (index > 0) {
			const newMedia = [...selectedMedia];
			[newMedia[index - 1], newMedia[index]] = [
				newMedia[index],
				newMedia[index - 1],
			];
			selectedMedia = newMedia;
		}
	}

	function moveMediaDown(index: number) {
		if (index < selectedMedia.length - 1) {
			const newMedia = [...selectedMedia];
			[newMedia[index], newMedia[index + 1]] = [
				newMedia[index + 1],
				newMedia[index],
			];
			selectedMedia = newMedia;
		}
	}
</script>

<div class="space-y-4">
	<!-- Title Input -->
	<div>
		<label for="gallery-title" class="block text-sm font-medium mb-2">
			Gallery Title (Optional)
		</label>
		<input
			id="gallery-title"
			type="text"
			bind:value={title}
			placeholder="e.g., Summer Vacation Photos"
			class="w-full px-4 py-2 border border-gray-600 bg-surface rounded-lg focus:ring-2 focus:ring--accent focus:border-transparent"
		/>
	</div>

	<!-- Media Selection -->
	<div>
		<div class="flex items-center justify-between mb-3">
			<label class="block text-sm font-medium">
				Images ({selectedMedia.length})
			</label>
			<button
				type="button"
				onclick={() => showMediaPicker = true}
				class="flex items-center gap-2 px-3 py-1.5 bg-primary hover:accent rounded-lg transition-colors text-sm font-medium"
			>
				<Plus class="w-4 h-4" />
				Add Image
			</button>
		</div>

		{#if selectedMedia.length === 0}
			<div class="border-2 border-dashed border-gray-600 rounded-lg p-8 text-center">
				<div class="w-16 h-16 bg-gray-700 rounded-full flex items-center justify-center mx-auto mb-3">
					<Image class="w-8 h-8 text-gray-400" />
				</div>
				<p class="text-gray-400 mb-2">No images added yet</p>
				<p class="text-sm text-gray-500">Click "Add Image" to select images</p>
			</div>
		{:else}
			<div class="grid grid-cols-2 md:grid-cols-3 gap-3">
				{#each selectedMedia as media, index (media.id)}
					<div class="group relative bg-gray-800 rounded-lg overflow-hidden border border-gray-700 hover:border-accent transition-all">
						<div class="aspect-square bg-gray-700 relative">
							<img
								src={media.thumbnail_url || media.medium_url || media.url}
								alt={media.original_filename}
								class="w-full h-full object-cover"
							/>
							
							<!-- Controls overlay -->
							<div class="absolute inset-0  group-hover:bg-opacity-50 transition-all flex items-center justify-center gap-1">
								<button
									type="button"
									onclick={() => moveMediaUp(index)}
									disabled={index === 0}
									class="w-7 h-7 bg-white text-gray-900 rounded opacity-0 group-hover:opacity-100 transition-opacity hover:bg-gray-100 flex items-center justify-center disabled:opacity-30 disabled:cursor-not-allowed text-sm font-bold"
									title="Move up"
								>
									↑
								</button>
								<button
									type="button"
									onclick={() => moveMediaDown(index)}
									disabled={index === selectedMedia.length - 1}
									class="w-7 h-7 bg-white text-gray-900 rounded opacity-0 group-hover:opacity-100 transition-opacity hover:bg-gray-100 flex items-center justify-center disabled:opacity-30 disabled:cursor-not-allowed text-sm font-bold"
									title="Move down"
								>
									↓
								</button>
								<button
									type="button"
									onclick={() => removeMedia(media.id)}
									class="w-7 h-7 bg-red-600 text-white rounded opacity-0 group-hover:opacity-100 transition-opacity hover:bg-red-700 flex items-center justify-center"
									title="Remove"
								>
									<X class="w-4 h-4" />
								</button>
							</div>

							<!-- Position badge -->
							<div class="absolute top-2 left-2 w-6 h-6 bg-primary text-white rounded-full flex items-center justify-center text-xs font-bold shadow-lg">
								{index + 1}
							</div>
						</div>
						<div class="p-2">
							<p class="text-xs text-gray-400 truncate">{media.original_filename}</p>
						</div>
					</div>
				{/each}
			</div>
		{/if}

		<!-- Validation warning -->
		{#if selectedMedia.length === 0}
			<p class="text-sm text-amber-500 flex items-center gap-2 mt-2">
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
		onclose={() => showMediaPicker = false}
	/>
{/if}
