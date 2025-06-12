<script lang="ts">
	import { updateMedia, deleteMedia } from '$lib/api/media';
	import { Ban, Pencil, Save, Trash2 } from '@lucide/svelte';

	let { item, onrefresh } = $props();
	let title = $state(item.title ?? '');
	let altText = $state(item.alt_text ?? '');

	let editMode = $state(false);

	async function save() {
		await updateMedia(item.id, { title, alt_text: altText });
		console.log(item);
		onrefresh?.();
	}

	async function remove() {
		if (confirm('Delete this media item?')) {
			await deleteMedia(item.id);
			onrefresh?.();
		}
	}
</script>

<div class="relative flex flex-col gap-2 rounded border">
	<div class="relative">
		<img src={item.url} alt={altText} class=" aspect-video h-64 w-full rounded object-cover" />
		{#if editMode}
			<div class="absolute bottom-0 grid w-full space-y-2 bg-white p-2 opacity-80">
				<div class="relative">
					<label
						for="name"
						class="absolute -top-2 left-2 inline-block rounded-lg bg-white px-1 text-xs font-medium text-gray-900"
						>Name</label
					>
					<input
						type="text"
						name="name"
						id="name"
						class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
						bind:value={title}
						placeholder="Title"
					/>
				</div>
				<div class="relative">
					<label
						for="altText"
						class="absolute -top-2 left-2 inline-block rounded-lg bg-white px-1 text-xs font-medium text-gray-900"
						>Alt Text</label
					>
					<input
						type="text"
						name="altText"
						id="name"
						class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
						bind:value={altText}
						placeholder="Alt Text"
					/>
				</div>
			</div>
		{:else}
			<div class="absolute bottom-0 w-full bg-white p-2 opacity-80">
				<p class="w-64 truncate">{title}</p>
				<p>{altText}</p>
			</div>
		{/if}
	</div>
	<div class="absolute -top-4 right-0 flex justify-center">
		<span class="group isolate inline-flex -space-x-px rounded-md shadow-xs">
			<button
				type="button"
				onclick={() => (editMode = !editMode)}
				class="relative inline-flex items-center rounded-l-md bg-white px-3 py-2 text-gray-400 ring-1 ring-gray-300 ring-inset hover:cursor-pointer hover:bg-gray-50 hover:text-indigo-600 focus:z-10"
			>
				<span class="sr-only">Edit</span>
				{#if editMode}
					<Ban size={16} />
				{:else}
					<Pencil size={16} />
				{/if}
			</button>
			{#if editMode}
				<button
					type="button"
					onclick={save}
					class="relative inline-flex items-center bg-white px-3 py-2 text-gray-400 ring-1 ring-gray-300 ring-inset hover:cursor-pointer hover:bg-gray-50 hover:text-indigo-600 focus:z-10"
				>
					<span class="sr-only">Save</span>
					<Save size={16} />
				</button>
			{/if}
			<button
				type="button"
				onclick={remove}
				class="relative inline-flex items-center rounded-r-md bg-white px-3 py-2 text-gray-400 ring-1 ring-gray-300 ring-inset hover:cursor-pointer hover:bg-gray-50 hover:text-red-600 focus:z-10"
			>
				<span class="sr-only">Delete</span>
				<Trash2 size={16} />
			</button>
		</span>
	</div>
</div>
