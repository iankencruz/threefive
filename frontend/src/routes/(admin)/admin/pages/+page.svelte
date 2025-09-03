<script lang="ts">
	import { toast } from 'svelte-sonner';
	import PageForm from '$lib/components/Forms/PageForm.svelte';
	import type { Page } from '$lib/types';
	import { getPages } from '$lib/api/pages';
	import { goto } from '$app/navigation';
	import EmptyState from '$lib/components/Overlays/EmptyState.svelte';
	import { auth } from '$lib/store/auth.svelte';
	import { formatDate, slugify } from '$lib/utils/utilities';
	import { onMount } from 'svelte';
	import Drawers from '$lib/components/Overlays/Drawers.svelte';
	import { page } from '$app/state';
	import { ChevronDownIcon, ChevronUpIcon } from '@lucide/svelte';

	let newPage = $state<Page>({
		id: crypto.randomUUID() as Page['id'],
		slug: '',
		title: '',
		cover_image_id: null,
		seo_title: '',
		seo_description: '',
		seo_canonical: '',
		content: '',
		is_draft: true,
		is_published: false,
		created_at: new Date().toISOString(),
		updated_at: new Date().toISOString()
	});

	// track current sort state
	let sortField = $state<'title' | 'status' | 'created_at' | 'updated_at'>('updated_at');
	let sortDirection = $state<'asc' | 'desc'>('desc');

	onMount(() => {
		const sortParam = page.url.searchParams.get('sort') || 'desc';
		pageContent = getPages(sortParam as 'asc' | 'desc');
	});

	// when URL changes, re-fetch
	$effect(() => {
		const sortParam = page.url.searchParams.get('sort');
		if (sortParam) {
			const [f, d] = sortParam.split(':');
			if (f === 'title' || f === 'created_at' || f === 'updated_at') {
				sortField = f;
			} else {
				sortField = 'updated_at';
			}
			sortDirection = d === 'asc' ? 'asc' : 'desc';
		}
		pageContent = getPages(`${sortField}:${sortDirection}`);
	});

	$effect(() => {
		newPage.slug = slugify(newPage.title);
	});

	function toggleSort(field: typeof sortField) {
		if (sortField === field) {
			// flip direction if same field
			sortDirection = sortDirection === 'asc' ? 'desc' : 'asc';
		} else {
			sortField = field;
			sortDirection = 'asc'; // default new field to asc
		}
		goto(`/admin/pages?sort=${sortField}:${sortDirection}`, { replaceState: true });
	}

	let drawerOpen = $state(false);

	$inspect('newPage:', newPage);

	let pageContent = $state<Promise<Page[]> | null>(null);

	function openCreateDrawer() {
		drawerOpen = true;
	}

	function closeDrawer() {
		drawerOpen = false;
	}

	async function handleCreate() {
		try {
			const res = await fetch('/api/v1/admin/pages', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(newPage)
			});

			const json = await res.json();
			if (!res.ok) {
				toast.error(json.message || 'Failed to save');
				return;
			}

			drawerOpen = false;
			toast.success('Page created');
			pageContent = getPages(); // ✅ refresh pages list
		} catch (err) {
			console.error(err);
			toast.error('Create failed');
		}
	}

	async function handleDelete(slug: string) {
		if (!confirm('Delete project?')) return;

		try {
			const res = await fetch(`/api/v1/admin/pages/${slug}`, {
				method: 'DELETE'
			});

			if (!res.ok) {
				throw new Error('Failed to delete project');
			}
			toast.success('Project deleted');
			pageContent = getPages();
		} catch (error) {
			console.error('❌ Failed to delete project', error);
			toast.error('Failed to delete project. Please try again.');
		}
	}
</script>

