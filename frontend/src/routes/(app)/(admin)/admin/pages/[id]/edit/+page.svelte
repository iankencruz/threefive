<!-- frontend/src/routes/admin/pages/[id]/edit/+page.svelte -->
<script lang="ts">
	import { goto } from '$app/navigation';
	import { PUBLIC_API_URL } from '$env/static/public';
	import BlockEditor from '$components/blocks/BlockEditor.svelte';
	import SEOFields from '$components/admin/shared/SEOField.svelte';
	import { toast } from 'svelte-sonner';
	import type { PageData } from './$types';
	import type { SEOData } from '$lib/types/seo';
	import { capitalize } from '$src/lib/utilities';

	let { data }: { data: PageData } = $props();

	let formData = $state<{
		title: string;
		slug: string;
		status: 'draft' | 'published' | 'archived';
		blocks: any[];
		seo: SEOData;
	}>({
		title: data.page.title || '',
		slug: data.page.slug || '',
		status: data.page.status || ('draft' as 'draft' | 'published' | 'archived'),
		blocks: data.page.blocks || [],
		seo: {
			meta_title: data.page.seo?.meta_title || '',
			meta_description: data.page.seo?.meta_description || '',
			og_title: data.page.seo?.og_title || '',
			og_description: data.page.seo?.og_description || '',
			robots_index: data.page.seo?.robots_index ?? true,
			robots_follow: data.page.seo?.robots_follow ?? true,
			canonical_url: data.page.seo?.canonical_url || ''
		}
	});

	let errors = $state<Record<string, string>>({});
	let loading = $state(false);
	let currentTab = $state<'content' | 'seo'>('content');
	let slugManuallyEdited = $state(false);

	// Auto-generate slug from title
	$effect(() => {
		if (formData.title && !slugManuallyEdited) {
			formData.slug = formData.title
				.toLowerCase()
				.replace(/[^a-z0-9]+/g, '-')
				.replace(/^-|-$/g, '');
		}
	});

	// Auto-fill SEO meta title from page title (only if empty)
	$effect(() => {
		if (formData.title && !formData.seo.meta_title) {
			formData.seo.meta_title = formData.title;
		}
	});

	const handleSubmit = async () => {
		loading = true;
		errors = {};

		try {
			const payload = {
				title: formData.title,
				slug: formData.slug,
				page_type: 'generic',
				status: formData.status,
				blocks: formData.blocks,
				seo: formData.seo.meta_title || formData.seo.meta_description ? formData.seo : undefined
			};

			const response = await fetch(`${PUBLIC_API_URL}/api/v1/admin/pages/${data.page.id}`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json'
				},
				credentials: 'include',
				body: JSON.stringify(payload)
			});

			const result = await response.json();

			if (!response.ok) {
				if (result.errors && Array.isArray(result.errors)) {
					const newErrors: Record<string, string> = {};

					result.errors.forEach((err: { field: string; message: string }) => {
						newErrors[err.field] = err.message;
						toast.error(`${capitalize(err.field)}: ${err.message}`);
					});

					errors = newErrors;
				} else {
					toast.error(result.message || 'Failed to update page');
				}
				return;
			}

			toast.success('Page updated successfully!');
			goto(`/admin/pages/${data.page.id}/edit`);
		} catch (error) {
			console.error('Error updating page:', error);
			toast.error('An unexpected error occurred');
		} finally {
			loading = false;
		}
	};

	const handleDelete = async () => {
		if (!confirm('Are you sure you want to delete this page?')) return;

		try {
			const response = await fetch(`${PUBLIC_API_URL}/api/v1/admin/pages/${data.page.id}`, {
				method: 'DELETE',
				credentials: 'include'
			});

			if (!response.ok) {
				throw new Error('Failed to delete page');
			}

			toast.success('Page successfully deleted');
			goto('/admin/pages');
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete page');
		}
	};
</script>

