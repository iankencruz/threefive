<script lang="ts">
	let { page, totalPages, onchange, pageSize, totalMedia } = $props<{
		page: number;
		totalPages: number;
		pageSize?: number;
		totalMedia?: number;
		onchange: (newPage: number) => void;
	}>();

	// Set default if value missing
	pageSize ??= 10;

	let start = $state<number>(totalMedia === 0 ? 0 : (page - 1) * pageSize + 1);
	let end = $derived<number>(Math.min(start + pageSize - 1, totalMedia ?? 0));

	$effect(() => {
		start = totalMedia === 0 ? 0 : (page - 1) * pageSize + 1;
		end = Math.min(start + pageSize - 1, totalMedia ?? 0);
	});

	function paginate(newPage: number) {
		if (newPage >= 1 && newPage <= totalPages && newPage !== page) {
			start = (newPage - 1) * pageSize + 1;
			end = Math.min(start + pageSize - 1, totalMedia ?? 0);
			onchange(newPage);
		}
	}
	function getVisiblePages(): (number | 'dots')[] {
		const pages: (number | 'dots')[] = [];

		if (totalPages <= 7) {
			for (let i = 1; i <= totalPages; i++) pages.push(i);
			return pages;
		}

		pages.push(1); // Always show first page

		if (page > 4) pages.push('dots');

		const start = Math.max(2, page - 1);
		const end = Math.min(totalPages - 1, page + 1);

		for (let i = start; i <= end; i++) {
			pages.push(i);
		}

		if (page < totalPages - 3) pages.push('dots');

		pages.push(totalPages); // Always show last page

		return pages;
	}
</script>

<div
	class=" mt-8 flex w-full items-center justify-between border-t border-gray-200 bg-white px-4 py-3 sm:px-6"
>
	<div class="flex flex-1 justify-between sm:hidden">
		<button
			onclick={() => paginate(page - 1)}
			class="relative inline-flex items-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
			>Previous</button
		>
		<button
			onclick={() => paginate(page + 1)}
			class="relative ml-3 inline-flex items-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
			>Next</button
		>
	</div>
	<div class="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
		<div>
			<p class="text-sm text-gray-700">
				Showing
				<span class="font-medium">{start}</span>
				to
				<span class="font-medium">{end}</span>
				of
				<span class="font-medium">{totalMedia}</span>
				results
			</p>
		</div>
		<div>
			<nav class="isolate inline-flex -space-x-px rounded-md shadow-xs" aria-label="Pagination">
				<button
					onclick={() => paginate(page - 1)}
					class="relative inline-flex items-center rounded-l-md px-2 py-2 text-gray-400 ring-1 ring-gray-300 ring-inset hover:bg-gray-50 focus:z-20 focus:outline-offset-0"
				>
					<span class="sr-only">Previous</span>
					<svg
						class="size-5"
						viewBox="0 0 20 20"
						fill="currentColor"
						aria-hidden="true"
						data-slot="icon"
					>
						<path
							fill-rule="evenodd"
							d="M11.78 5.22a.75.75 0 0 1 0 1.06L8.06 10l3.72 3.72a.75.75 0 1 1-1.06 1.06l-4.25-4.25a.75.75 0 0 1 0-1.06l4.25-4.25a.75.75 0 0 1 1.06 0Z"
							clip-rule="evenodd"
						/>
					</svg>
				</button>
				<!-- Current: "z-10 bg-indigo-600 text-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600", Default: "text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:outline-offset-0" -->

				{#each getVisiblePages() as p}
					{#if p === 'dots'}
						<span
							class="relative inline-flex items-center px-4 py-2 text-sm font-semibold text-gray-700 ring-1 ring-gray-300 ring-inset"
							>...</span
						>
					{:else}
						<button
							onclick={() => paginate(p)}
							class={`relative inline-flex items-center px-4 py-2 text-sm font-semibold ring-1 ring-gray-300 ring-inset focus:z-20 ${p === page ? 'z-10 bg-indigo-600 text-white' : 'text-gray-900 hover:bg-gray-50'}`}
							aria-current={p === page ? 'page' : undefined}
						>
							{p}
						</button>
					{/if}
				{/each}
				<button
					onclick={() => paginate(page + 1)}
					class="relative inline-flex items-center rounded-r-md px-2 py-2 text-gray-400 ring-1 ring-gray-300 ring-inset hover:bg-gray-50 focus:z-20 focus:outline-offset-0"
				>
					<span class="sr-only">Next</span>
					<svg
						class="size-5"
						viewBox="0 0 20 20"
						fill="currentColor"
						aria-hidden="true"
						data-slot="icon"
					>
						<path
							fill-rule="evenodd"
							d="M8.22 5.22a.75.75 0 0 1 1.06 0l4.25 4.25a.75.75 0 0 1 0 1.06l-4.25 4.25a.75.75 0 0 1-1.06-1.06L11.94 10 8.22 6.28a.75.75 0 0 1 0-1.06Z"
							clip-rule="evenodd"
						/>
					</svg>
				</button>
			</nav>
		</div>
	</div>
</div>