{#await pageContent}
	<p>Loading pages...</p>
{:then data}
	<section class="py-6">
		<div class="mb-6 flex items-center justify-between">
			<h1 class=" text-2xl font-semibold text-gray-900">Pages</h1>
			<button
				onclick={openCreateDrawer}
				class="rounded-md border bg-indigo-600 px-3 py-1.5 text-sm text-black text-white hover:bg-black/10"
			>
				+ New Page
			</button>
		</div>

		{#if !data}
			<!-- {#snippet icon()} -->
			<!-- 	<PaperclipIcon /> -->
			<!-- {/snippet} -->
			<div class="mt-12">
				<EmptyState
					action={openCreateDrawer}
					title={'No Projects'}
					text={'Get Started By Creating a new project.'}
				/>
			</div>
		{:else}
			<table class="relative min-w-full divide-y divide-gray-300">
				<thead>
					<tr>
						<th
							scope="col"
							class="cursor-pointer px-3 py-3.5 text-left text-sm font-semibold text-gray-900"
							onclick={() => toggleSort('title')}
						>
							Title
							{#if sortField === 'title'}
								<span>{sortDirection === 'asc' ? '▲' : '▼'}</span>
							{/if}
						</th>
						<th
							scope="col"
							class="cursor-pointer px-3 py-3.5 text-left text-sm font-semibold text-gray-900"
							onclick={() => toggleSort('status')}
						>
							Status
							{#if sortField === 'status'}
								<span>{sortDirection === 'asc' ? '▲' : '▼'}</span>
							{/if}
						</th>
						<th
							scope="col"
							class="cursor-pointer px-3 py-3.5 text-left text-sm font-semibold text-gray-900"
							onclick={() => toggleSort('created_at')}
						>
							Created At
							{#if sortField === 'created_at'}
								<span>{sortDirection === 'asc' ? '▲' : '▼'}</span>
							{/if}
						</th>
						<th
							scope="col"
							class="cursor-pointer px-3 py-3.5 text-left text-sm font-semibold text-gray-900"
							onclick={() => toggleSort('updated_at')}
						>
							Updated At
							{#if sortField === 'updated_at'}
								<span>{sortDirection === 'asc' ? '▲' : '▼'}</span>
							{/if}
						</th>
						<th scope="col" class="py-3.5 pr-0 pl-3">
							Actions
							<span class="sr-only">Action</span>
						</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-200 bg-white">
					{#each data as page}
						<tr>
							<td class="py-4 pr-3 pl-4 text-sm font-medium whitespace-nowrap text-gray-900 sm:pl-0"
								>{page.title}</td
							>
							<td class="px-3 py-4 text-sm whitespace-nowrap text-gray-500"
								>{#if page.is_published}
									<p
										class="mt-0.5 w-fit rounded-md bg-green-50 px-1.5 py-0.5 text-xs font-medium whitespace-nowrap text-green-700 ring-1 ring-green-600/20 ring-inset"
									>
										Published
									</p>
								{:else}
									<p
										class="mt-0.5 w-fit rounded-md bg-yellow-50 px-1.5 py-0.5 text-xs font-medium whitespace-nowrap text-yellow-800 ring-1 ring-yellow-600/20 ring-inset"
									>
										Draft
									</p>
								{/if}</td
							>
							<td class="px-3 py-4 text-sm whitespace-nowrap text-gray-500"
								><time datetime="2023-03-17T00:00Z">{formatDate(page.created_at, 'relative')}</time>
							</td>
							<td class="px-3 py-4 text-sm whitespace-nowrap text-gray-500"
								><time datetime="2023-03-17T00:00Z">{formatDate(page.updated_at, 'relative')}</time>
							</td>
							<td class="py-4 pr-4 pl-3 text-center text-sm whitespace-nowrap sm:pr-0">
								<a href={`/admin/pages/${page.slug}`} class="text-indigo-600 hover:text-indigo-900"
									>Edit<span class="sr-only">Edit Action</span></a
								>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		{/if}
	</section>
{:catch error}
	<p>Something went wrong {error}</p>
{/await}

<Drawers
	title="New Page"
	description="Fill in the details to create a new project."
	open={drawerOpen}
	onclose={closeDrawer}
	onsubmit={handleCreate}
>
	<div>
		<label class="mb-1 block text-sm font-medium text-gray-700" for="title">Title</label>
		<input
			id="title"
			type="text"
			bind:value={newPage.title}
			class="w-full rounded-md border px-3 py-2 text-sm shadow-sm"
			required
		/>
	</div>
	<div>
		<label for="slug" class="mb-1 block text-sm font-medium text-gray-700">Slug</label>
		<input
			id="slug"
			type="text"
			bind:value={newPage.slug}
			class="w-full rounded-md border px-3 py-2 text-sm shadow-sm"
			required
		/>
	</div>

	<!-- Create a divider line -->
	<div class="my-4 border-t text-gray-300"></div>
</Drawers>
