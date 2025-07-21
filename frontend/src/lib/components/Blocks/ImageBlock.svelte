<script lang="ts">
	import { onMount } from 'svelte';
	import type { MediaItem, Block } from '$lib/types';
	import type { ImageBlock } from '$lib/types';
	import { getMediaById } from '$lib/api/media'; // 🧠 Add this API if not defined
	import LinkMediaModal from '$lib/components/Media/LinkMediaModal.svelte';
	import { Image } from '@lucide/svelte';

	let { block, onupdate } = $props<{
		block: Block;
		onupdate: (block: Block) => void;
	}>();

	let modalOpen = $state(false);
	let linkedMedia: MediaItem | null = $state(null);

	let config = $state<ImageBlock['props']>({
		media_id: block.props?.media_id ?? '',
		alt_text: block.props?.alt_text ?? '',
		align: block.props?.align ?? 'center',
		object_fit: block.props?.object_fit ?? 'cover'
	});

	onMount(async () => {
		if (config.media_id && !linkedMedia) {
			try {
				linkedMedia = await getMediaById(config.media_id);
			} catch (err) {
				console.error('❌ Failed to fetch media by ID:', err);
			}
		}
	});

	function handleChange(e: Event) {
		const target = e.target as HTMLInputElement | HTMLSelectElement;
		const name = target.name as keyof ImageBlock;
		config = { ...config, [name]: target.value };

		onupdate({
			...block,
			props: { ...block.props, ...config }
		});
	}

	function handleMediaLinked(media: MediaItem) {
		config = { ...config, media_id: media.id };
		linkedMedia = media;

		onupdate({
			...block,
			props: { ...block.props, ...config }
		});

		modalOpen = false;
	}

	function clearMedia() {
		config = { ...config, media_id: '' };
		linkedMedia = null;
		onupdate({ ...block, props: { ...block.props, ...config } });
	}
</script>

<div class="flex h-full w-full flex-col gap-8 pt-4 md:flex-row">
	<!-- Media Section -->
	<div class="h-full w-4/6">
		{#if linkedMedia}
			<div class="relative h-64 w-full">
				<img
					src={linkedMedia.thumbnail_url || linkedMedia.url}
					alt={linkedMedia.alt_text || linkedMedia.title || 'Media'}
					class=" h-full w-full rounded border object-cover"
				/>
				<button
					onclick={clearMedia}
					class="absolute top-2 right-4 rounded-full bg-white p-1 px-2 text-xs text-gray-600 shadow hover:bg-gray-100"
				>
					Clear
				</button>
			</div>
		{:else}
			<div
				class=" flex h-full justify-center rounded-lg border border-dashed border-gray-700/25 px-6 py-10"
			>
				<div class="text-center">
					<Image size={24} />
					<div class="mt-4 flex text-sm/6 text-gray-400">
						<label
							for="file-upload"
							class="relative cursor-pointer rounded-md bg-gray-900 font-semibold text-white focus-within:ring-2 focus-within:ring-indigo-600 focus-within:ring-offset-2 focus-within:ring-offset-gray-900 focus-within:outline-hidden hover:text-indigo-500"
						>
							<button
								type="button"
								class="rounded bg-gray-200 px-2 py-1 text-xs text-gray-600
								hover:bg-gray-300"
								onclick={() => (modalOpen = true)}
							>
								Select
							</button>
							<input
								id="file-upload"
								type="text"
								class="sr-only"
								name="media_id"
								value={config.media_id}
								oninput={handleChange}
							/>
						</label>
						<p class="pl-1">or drag and drop</p>
					</div>
					<p class="text-xs/5 text-gray-400">PNG, JPG, GIF up to 10MB</p>
				</div>
			</div>
		{/if}
	</div>

	<div class=" w-2/6 space-y-2">
		<!-- Alt text -->
		<div>
			<label for="alt_text" class="block text-xs font-light text-gray-700">Alt Text</label>
			<input
				type="text"
				name="alt_text"
				class="w-full rounded border px-3 py-1"
				placeholder="Alt text"
				value={config.alt_text}
				oninput={handleChange}
			/>
		</div>

		<!-- Alignment -->
		<div>
			<label for="align" class="block text-xs font-light text-gray-700">Alignment</label>
			<select name="align" class="w-full rounded border px-3 py-1" onchange={handleChange}>
				<option value="left" selected={config.align === 'left'}>Left</option>
				<option value="center" selected={config.align === 'center'}>Center</option>
				<option value="right" selected={config.align === 'right'}>Right</option>
			</select>
		</div>

		<!-- Object Fit -->
		<div>
			<label for="object_fit" class="block text-xs font-light text-gray-700">Object Fit</label>
			<select
				id="object_fit"
				name="object_fit"
				class="w-full rounded border px-3 py-1"
				onchange={handleChange}
			>
				<option value="cover" selected={config.object_fit === 'cover'}>Cover</option>
				<option value="contain" selected={config.object_fit === 'contain'}>Contain</option>
				<option value="fill" selected={config.object_fit === 'fill'}>Fill</option>
				<option value="scale-down" selected={config.object_fit === 'scale-down'}>Scale Down</option>
				<option value="none" selected={config.object_fit === 'none'}>None</option>
			</select>
		</div>
	</div>
</div>

<LinkMediaModal
	open={modalOpen}
	context={{ type: 'block', id: block.id }}
	onclose={() => (modalOpen = false)}
	onlinked={handleMediaLinked}
	selectOnly={true}
	linkedMediaIds={[config.media_id]}
/>
