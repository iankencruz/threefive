<script lang="ts">
	import { createEventDispatcher } from "svelte";
	import { Image, X, Plus, GripVertical } from "lucide-svelte";
	import MediaPicker from "$lib/components/ui/MediaPicker.svelte";
	import type { Media } from "$api/media";

	interface Props {
		block?: {
			id?: string;
			type: "gallery";
			data: {
				title?: string;
				media_ids?: string[];
				media?: Media[];
			};
		};
	}

	let { block }: Props = $props();

	const dispatch = createEventDispatcher();

	let title = $state(block?.data?.title || "");
	let selectedMedia = $state<Media[]>(block?.data?.media || []);
	let showMediaPicker = $state(false);

	// Reactive statement to sync media_ids with selectedMedia
	$effect(() => {
		if (block?.data?.media) {
			selectedMedia = [...block.data.media];
		}
	});

	function handleMediaSelect(mediaId: string, media: Media) {
		// Check if media already exists
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

	function getBlockData() {
		return {
			id: block?.id,
			type: "gallery" as const,
			data: {
				title: title || undefined,
				media_ids: selectedMedia.map((m) => m.id),
			},
		};
	}

	function handleSave() {
		dispatch("save", getBlockData());
	}

	function handleDelete() {
		dispatch("delete");
	}

	// Expose getBlockData for parent component
	export { getBlockData };
</script>

<div class="border-2 border-slate-200 rounded-xl p-6 bg-white">
	<div class="flex items-center justify-between mb-4">
		<div class="flex items-center gap-3">
			<div class="w-10 h-10 bg-purple-500 rounded-lg flex items-center justify-center">
				<Image class="w-5 h-5 text-white" />
			</div>
			<div>
				<h3 class="font-semibold text-slate-900">Gallery Block</h3>
				<p class="text-sm text-slate-500">Add multiple images</p>
			</div>
		</div>
		<button
			type="button"
			onclick={handleDelete}
			class="text-red-600 hover:text-red-700 text-sm font-medium"
		>
			Remove Block
		</button>
	</div>

	<div class="space-y-4">
		<!-- Title Input -->
		<div>
			<label for="gallery-title" class="block text-sm font-medium text-slate-700 mb-2">
				Gallery Title (Optional)
			</label>
			<input
				id="gallery-title"
				type="text"
				bind:value={title}
				placeholder="e.g., Summer Vacation Photos"
				class="w-full px-4 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
			/>
		</div>

		<!-- Media Selection -->
		<div>
			<div class="flex items-center justify-between mb-3">
				<label class="block text-sm font-medium text-slate-700">
					Images ({selectedMedia.length})
				</label>
				<button
					type="button"
					onclick={() => showMediaPicker = true}
					class="flex items-center gap-2 px-4 py-2 bg-purple-500 text-white rounded-lg hover:bg-purple-600 transition-colors text-sm font-medium"
				>
					<Plus class="w-4 h-4" />
					Add Image
				</button>
			</div>

			{#if selectedMedia.length === 0}
				<div class="border-2 border-dashed border-slate-300 rounded-lg p-8 text-center">
					<div class="w-16 h-16 bg-slate-100 rounded-full flex items-center justify-center mx-auto mb-3">
						<Image class="w-8 h-8 text-slate-400" />
					</div>
					<p class="text-slate-600 mb-2">No images added yet</p>
					<p class="text-sm text-slate-500">Click "Add Image" to select images from your media library</p>
				</div>
			{:else}
				<div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
					{#each selectedMedia as media, index (media.id)}
						<div class="group relative bg-slate-50 rounded-lg overflow-hidden border-2 border-slate-200 hover:border-purple-400 transition-all">
							<div class="aspect-square bg-slate-200 relative">
								<img
									src={media.thumbnail_url || media.medium_url || media.url}
									alt={media.original_filename}
									class="w-full h-full object-cover"
								/>
								
								<!-- Overlay with controls -->
								<div class="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-40 transition-all flex items-center justify-center gap-2">
									<button
										type="button"
										onclick={() => moveMediaUp(index)}
										disabled={index === 0}
										class="w-8 h-8 bg-white text-slate-700 rounded-lg opacity-0 group-hover:opacity-100 transition-opacity hover:bg-slate-100 flex items-center justify-center disabled:opacity-50 disabled:cursor-not-allowed"
										title="Move up"
									>
										↑
									</button>
									<button
										type="button"
										onclick={() => moveMediaDown(index)}
										disabled={index === selectedMedia.length - 1}
										class="w-8 h-8 bg-white text-slate-700 rounded-lg opacity-0 group-hover:opacity-100 transition-opacity hover:bg-slate-100 flex items-center justify-center disabled:opacity-50 disabled:cursor-not-allowed"
										title="Move down"
									>
										↓
									</button>
									<button
										type="button"
										onclick={() => removeMedia(media.id)}
										class="w-8 h-8 bg-red-500 text-white rounded-lg opacity-0 group-hover:opacity-100 transition-opacity hover:bg-red-600 flex items-center justify-center"
										title="Remove"
									>
										<X class="w-4 h-4" />
									</button>
								</div>

								<!-- Position indicator -->
								<div class="absolute top-2 left-2 w-6 h-6 bg-purple-500 text-white rounded-full flex items-center justify-center text-xs font-bold">
									{index + 1}
								</div>
							</div>
							<div class="p-2">
								<p class="text-xs text-slate-600 truncate">{media.original_filename}</p>
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>

		<!-- Validation Messages -->
		{#if selectedMedia.length === 0}
			<p class="text-sm text-amber-600 flex items-center gap-2">
				<span class="text-lg">!</span>
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
