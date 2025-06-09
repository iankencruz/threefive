<script lang="ts">
	import { updateMedia, deleteMedia } from '$lib/api/media';

	let { item, onrefresh } = $props();
	let title = $state(item.title ?? '');
	let altText = $state(item.alt_text ?? '');

	async function save() {
		await updateMedia(item.id, { title, alt_text: altText });
		onrefresh?.();
	}

	async function remove() {
		if (confirm('Delete this media item?')) {
			await deleteMedia(item.id);
			onrefresh?.();
		}
	}
</script>

<div class="flex flex-col gap-2 rounded border p-2">
	<img src={item.url} alt={altText} class="aspect-video w-full rounded object-cover" />

	<input type="text" bind:value={title} placeholder="Title" class="input input-sm" />
	<input type="text" bind:value={altText} placeholder="Alt text" class="input input-sm" />

	<div class="flex items-center justify-between">
		<button onclick={save} class="btn btn-sm">💾 Save</button>
		<button onclick={remove} class="btn btn-sm btn-error">🗑 Delete</button>
	</div>
</div>
