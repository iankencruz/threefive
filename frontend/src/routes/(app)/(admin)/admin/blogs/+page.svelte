<!-- frontend/src/routes/admin/blogs/+page.svelte -->
<script lang="ts">
	import { goto } from '$app/navigation';
	import type { PageData } from './$types';
	import { EyeIcon, SquarePenIcon, Layers, Plus } from 'lucide-svelte';
	import { page } from '$app/state';

	let { data }: { data: PageData } = $props();

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
</script>

<svelte:head>
	<title>Admin: Blogs</title>
</svelte:head>

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

<div class="mx-auto max-w-7xl">
	<div class="mb-8 flex items-center justify-between">
		<h1 class="">Blogs</h1>
		<button
			class="flex items-center gap-2 rounded-lg bg-primary px-4 py-2 font-medium text-white transition-colors hover:bg-primary/90"
			onclick={() => goto('/admin/blogs/new')}
		>
			<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
			</svg>
			New Blog Post
		</button>
	</div>

	{#if data.blogs.length === 0}
		<!-- Empty State -->
		<div class="rounded-lg border-2 border-dashed border-gray-700 bg-surface py-20 text-center">
			<Layers class="mx-auto mb-4 h-12 w-12 text-gray-400" />
			<h3 class="mb-2 text-lg font-medium text-gray-200">No blog posts yet</h3>
			<p class="mb-6 text-gray-400">Get started by creating your first blog post.</p>
			<button
				onclick={() => goto('/admin/blogs/new')}
				class="inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2 font-medium text-white transition-colors hover:bg-primary/90"
			>
				<Plus size={18} />
				Create Blog Post
			</button>
		</div>
	{:else}
		<!-- Blogs Table -->
		<div class="overflow-hidden rounded-lg bg-surface shadow">
			<table class="min-w-full divide-y divide-gray-700">
				<thead class="bg-surface">
					<tr>
						<th class="px-6 py-3 text-left text-xs font-medium tracking-wider uppercase">
							Title
						</th>
						<th class="px-6 py-3 text-left text-xs font-medium tracking-wider uppercase">
							Description
						</th>
						<th class="px-6 py-3 text-left text-xs font-medium tracking-wider uppercase">
							Reading Time
						</th>
						<th class="px-6 py-3 text-left text-xs font-medium tracking-wider uppercase">
							Status
						</th>
						<th class="px-6 py-3 text-left text-xs font-medium tracking-wider uppercase">
							Updated
						</th>
						<th class="px-6 py-3 text-left text-xs font-medium tracking-wider uppercase">
							Actions
						</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-700 bg-surface">
					{#each data.blogs as blog}
						<tr class="hover:bg-white/5">
							<td class="px-6 py-4">
								<div class="flex flex-col">
									<div class="flex items-center gap-2">
										<span class="font-medium">{blog.title}</span>
										{#if blog.is_featured}
											<span class="rounded bg-yellow-900/30 px-2 py-0.5 text-xs text-yellow-300">
												Featured
											</span>
										{/if}
									</div>
									<span class="text-sm text-gray-400">/{blog.slug}</span>
								</div>
							</td>
							<td class="px-6 py-4">
								{#if blog.description}
									<p class="line-clamp-2 max-w-xs text-sm text-gray-300">
										{blog.description}
									</p>
								{:else}
									<span class="text-sm text-gray-500">—</span>
								{/if}
							</td>
							<td class="px-6 py-4 text-sm text-gray-300">
								{#if blog.reading_time}
									{blog.reading_time} min
								{:else}
									—
								{/if}
							</td>
							<td class="px-6 py-4 whitespace-nowrap">
								<span
									class="inline-flex items-center rounded-full px-3 py-1 text-xs font-medium capitalize {getStatusColor(
										blog.status
									)}"
								>
									{blog.status}
								</span>
							</td>
							<td class="px-6 py-4 text-sm whitespace-nowrap">
								{formatDate(blog.updated_at)}
							</td>
							<td class="px-6 py-4 text-sm font-medium whitespace-nowrap">
								<div class="flex items-center gap-2">
									<a
										class="rounded-lg p-2 transition-colors hover:bg-gray-700"
										href={`/admin/blogs/${blog.id}/preview`}
										title="Preview"
									>
										<EyeIcon class="h-4 w-4" />
									</a>
									<button
										class="rounded-lg p-2 transition-colors hover:bg-gray-700"
										onclick={() => goto(`/admin/blogs/${blog.id}/edit`)}
										title="Edit"
									>
										<SquarePenIcon class="h-4 w-4" />
									</button>

									<!-- View Public Page (only if published) -->
									{#if blog.status === 'published'}
										<a
											class="rounded-lg p-2 text-gray-600 transition-colors hover:bg-gray-100 hover:text-gray-900"
											href={`/blogs/${blog.slug}`}
											target="_blank"
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
