<!-- frontend/src/routes/admin/pages/+page.svelte -->
<script lang="ts">
	import { goto } from '$app/navigation';
	import type { PageData } from './$types';
	import { browser } from '$app/environment';
	import { EyeIcon, SquarePenIcon, Layers, Plus } from 'lucide-svelte';
	import { page } from '$app/state';

	let { data }: { data: PageData } = $props();

	// Get current page_type from URL
	const currentPageType = $derived(page.url.searchParams.get('page_type') || 'all');

	const formatDate = (dateString: string) => {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric'
		});
	};

	const getStatusColor = (status: string) => {
		switch (status) {
			case 'published':
				return 'bg-green-100 text-green-800';
			case 'draft':
				return 'bg-yellow-100 text-yellow-800';
			case 'archived':
				return 'bg-gray-100 text-gray-800';
			default:
				return 'bg-gray-100 text-gray-800';
		}
	};

	const getTypeColor = (type: string) => {
		switch (type) {
			case 'project':
				return 'color-project';
			case 'blog':
				return 'color-blog';
			case 'generic':
				return 'color-generic';
			default:
				return 'color-generic';
		}
	};

	

	const headerItems = ['title', 'status', 'updated', 'actions'];
</script>

{#snippet pagination(data: PageData)}
	{#if data.pagination && data.pagination.total_pages > 1}
		<div class="mt-8 flex items-center justify-center gap-4">
			<button
				class="rounded-lg border border-gray-300 bg-surface px-4 py-2 font-medium text-gray-700 transition-colors hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50"
				disabled={data.pagination.page === 1}
				onclick={() => {
					const params = new URLSearchParams(page.url.searchParams);
					params.set('page', (data.pagination.page - 1).toString());
					goto(`/admin/pages?${params.toString()}`);
				}}
			>
				Previous
			</button>
			<span class="text-sm text-gray-600">
				Page {data.pagination.page} of {data.pagination.total_pages}
			</span>
			<button
				class="rounded-lg border border-gray-300 bg-surface px-4 py-2 font-medium text-gray-700 transition-colors hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50"
				disabled={data.pagination.page === data.pagination.total_pages}
				onclick={() => {
					const params = new URLSearchParams(page.url.searchParams);
					params.set('page', (data.pagination.page + 1).toString());
					goto(`/admin/pages?${params.toString()}`);
				}}
			>
				Next
			</button>
		</div>
	{/if}
{/snippet}

{#snippet thead(header: string[])}
	<thead class="bg-surface">
		<tr>
			{#each header as h}
				<th class="bg-black px-6 py-3 text-left text-xs font-medium tracking-wider uppercase">
					{h}
				</th>
			{/each}
		</tr>
	</thead>
{/snippet}

<div class="mx-auto max-w-7xl">
	<div class="mb-8 flex items-center justify-between">
		<h1 class="">Pages</h1>
		<button
			class="flex items-center gap-2 rounded-lg bg-primary px-4 py-2 font-medium text-white transition-colors hover:bg-primary"
			onclick={() => goto('/admin/pages/new')}
		>
			<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
			</svg>
			New Page
		</button>
	</div>

	<!-- Table -->
	{#if !data.pages || data.pages.data.length === 0}
		<div
			class="flex flex-col items-center justify-center rounded-lg border-2 border-dashed border-foreground-muted bg-surface px-8 py-16"
		>
			<Layers class="mx-auto mb-4 h-12 w-12 text-gray-400" />
			<h3 class="mb-2 text-xl font-semibold">No pages yet</h3>
			<p class=" mb-6">Get started by creating your first page</p>
			<button
				class="inline-flex items-center gap-2 cursor-pointer rounded-lg bg-primary px-6 py-2 font-medium text-white transition-colors hover:bg-primary"
				onclick={() => goto('/admin/pages/new')}
			>
				<Plus size={18} />
				Create Page
			</button>
		</div>
	{:else}
		<div class="overflow-hidden rounded-lg bg-surface shadow">
			<table class="min-w-full divide-y divide-gray-700">
				{@render thead(headerItems)}
				<tbody class="divide-y divide-gray-200 bg-surface">
					{#each data.pages.data as page}
						<tr class="">
							<td class="px-6 py-4 whitespace-nowrap">
								<div class="flex flex-col">
									<span class="font-medium">{page.title}</span>
									<span class="text-sm">/{page.slug}</span>
								</div>
							</td>
							<td class="px-6 py-4 whitespace-nowrap">
								<span
									class="inline-flex items-center rounded-full px-3 py-1 text-xs font-medium capitalize {getStatusColor(
										page.status
									)}"
								>
									{page.status}
								</span>
							</td>
							<td class="px-6 py-4 text-sm whitespace-nowrap">
								{formatDate(page.updated_at)}
							</td>
							<td class="px-6 py-4 text-sm font-medium whitespace-nowrap">
								<div class="flex items-center gap-2">
									<!-- Preview Button -->
									<a
										class="rounded-lg p-2 transition-colors hover:bg-gray-400"
										href={`/admin/pages/${page.id}/preview`}
                    target="_blank"
										aria-label="Preview page"
										title="Preview in new tab"
									>
										<EyeIcon size={20} />
									</a>

									<!-- Edit Button -->
									<button
										class="rounded-lg p-2 transition-colors hover:bg-gray-400"
										onclick={() => goto(`/admin/pages/${page.id}/edit`)}
										aria-label="Edit page"
									>
										<SquarePenIcon size={16} />
									</button>

									<!-- View Public Page (only if published) -->
									{#if page.status === 'published'}
										<a
											class="rounded-lg p-2 text-gray-600 transition-colors hover:bg-gray-100 hover:text-gray-900"
                      href={`/${page.slug}`}
											aria-label="View page"
										>
											<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path
													stroke-linecap="round"
													stroke-linejoin="round"
													stroke-width="2"
													d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
												/>
											</svg>
										</a>
									{/if}
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		{@render pagination(data)}
	{/if}
</div>
