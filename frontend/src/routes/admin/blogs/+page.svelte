<!-- frontend/src/routes/admin/blogs/+page.svelte -->
<script lang="ts">
	import { goto } from "$app/navigation";
	import type { PageData } from "./$types";
	import { EyeIcon, SquarePenIcon, Layers } from "lucide-svelte";
	import { page } from "$app/state";

	let { data }: { data: PageData } = $props();

	const formatDate = (dateString: string) => {
		return new Date(dateString).toLocaleDateString("en-US", {
			year: "numeric",
			month: "short",
			day: "numeric",
		});
	};

	const getStatusColor = (status: string) => {
		switch (status) {
			case "published":
				return "bg-green-100 text-green-800";
			case "draft":
				return "bg-yellow-100 text-yellow-800";
			case "archived":
				return "bg-gray-100 text-gray-800";
			default:
				return "bg-gray-100 text-gray-800";
		}
	};
</script>


{#snippet pagination(data: PageData)}
		{#if data.pagination && data.pagination.total_pages > 1}
			<div class="flex items-center justify-center gap-4 mt-8">
				<button
					class="px-4 py-2 border border-gray-300 rounded-lg bg-surface text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed font-medium transition-colors"
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
					class="px-4 py-2 border border-gray-300 rounded-lg bg-surface text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed font-medium transition-colors"
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



<div class="max-w-7xl mx-auto">
	<div class="flex justify-between items-center mb-8">
		<h1 class="">Blogs</h1>
		<button
			class="flex items-center gap-2 bg-primary hover:bg-primary/90 text-white px-4 py-2 rounded-lg font-medium transition-colors"
			onclick={() => goto('/admin/blogs/new')}
		>
			<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M12 4v16m8-8H4"
				/>
			</svg>
			New Blog Post
		</button>
	</div>

	{#if data.blogs.length === 0}
		<!-- Empty State -->
		<div
			class="text-center py-20 bg-surface rounded-lg border-2 border-dashed border-gray-700"
		>

      <Layers class="mx-auto h-12 w-12 text-gray-400 mb-4"/>
			<h3 class="text-lg font-medium text-gray-200 mb-2">No blog posts yet</h3>
			<p class="text-gray-400 mb-6">Get started by creating your first blog post.</p>
			<button
				onclick={() => goto('/admin/blogs/new')}
				class="inline-flex items-center gap-2 bg-primary hover:bg-primary/90 text-white px-4 py-2 rounded-lg font-medium transition-colors"
			>
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M12 4v16m8-8H4"
					/>
				</svg>
				Create Blog Post
			</button>
		</div>
	{:else}
		<!-- Blogs Table -->
		<div class="bg-surface rounded-lg shadow overflow-hidden">
			<table class="min-w-full divide-y divide-gray-700">
				<thead class="bg-surface">
					<tr>
						<th
							class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider"
						>
							Title
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider"
						>
							Excerpt
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider"
						>
							Reading Time
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider"
						>
							Status
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider"
						>
							Updated
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider"
						>
							Actions
						</th>
					</tr>
				</thead>
				<tbody class="bg-surface divide-y divide-gray-700">
					{#each data.blogs as blog}
						<tr class="hover:bg-white/5">
							<td class="px-6 py-4">
								<div class="flex flex-col">
									<div class="flex items-center gap-2">
										<span class="font-medium">{blog.title}</span>
										{#if blog.is_featured}
											<span
												class="px-2 py-0.5 bg-yellow-900/30 text-yellow-300 text-xs rounded"
											>
												Featured
											</span>
										{/if}
									</div>
									<span class="text-sm text-gray-400">/{blog.slug}</span>
								</div>
							</td>
							<td class="px-6 py-4">
								{#if blog.excerpt}
									<p class="text-sm text-gray-300 line-clamp-2 max-w-xs">
										{blog.excerpt}
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
									class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium capitalize {getStatusColor(
										blog.status,
									)}"
								>
									{blog.status}
								</span>
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm">
								{formatDate(blog.updated_at)}
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
								<div class="flex items-center gap-2">
									<button
										class="p-2 hover:bg-gray-700 rounded-lg transition-colors"
										onclick={() =>
											window.open(`/preview/pages/${blog.id}`, "_blank")}
										title="Preview"
									>
										<EyeIcon class="w-4 h-4" />
									</button>
									<button
										class="p-2 hover:bg-gray-700 rounded-lg transition-colors"
										onclick={() => goto(`/admin/blogs/${blog.id}/edit`)}
										title="Edit"
									>
										<SquarePenIcon class="w-4 h-4" />
									</button>
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
