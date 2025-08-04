<script lang="ts">
	import { onMount } from 'svelte';
	import { createProject, getProjects, deleteProjectBySlug } from '$lib/api/projects';
	import { toast } from 'svelte-sonner';
	import Drawers from '$lib/components/Overlays/Drawers.svelte';
	import EmptyState from '$lib/components/Overlays/EmptyState.svelte';
	import { formatDate, slugify } from '$lib/utils/utilities';
	import { FolderPlus, PaperclipIcon } from '@lucide/svelte';
	import { Paperclip } from 'lucide';
	import { auth } from '$lib/store/auth.svelte';

	let drawerOpen = $state(false);
	let title = $state('');
	let slug = $state('');
	let description = $state('');
	let projects = $state<any[]>([]);
	let loading = $state(true);

	$effect(() => {
		slug = slugify(title);
	});

	onMount(() => {
		fetchProjects();
	});

	function openCreateDrawer() {
		title = '';
		description = '';
		drawerOpen = true;
	}

	function closeDrawer() {
		drawerOpen = false;
	}

	async function fetchProjects() {
		loading = true;
		try {
			projects = await getProjects();
		} catch (err) {
			console.error('❌ Failed to load projects', err);
			toast.error('Failed to load projects');
		} finally {
			loading = false;
		}
	}

	async function deleteProject(slug: string) {
		if (!confirm('Delete project?')) return;

		try {
			await deleteProjectBySlug(slug);
			toast.success('Project deleted');
			await fetchProjects();
		} catch (error) {
			console.error('❌ Failed to delete project', error);
			toast.error('Failed to delete project. Please try again.');
		}
	}

	async function handleSubmit() {
		try {
			await createProject({ title, slug, description });
			drawerOpen = false;
			toast.success('Project created successfully!');
			await fetchProjects();
		} catch (error) {
			console.error('❌', error);
			toast.error('Failed to create project. Please try again.');
		}
	}
</script>

<section class="py-6">
	<div class="mb-6 flex items-center justify-between">
		<h1 class="text-2xl font-semibold text-gray-900">Blogs</h1>
		<button
			onclick={openCreateDrawer}
			class="rounded-md border bg-indigo-600 px-3 py-1.5 text-sm text-black text-white hover:bg-black/10"
		>
			+ New Project
		</button>
	</div>

	{#if loading}
		<div class="text-gray-500">Loading projects...</div>
	{:else if !projects}
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
			{#each projects as project}
				<li
					class="flex items-center justify-between gap-x-6 border-t border-gray-200 py-5 first:border-none"
				>
					<div class="min-w-0">
						<div class="flex items-start gap-x-3">
							<p class="text-sm/6 font-semibold text-gray-900">{project.title}</p>

							{#if project.is_published}
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
								<time datetime="2023-03-17T00:00Z"
									>{formatDate(project.created_at, 'relative')}</time
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
							href={`/admin/projects/${project.slug}`}
							class="hover hidden rounded-md bg-white px-2.5 py-1.5 text-sm font-semibold text-gray-900 shadow-xs ring-1 ring-gray-300 ring-inset hover:bg-gray-50 sm:block"
							>View project<span class="sr-only">View Project</span></a
						>
						<button
							onclick={() => deleteProject(project.id)}
							class=" inline-flex items-center rounded-md bg-red-500 px-2 py-1.5 text-sm font-semibold text-white shadow-xs ring-1 ring-gray-300 ring-inset hover:bg-red-700"
						>
							Delete
							<span class="sr-only">Delete Project</span>
						</button>
					</div>
				</li>
			{/each}
		</ul>
	{/if}
</section>

<Drawers
	title="New Project"
	description="Fill in the details to create a new project."
	open={drawerOpen}
	onclose={closeDrawer}
	onsubmit={handleSubmit}
>
	<div>
		<label class="mb-1 block text-sm font-medium text-gray-700" for="title">Title</label>
		<input
			id="title"
			type="text"
			bind:value={title}
			class="w-full rounded-md border px-3 py-2 text-sm shadow-sm"
			required
		/>
	</div>
	<div>
		<label for="slug" class="mb-1 block text-sm font-medium text-gray-700">Slug</label>
		<input
			id="slug"
			type="text"
			bind:value={slug}
			class="w-full rounded-md border px-3 py-2 text-sm shadow-sm"
			required
		/>
	</div>
	<div>
		<label for="description" class="mb-1 block text-sm font-medium text-gray-700">Description</label
		>
		<textarea
			id="description"
			rows="4"
			bind:value={description}
			class="w-full rounded-md border px-3 py-2 text-sm shadow-sm"
		></textarea>
	</div>

	<!-- Create a divider line -->
	<div class="my-4 border-t text-gray-300"></div>
</Drawers>