<div class="  mx-auto max-w-7xl">
	<div class="relative mb-8 flex items-center justify-between gap-4">
		<div class="flex items-center gap-4">
			<button
				onclick={() => goto('/admin/pages')}
				class="rounded-lg p-2 transition-colors hover:bg-gray-700"
				aria-label="Back to Pages List"
			>
				<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M10 19l-7-7m0 0l7-7m-7 7h18"
					/>
				</svg>
			</button>
			<h1 class="">Edit Page</h1>
		</div>
	</div>

	<form
		onsubmit={(e) => {
			e.preventDefault();
			handleSubmit();
		}}
		class=" space-y-6"
	>
		<!-- Main Content Card -->
		<div class=" rounded-lg bg-surface shadow-lg">
			<!-- Tabs Navigation -->
			<div class="border-b border-gray-700">
				<nav class="flex px-6" aria-label="Tabs">
					<button
						type="button"
						onclick={() => (currentTab = 'content')}
						class="border-b-2 px-1 py-4 text-sm font-medium transition-colors {currentTab ===
						'content'
							? 'border-primary text-primary'
							: 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-200'}"
					>
						Content
					</button>
					<button
						type="button"
						onclick={() => (currentTab = 'seo')}
						class="ml-8 border-b-2 px-1 py-4 text-sm font-medium transition-colors {currentTab ===
						'seo'
							? 'border-primary text-primary'
							: 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-200'}"
					>
						SEO
					</button>
				</nav>
			</div>

			<div class="p-6">
				{#if currentTab === 'content'}
					<!-- Basic Info -->
					<div class=" mb-8 space-y-6">
						<div class="grid grid-cols-2 gap-6">
							<div>
								<label for="title" class="mb-2 block font-medium">
									Title <span class="text-red-500">*</span>
								</label>
								<input
									name="title"
									type="text"
									bind:value={formData.title}
									required
									class="form-input"
									placeholder="Enter page title"
								/>
								{#if errors.title}
									<p class="mt-1 text-sm text-red-600">{errors.title}</p>
								{/if}
							</div>

							<div>
								<label for="slug" class="mb-2 block text-sm font-medium">
									Slug <span class="text-red-500">*</span>
								</label>
								<input
									name="slug"
									type="text"
									bind:value={formData.slug}
									oninput={() => (slugManuallyEdited = true)}
									required
									class="form-input"
									placeholder="page-slug"
								/>
								{#if errors.slug}
									<p class="mt-1 text-sm text-red-600">{errors.slug}</p>
								{/if}
							</div>
						</div>

						<div>
							<label for="status" class="mb-2 block text-sm font-medium">
								Status <span class="text-red-500">*</span>
							</label>
							<select name="status" bind:value={formData.status} class="form-input">
								<option value="draft">Draft</option>
								<option value="published">Published</option>
								<option value="archived">Archived</option>
							</select>
						</div>
					</div>

					<!-- Blocks Section -->
					<div class="border-t border-gray-700 pt-8">
						<BlockEditor bind:blocks={formData.blocks} />
					</div>
				{:else if currentTab === 'seo'}
					<!-- SEO Fields -->
					<SEOFields bind:seo={formData.seo} onchange={(updated) => (formData.seo = updated)} />
				{/if}
			</div>
		</div>

		<!-- Footer Actions -->
		<div
			class="fixed right-0 bottom-0 flex w-full justify-end gap-4 border-t border-gray-400 bg-background p-3"
		>
			<div class="mr-4 flex gap-4">
				<button
					type="button"
					onclick={handleDelete}
					class="btn border-gray-400 underline-offset-6 hover:underline">Delete</button
				>

				<button
					type="submit"
					disabled={loading}
					class="rounded-lg bg-primary px-6 py-2 text-white transition-colors hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-50"
				>
					{loading ? 'Updating...' : 'Save'}
				</button>
			</div>
		</div>
	</form>
</div>
