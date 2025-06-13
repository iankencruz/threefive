<!-- MediaList.svelte -->
<script lang="ts">
	import dayjs from 'dayjs';
	import MediaItem from './MediaItem.svelte';
	import { Ban, Pencil, Save, Trash2 } from '@lucide/svelte';

	let { media, refresh } = $props();
	console.log('Media List View');
</script>

<div class="flex flex-col gap-4">
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
				<div class="flex justify-between gap-4 rounded border bg-white p-4 shadow-sm">
					<div class="relative flex-shrink-0">
						<img src={item.url} alt={altText} class="h-24 w-32 rounded object-cover" />
					</div>

					<div class=" min-w-0 grow">
						{#if editMode}
							<div class="space-y-3">
								<div class="relative">
									<label
										for="name-list-{item.id}"
										class="mb-1 block text-sm font-medium text-gray-700">Title</label
									>
									<input
										type="text"
										id="name-list-{item.id}"
										class="block w-full rounded-md border-gray-300 px-3 py-2 text-sm focus:border-indigo-600 focus:ring-indigo-600"
										value={title}
										oninput={(e: Event) => updateTitle((e.target as HTMLInputElement).value)}
										placeholder="Enter title"
									/>
								</div>
								<div class="relative">
									<label
										for="altText-list-{item.id}"
										class="mb-1 block text-sm font-medium text-gray-700">Alt Text</label
									>
									<input
										type="text"
										id="altText-list-{item.id}"
										class="block w-full rounded-md border-gray-300 px-3 py-2 text-sm focus:border-indigo-600 focus:ring-indigo-600"
										value={altText}
										oninput={(e: Event) => updateAltText((e.target as HTMLInputElement).value)}
										placeholder="Enter alt text"
									/>
								</div>
							</div>
						{:else}
							<div>
								<h3 class="truncate text-lg font-medium text-gray-900">{title || 'Untitled'}</h3>
								<p class="mt-1 text-sm text-gray-600">{altText || 'No alt text'}</p>

								<div class="mt-2 grid grid-cols-2 gap-4">
									<span class=" text-xs text-gray-400">ID: {item.id}</span>
									<span class=" text-xs text-gray-400"
										>Last Updated: {dayjs(item.updated_at).format('DD-MM-YYYY : mm:ss')}</span
									>
								</div>
							</div>
						{/if}
					</div>

					<div class="flex flex-shrink-0 flex-col gap-2">
						<button
							type="button"
							onclick={toggleEdit}
							class="inline-flex items-center gap-2 rounded-md bg-white px-3 py-2 text-sm font-medium text-gray-700 ring-1 ring-gray-300 hover:bg-gray-50 hover:text-indigo-600"
						>
							{#if editMode}
								<Ban size={16} />
								Cancel
							{:else}
								<Pencil size={16} />
								Edit
							{/if}
						</button>

						{#if editMode}
							<button
								type="button"
								onclick={save}
								class="inline-flex items-center gap-2 rounded-md bg-indigo-600 px-3 py-2 text-sm font-medium text-white hover:bg-indigo-700 focus:ring-2 focus:ring-indigo-500 focus:outline-none"
							>
								<Save size={16} />
								Save
							</button>
						{/if}

						<button
							type="button"
							onclick={remove}
							class="inline-flex items-center gap-2 rounded-md bg-red-600 px-3 py-2 text-sm font-medium text-white hover:bg-red-700 focus:ring-2 focus:ring-red-500 focus:outline-none"
						>
							<Trash2 size={16} />
							Delete
						</button>
					</div>
				</div>
			{/snippet}
		</MediaItem>
	{/each}
</div>
