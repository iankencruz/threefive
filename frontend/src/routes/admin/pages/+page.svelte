<!-- frontend/src/routes/admin/pages/+page.svelte -->
<script lang="ts">
	import { goto } from "$app/navigation";
	import type { PageData } from "./$types";
	import { PUBLIC_API_URL } from "$env/static/public";
	import { browser } from "$app/environment";
	import { EyeIcon, SquarePenIcon } from "lucide-svelte";

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

	const getTypeColor = (type: string) => {
		switch (type) {
			case "project":
				return "bg-blue-100 text-blue-800";
			case "blog":
				return "bg-purple-100 text-purple-800";
			case "generic":
				return "bg-gray-100 text-gray-800";
			default:
				return "bg-gray-100 text-gray-800";
		}
	};

	function navigateToExternal(url: string) {
		if (browser) {
			// Ensure this runs only in the browser environment
			window.location.href = url;
		}
	}
</script>

<div class="max-w-7xl mx-auto">
	<div class="flex justify-between items-center mb-8">
		<h1 class="">Pages</h1>
		<button
			class="flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium transition-colors"
			onclick={() => goto('/admin/pages/new')}
		>
			<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
			</svg>
			New Page
		</button>
	</div>

	{#if !data.pages || data.pages.length === 0}
		<div class="flex flex-col items-center justify-center py-16 px-8 bg-surface rounded-lg border-2 border-dashed border-gray-300">
			<svg class="w-16 h-16 text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
				/>
			</svg>
			<h3 class="text-xl font-semibold text-gray-900 mb-2">No pages yet</h3>
			<p class="text-gray-600 mb-6">Get started by creating your first page</p>
			<button
				class="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg font-medium transition-colors"
				onclick={() => goto('/admin/pages/new')}
			>
				Create Page
			</button>
		</div>
	{:else}
		<div class="bg-surface rounded-lg shadow overflow-hidden">
			<table class="min-w-full divide-y divide-gray-700">
				<thead class="bg-surface">
					<tr>
						<th class="px-6 py-3 text-left text-xs font-medium  uppercase tracking-wider">
							Title
						</th>
						<th class="px-6 py-3 text-left text-xs font-medium   uppercase tracking-wider">
							Type
						</th>
						<th class="px-6 py-3 text-left text-xs font-medium   uppercase tracking-wider">
							Status
						</th>
						<th class="px-6 py-3 text-left text-xs font-medium   uppercase tracking-wider">
							Author
						</th>
						<th class="px-6 py-3 text-left text-xs font-medium   uppercase tracking-wider">
							Updated
						</th>
						<th class="px-6 py-3 text-left text-xs font-medium   uppercase tracking-wider">
							Actions
						</th>
					</tr>
				</thead>
				<tbody class="bg-surface divide-y divide-gray-200">
					{#each data.pages as page}
						<tr class="">
							<td class="px-6 py-4 whitespace-nowrap">
								<div class="flex flex-col">
									<span class="font-medium ">{page.title}</span>
									<span class="text-sm  ">/{page.slug}</span>
								</div>
							</td>
							<td class="px-6 py-4 whitespace-nowrap">
								<span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium capitalize {getTypeColor(page.page_type)}">
									{page.page_type}
								</span>
							</td>
							<td class="px-6 py-4 whitespace-nowrap">
								<span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium capitalize {getStatusColor(page.status)}">
									{page.status}
								</span>
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm ">
								{page.author_id}
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm ">
								{formatDate(page.updated_at)}
							</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <div class="flex items-center gap-2">
                      <!-- Preview Button -->
                      <button
                          class="p-2   hover:bg-gray-400 rounded-lg transition-colors"
                          onclick={() => window.open(`/admin/pages/${page.id}/preview`, '_blank')}
                          aria-label="Preview page"
                          title="Preview in new tab"
                      >
                          <EyeIcon size={20}/>
                      </button>
                      
                      <!-- Edit Button -->
                      <button
                          class="p-2   hover:bg-gray-400 rounded-lg transition-colors"
                          onclick={() => goto(`/admin/pages/${page.id}/edit`)}
                          aria-label="Edit page"
                      >
                          <SquarePenIcon size={16}/>
                      </button>
                      
                      <!-- View Public Page (only if published) -->
                      {#if page.status === 'published'}
                          <button
                              class="p-2 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-colors"
                              onclick={() => navigateToExternal(`/${page.slug}`)}
                              aria-label="View page"
                          >
                      {console.log("public page: ", `${PUBLIC_API_URL}/${page.slug}`)}
                              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                  <path
                                      stroke-linecap="round"
                                      stroke-linejoin="round"
                                      stroke-width="2"
                                      d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
                                  />
                              </svg>
                          </button>
                      {/if}
                  </div>
              </td>					

            </tr>
					{/each}
				</tbody>
			</table>
		</div>

		{#if data.pagination && data.pagination.total_pages > 1}
			<div class="flex items-center justify-center gap-4 mt-8">
				<button
					class="px-4 py-2 border border-gray-300 rounded-lg bg-surface text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed font-medium transition-colors"
					disabled={data.pagination.page === 1}
					onclick={() => goto(`/admin/pages?page=${data.pagination.page - 1}`)}
				>
					Previous
				</button>
				<span class="text-sm text-gray-600">
					Page {data.pagination.page} of {data.pagination.total_pages}
				</span>
				<button
					class="px-4 py-2 border border-gray-300 rounded-lg bg-surface text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed font-medium transition-colors"
					disabled={data.pagination.page === data.pagination.total_pages}
					onclick={() => goto(`/admin/pages?page=${data.pagination.page + 1}`)}
				>
					Next
				</button>
			</div>
		{/if}
	{/if}
</div>
