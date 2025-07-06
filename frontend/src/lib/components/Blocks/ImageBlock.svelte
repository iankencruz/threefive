<script lang="ts">
	import { onMount, tick } from 'svelte';
	import type { MediaItem, Block } from '$lib/types';
	import type { ImageBlock } from '$lib/types';
	import { getMediaById } from '$lib/api/media'; // 🧠 Add this API if not defined
	import LinkMediaModal from '$lib/components/Media/LinkMediaModal.svelte';

	let { block, onupdate } = $props<{
		block: Block;
		onupdate: (block: Block) => void;
	}>();

	let modalOpen = $state(false);
	let linkedMedia: MediaItem | null = $state(null);

	let config: ImageBlock = {
		media_id: block.props?.media_id || '',
		alt: block.props?.alt || '',
		size: block.props?.size || 'medium',
		alignment: block.props?.alignment || 'center',
		object_fit: block.props?.object_fit || 'cover',
		object_position: block.props?.object_position || 'center'
	};

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

<div class="space-y-4 border-t border-gray-200 pt-4">
	<!-- Media Section -->
	<div>
		<label class="mb-1 block text-xs font-light text-gray-700">Media</label>

		{#if linkedMedia}
			<div class="relative h-64 w-64">
				<img
					src={linkedMedia.thumbnail_url || linkedMedia.url}
					alt={linkedMedia.alt_text || linkedMedia.title || 'Media'}
					class=" h-64 w-64 rounded border object-cover"
				/>
				<button
					onclick={clearMedia}
					class="absolute top-1 right-1 rounded-full bg-white p-1 px-2 text-xs text-gray-600 shadow hover:bg-gray-100"
				>
					Clear
				</button>
			</div>
		{:else}
			<div class="flex items-center gap-2">
				<input
					type="text"
					name="media_id"
					class="w-full rounded border px-3 py-2"
					value={config.media_id}
					oninput={handleChange}
				/>
				<button
					type="button"
					class="rounded bg-gray-200 px-2 py-1 text-xs text-gray-600 hover:bg-gray-300"
					onclick={() => (modalOpen = true)}
				>
					Select
				</button>
			</div>
		{/if}
	</div>

	<!-- Alt text -->
	<div>
		<label class="block text-xs font-light text-gray-700">Alt Text</label>
		<input
			type="text"
			name="alt"
			class="w-full rounded border px-3 py-2"
			placeholder="Alt text"
			value={config.alt}
			oninput={handleChange}
		/>
	</div>

	<!-- Size -->
	<div>
		<label class="block text-xs font-light text-gray-700">Image Size</label>
		<select name="size" class="w-full rounded border px-3 py-2" onchange={handleChange}>
			<option value="small" selected={config.size === 'small'}>Small</option>
			<option value="medium" selected={config.size === 'medium'}>Medium</option>
			<option value="large" selected={config.size === 'large'}>Large</option>
		</select>
	</div>

	<!-- Alignment -->
	<div>
		<label class="block text-xs font-light text-gray-700">Alignment</label>
		<select name="alignment" class="w-full rounded border px-3 py-2" onchange={handleChange}>
			<option value="left" selected={config.alignment === 'left'}>Left</option>
			<option value="center" selected={config.alignment === 'center'}>Center</option>
			<option value="right" selected={config.alignment === 'right'}>Right</option>
		</select>
	</div>

	<!-- Object Fit -->
	<div>
		<label class="block text-xs font-light text-gray-700">Object Fit</label>
		<select name="object_fit" class="w-full rounded border px-3 py-2" onchange={handleChange}>
			<option value="cover" selected={config.object_fit === 'cover'}>Cover</option>
			<option value="contain" selected={config.object_fit === 'contain'}>Contain</option>
		</select>
	</div>

	<!-- Object Position -->
	<div>
		<label class="block text-xs font-light text-gray-700">Object Position</label>
		<select name="object_position" class="w-full rounded border px-3 py-2" onchange={handleChange}>
			<option value="top" selected={config.object_position === 'top'}>Top</option>
			<option value="center" selected={config.object_position === 'center'}>Center</option>
			<option value="bottom" selected={config.object_position === 'bottom'}>Bottom</option>
		</select>
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
