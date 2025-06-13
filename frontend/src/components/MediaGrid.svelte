<!-- MediaGrid.svelte -->
<script lang="ts">
	import MediaItem from './MediaItem.svelte';
	import { Ban, Pencil, Save, Trash2 } from '@lucide/svelte';

	let { media, refresh } = $props();
	console.log('Media Grid View');
</script>

<div class="grid grid-cols-2 gap-4 md:grid-cols-3 lg:grid-cols-4">
	{#each media as item (item.id)}
		<MediaItem {item} onrefresh={refresh}>
			{#snippet children({
				item,
				title,
				altText,
				editMode,
				save,
				remove,
				toggleEdit,
				updateTitle,
				updateAltText
			})}
				<div class="relative flex flex-col gap-2 rounded border">
					<div class="relative">
						<img
							src={item.url}
							alt={altText}
							class="aspect-video h-64 w-full rounded object-cover"
						/>
						{#if editMode}
							<div class="absolute bottom-0 grid w-full space-y-2 bg-white p-2 opacity-80">
								<div class="relative">
									<label
										for="name-{item.id}"
										class="absolute -top-2 left-2 inline-block rounded-lg bg-white px-1 text-xs font-medium text-gray-900"
										>Name</label
									>
									<input
										type="text"
										id="name-{item.id}"
										class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
										value={title}
										oninput={(e: Event) => updateTitle((e.target as HTMLInputElement).value)}
										placeholder="Title"
									/>
								</div>
								<div class="relative">
									<label
										for="altText-{item.id}"
										class="absolute -top-2 left-2 inline-block rounded-lg bg-white px-1 text-xs font-medium text-gray-900"
										>Alt Text</label
									>
									<input
										type="text"
										id="altText-{item.id}"
										class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
										value={altText}
										oninput={(e: Event) => updateAltText((e.target as HTMLInputElement).value)}
										placeholder="Alt Text"
									/>
								</div>
							</div>
						{:else}
							<div class="absolute bottom-0 w-full bg-white p-2 opacity-80">
								<p class="w-64 truncate">{title}</p>
								<p class="truncate text-sm text-gray-600">{altText}</p>
							</div>
						{/if}
					</div>
					<div class="absolute -top-4 right-0 flex justify-center">
						<span class="group isolate inline-flex -space-x-px rounded-md shadow-xs">
							<button
								type="button"
								onclick={toggleEdit}
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
			{/snippet}
		</MediaItem>
	{/each}
</div>
