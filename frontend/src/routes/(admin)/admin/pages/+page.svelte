<script lang="ts">
	import { toast } from 'svelte-sonner';
	import PageForm from '$lib/components/Forms/PageForm.svelte';
	import type { Page } from '$lib/types';
	import { getPages } from '$lib/api/pages';
	import { goto } from '$app/navigation';
	import EmptyState from '$lib/components/Overlays/EmptyState.svelte';
	import { auth } from '$lib/store/auth.svelte';
	import { formatDate } from '$lib/utils/utilities';
	import { onMount } from 'svelte';
	import Drawers from '$lib/components/Overlays/Drawers.svelte';

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

	let drawerOpen = $state(false);
	// let title = $state('');
	// let slug = $state('');
	// let description = $state('');
	// let projects = $state<any[]>([]);
	// let loading = $state(true);

	$inspect('newPage:', newPage);

	let pageContent = $state<Promise<Page[]> | null>(null);

	onMount(async () => {
		pageContent = getPages();
	});

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
			<h1 class="text-2xl font-semibold text-gray-900">Pages</h1>
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
			<ul class="space-y-2">
				{#each data as page}
					<li
						class="flex items-center justify-between gap-x-6 border-t border-gray-200 py-5 first:border-none"
					>
						<div class="min-w-0">
							<div class="flex items-start gap-x-3">
								<p class="text-sm/6 font-semibold text-gray-900">{page.title}</p>

								{#if page.is_published}
									<p
										class="mt-0.5 rounded-md bg-green-50 px-1.5 py-0.5 text-xs font-medium whitespace-nowrap text-green-700 ring-1 ring-green-600/20 ring-inset"
									>
										Published
									</p>
								{:else}
									<p
										class="mt-0.5 rounded-md bg-yellow-50 px-1.5 py-0.5 text-xs font-medium whitespace-nowrap text-yellow-800 ring-1 ring-yellow-600/20 ring-inset"
									>
										Draft
									</p>
								{/if}
							</div>
							<div class="mt-1 flex items-center gap-x-2 text-xs/5 text-gray-500">
								<p class="whitespace-nowrap">
									<strong>Created</strong>:
									<time datetime="2023-03-17T00:00Z">{formatDate(page.created_at, 'relative')}</time
									>
								</p>
								<svg viewBox="0 0 2 2" class="size-0.5 fill-current">
									<circle cx="1" cy="1" r="1" />
								</svg>
								<p class="truncate">
									Created by <span>{auth.user?.first_name} {auth.user?.last_name}</span>
								</p>
							</div>
						</div>
						<div class="flex flex-none items-center gap-x-2">
							<a
								href={`/admin/pages/${page.slug}`}
								class="hover hidden rounded-md bg-white px-2.5 py-1.5 text-sm font-semibold text-gray-900 shadow-xs ring-1 ring-gray-300 ring-inset hover:bg-gray-50 sm:block"
								>View page<span class="sr-only">View page</span></a
							>
							<button
								onclick={() => handleDelete(page.slug)}
								class=" inline-flex items-center rounded-md bg-red-500 px-2 py-1.5 text-sm font-semibold text-white shadow-xs ring-1 ring-gray-300 ring-inset hover:bg-red-700"
							>
								Delete
								<span class="sr-only">Delete page</span>
							</button>
						</div>
					</li>
				{/each}
			</ul>
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
