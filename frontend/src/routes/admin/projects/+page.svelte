<!-- frontend/src/routes/admin/projects/+page.svelte -->
<script lang="ts">
	import { goto } from '$app/navigation';
	import type { PageData } from './$types';
	import { EyeIcon, Layers, SquarePenIcon } from 'lucide-svelte';
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
		<h1 class="">Projects</h1>
		<button
			class="flex items-center gap-2 rounded-lg bg-primary px-4 py-2 font-medium text-white transition-colors hover:bg-primary/90"
			onclick={() => goto('/admin/projects/new')}
		>
			<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
			</svg>
			New Project
		</button>
	</div>

	{#if !data.projects.data || data.projects.data.length === 0}
		<!-- Empty State -->
		<div class="rounded-lg border-2 border-dashed border-gray-700 bg-surface py-20 text-center">
			<Layers class="mx-auto mb-4 h-12 w-12 text-gray-400" />
			<h3 class="mb-2 text-lg font-medium text-gray-200">No projects yet</h3>
			<p class="mb-6 text-gray-400">Get started by creating your first project.</p>
			<button
				onclick={() => goto('/admin/projects/new')}
				class="inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2 font-medium text-white transition-colors hover:bg-primary/90"
			>
				<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M12 4v16m8-8H4"
					/>
				</svg>
				Create Project
			</button>
		</div>
	{:else}
		<!-- Projects Table -->
		<div class="overflow-hidden rounded-lg bg-surface shadow">
			<table class="min-w-full divide-y divide-gray-700">
				<thead class="bg-surface">
					<tr>
						<th class="px-6 py-3 text-left text-xs font-medium tracking-wider uppercase">
							Title
						</th>
						<th class="px-6 py-3 text-left text-xs font-medium tracking-wider uppercase">
							Client
						</th>
						<th class="px-6 py-3 text-left text-xs font-medium tracking-wider uppercase">
							Technologies
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
					{#each data.projects.data as project}
						<tr class="hover:bg-white/5">
							<td class="px-6 py-4">
								<div class="flex flex-col">
									<span class="font-medium">{project.title}</span>
									<span class="text-sm text-gray-400">/{project.slug}</span>
								</div>
							</td>
							<td class="px-6 py-4 text-sm text-gray-300">
								{project.client_name || '—'}
							</td>
							<td class="px-6 py-4">
								{#if project.technologies && project.technologies.length > 0}
									<div class="flex flex-wrap gap-1">
										{#each project.technologies.slice(0, 3) as tech}
											<span class="rounded bg-blue-900/30 px-2 py-0.5 text-xs text-blue-300">
												{tech}
											</span>
										{/each}
										{#if project.technologies.length > 3}
											<span class="px-2 py-0.5 text-xs text-gray-400">
												+{project.technologies.length - 3}
											</span>
										{/if}
									</div>
								{:else}
									<span class="text-sm text-gray-500">—</span>
								{/if}
							</td>
							<td class="px-6 py-4 whitespace-nowrap">
								<span
									class="inline-flex items-center rounded-full px-3 py-1 text-xs font-medium capitalize {getStatusColor(
										project.status
									)}"
								>
									{project.status}
								</span>
							</td>
							<td class="px-6 py-4 text-sm whitespace-nowrap">
								{formatDate(project.updated_at)}
							</td>
							<td class="px-6 py-4 text-sm font-medium whitespace-nowrap">
								<div class="flex items-center gap-2">
									<button
										class="rounded-lg p-2 transition-colors hover:bg-gray-700"
										onclick={() => window.open(`/preview/pages/${project.id}`, '_blank')}
										title="Preview"
									>
										<EyeIcon class="h-4 w-4" />
									</button>
									<button
										class="rounded-lg p-2 transition-colors hover:bg-gray-700"
										onclick={() => goto(`/admin/projects/${project.id}/edit`)}
										title="Edit"
									>
										<SquarePenIcon class="h-4 w-4" />
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
