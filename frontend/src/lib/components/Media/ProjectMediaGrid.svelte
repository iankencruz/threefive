<script lang="ts">
	import type { MediaItem } from '$src/lib/types';
	interface Props {
		media: MediaItem[];
	}

	let { media }: Props = $props();
</script>

<section class="mt-10 pb-16" aria-labelledby="media-heading">
	<h2 id="media-heading" class="mb-4 text-xl font-semibold text-gray-900">Media</h2>

	<ul
		role="list"
		class="grid grid-cols-2 gap-x-4 gap-y-8 sm:grid-cols-3 sm:gap-x-6 md:grid-cols-4 lg:grid-cols-3 xl:grid-cols-4 xl:gap-x-8"
	>
		{#each media as file (file?.id)}
			<li class="relative">
				<div
					class="group block w-full overflow-hidden rounded-lg bg-gray-100 ring-1 ring-gray-300 ring-inset"
				>
					<img
						src={file?.thumbnail_url ?? file?.url}
						alt={file.alt_text || ''}
						class="pointer-events-none aspect-video w-full object-cover group-hover:opacity-75"
					/>
					<button type="button" class="absolute inset-0 focus:outline-none">
						<span class="sr-only">View details for {file.title}</span>
					</button>
				</div>
				<p class="pointer-events-none mt-2 block truncate text-sm font-medium text-gray-900">
					{file.title}
				</p>
				<p class="pointer-events-none block text-sm font-medium text-gray-500">
					{file.file_size ?? ''}
				</p>
			</li>
		{/each}
	</ul>
</section>
