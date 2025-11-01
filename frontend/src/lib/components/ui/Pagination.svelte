<script lang="ts">
	interface Props {
		currentPage: number;
		totalPages: number;
		onPageChange: (page: number) => void;
	}

	let { currentPage, totalPages, onPageChange }: Props = $props();

	function changePage(page: number) {
		if (page < 1 || page > totalPages) return;
		onPageChange(page);
	}
</script>

{#if totalPages > 1}
	<nav aria-label="Pagination" class="isolate inline-flex -space-x-px rounded-md">
		<!-- Previous Button -->
		<button
			onclick={() => changePage(currentPage - 1)}
			disabled={currentPage === 1}
			class="relative inline-flex items-center rounded-l-md px-2 py-1 text-gray-400 inset-ring inset-ring-gray-700 hover:bg-white/5 focus:z-20 focus:outline-offset-0 disabled:opacity-30 disabled:cursor-not-allowed"
		>
			<span class="sr-only">Previous</span>
			<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
			</svg>
		</button>

		<!-- Generate page numbers to show -->
		{#each (() => {
			const pages = [];
			
			// Always add first page
			pages.push(1);
			
			// Add ellipsis if needed before middle pages
			if (currentPage > 3) {
				pages.push('...');
			}
			
			// Add middle pages (current - 1, current, current + 1)
			// But exclude page 1 and last page to avoid duplicates
			for (let i = Math.max(2, currentPage - 1); i <= Math.min(totalPages - 1, currentPage + 1); i++) {
				if (i < totalPages) { // Don't add if it's the last page
					pages.push(i);
				}
			}
			
			// Add ellipsis if needed after middle pages
			if (currentPage < totalPages - 2) {
				pages.push('...');
			}
			
			// Always add last page (if different from first page)
			if (totalPages > 1) {
				pages.push(totalPages);
			}
			
			return pages;
		})() as page}
			{#if page === '...'}
				<span class="relative inline-flex items-center px-4 py-1 text-sm font-semibold text-gray-400 inset-ring inset-ring-gray-700">
					...
				</span>
			{:else}
				<button
					onclick={() => changePage(page)}
					class="relative inline-flex items-center px-4 py-1 text-sm font-semibold focus:z-20 {currentPage === page 
						? 'bg-primary text-white' 
						: 'text-gray-300 inset-ring inset-ring-gray-700 hover:bg-white/5'}"
				>
					{page}
				</button>
			{/if}
		{/each}

		<!-- Next Button -->
		<button
			onclick={() => changePage(currentPage + 1)}
			disabled={currentPage === totalPages}
			class="relative inline-flex items-center rounded-r-md px-2 py-1 text-gray-400 inset-ring inset-ring-gray-700 hover:bg-white/5 focus:z-20 focus:outline-offset-0 disabled:opacity-30 disabled:cursor-not-allowed"
		>
			<span class="sr-only">Next</span>
			<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
			</svg>
		</button>
	</nav>
{/if}
