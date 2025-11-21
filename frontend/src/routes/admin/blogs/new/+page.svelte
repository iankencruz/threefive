<!-- frontend/src/routes/admin/blogs/new/+page.svelte -->
<script lang="ts">
	import { goto } from '$app/navigation';
	import { PUBLIC_API_URL } from '$env/static/public';
	import BlockEditor from '$components/blocks/BlockEditor.svelte';
	import SEOFields from '$components/admin/shared/SEOField.svelte';
	import { toast } from 'svelte-sonner';
	import type { SEOData } from '$types/seo';

	let formData = $state<{
		title: string;
		slug: string;
		status: 'draft' | 'published' | 'archived';
		excerpt: string;
		reading_time: number;
		is_featured: boolean;
		blocks: any[];
		seo: SEOData;
	}>({
		title: '',
		slug: '',
		status: 'draft' as 'draft' | 'published' | 'archived',
		excerpt: '',
		reading_time: 1,
		is_featured: false,
		blocks: [],
		seo: {
			meta_title: '',
			meta_description: '',
			og_title: '',
			og_description: '',
			robots_index: true,
			robots_follow: true,
			canonical_url: ''
		}
	});

	let errors = $state<Record<string, string>>({});
	let loading = $state(false);
	let currentTab = $state<'content' | 'seo' | 'blog'>('content');
	let slugManuallyEdited = $state(false);
	let seoTitleManuallyEdited = $state(false);

	// Auto-generate slug from title
	$effect(() => {
		if (formData.title && !slugManuallyEdited) {
			formData.slug = formData.title
				.toLowerCase()
				.replace(/[^a-z0-9]+/g, '-')
				.replace(/^-|-$/g, '');
		}
	});

	// Auto-fill SEO meta title from page title
	$effect(() => {
		if (formData.title && !seoTitleManuallyEdited) {
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
				status: formData.status,
				excerpt: formData.excerpt,
				reading_time: formData.reading_time,
				is_featured: formData.is_featured,
				blocks: formData.blocks,
				seo: formData.seo.meta_title || formData.seo.meta_description ? formData.seo : undefined
			};

			const response = await fetch(`${PUBLIC_API_URL}/api/v1/blogs`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				credentials: 'include',
				body: JSON.stringify(payload)
			});

			if (!response.ok) {
				const errorData = await response.json();
				if (errorData.errors) {
					errors = errorData.errors;
					toast.error('Please fix the validation errors');
				} else {
					toast.error(errorData.message || 'Failed to create blog');
				}
				return;
			}

			const result = await response.json();
			toast.success('Blog created successfully!');

			// Redirect to blogs list
			goto('/admin/blogs');
		} catch (error) {
			console.error('Error creating blog:', error);
			toast.error('An unexpected error occurred');
		} finally {
			loading = false;
		}
	};
</script>

<div class="mx-auto max-w-7xl">
	<div class="mb-8 flex items-center gap-4">
		<button
			onclick={() => goto('/admin/blogs')}
			class="rounded-lg p-2 transition-colors hover:bg-gray-700"
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
		<h1 class="">Create New Blog Post</h1>
	</div>

	<form
		onsubmit={(e) => {
			e.preventDefault();
			handleSubmit();
		}}
		class="space-y-6"
	>
		<!-- Main Content Card -->
		<div class="overflow-hidden rounded-lg bg-surface shadow-lg">
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
					<button
						type="button"
						onclick={() => (currentTab = 'blog')}
						class="ml-8 border-b-2 px-1 py-4 text-sm font-medium transition-colors {currentTab ===
						'blog'
							? 'border-primary text-primary'
							: 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-200'}"
					>
						Blog Details
					</button>
				</nav>
			</div>

			<div class="p-6">
				{#if currentTab === 'content'}
					<!-- Basic Info -->
					<div class="mb-8 space-y-6">
						<div class="grid grid-cols-2 gap-6">
							<div>
								<label class="mb-2 block font-medium">
									Title <span class="text-red-500">*</span>
								</label>
								<input
									type="text"
									bind:value={formData.title}
									required
									class="form-input"
									placeholder="Enter blog post title"
								/>
								{#if errors.title}
									<p class="mt-1 text-sm text-red-600">{errors.title}</p>
								{/if}
							</div>

							<div>
								<label class="mb-2 block text-sm font-medium">
									Slug <span class="text-red-500">*</span>
								</label>
								<input
									type="text"
									bind:value={formData.slug}
									oninput={() => (slugManuallyEdited = true)}
									required
									class="form-input"
									placeholder="blog-post-slug"
								/>
								{#if errors.slug}
									<p class="mt-1 text-sm text-red-600">{errors.slug}</p>
								{/if}
							</div>
						</div>

						<div>
							<label class="mb-2 block text-sm font-medium">
								Status <span class="text-red-500">*</span>
							</label>
							<select bind:value={formData.status} class="form-input">
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
				{:else if currentTab === 'blog'}
					<!-- Blog-Specific Fields -->
					<div class="space-y-6">
						<div>
							<label class="mb-2 block text-sm font-medium">Excerpt</label>
							<textarea
								bind:value={formData.excerpt}
								maxlength="500"
								rows="4"
								class="form-input"
								placeholder="Brief summary of the blog post (500 chars max)"
							></textarea>
							<p class="mt-1 text-sm text-gray-500">
								{formData.excerpt.length}/500 characters
							</p>
						</div>

						<div>
							<label class="mb-2 block text-sm font-medium">Reading Time (minutes)</label>
							<input
								type="number"
								bind:value={formData.reading_time}
								min="0"
								class="form-input"
								placeholder="Estimated reading time in minutes"
							/>
							<p class="mt-1 text-xs text-gray-500">
								Leave empty to auto-calculate based on word count
							</p>
						</div>

						<div>
							<label class="flex cursor-pointer items-center gap-3">
								<input
									type="checkbox"
									bind:checked={formData.is_featured}
									class="h-4 w-4 rounded border-gray-600 text-primary focus:ring-primary"
								/>
								<span class="text-sm font-medium">Featured Post</span>
							</label>
							<p class="mt-1 ml-7 text-xs text-gray-500">
								Featured posts appear prominently on the blog homepage
							</p>
						</div>
					</div>
				{/if}
			</div>
		</div>

		<!-- Footer Actions -->
		<div class="flex justify-end gap-4">
			<button
				type="button"
				onclick={() => goto('/admin/blogs')}
				class="rounded-lg border border-gray-600 px-6 py-2 transition-colors hover:bg-gray-700"
			>
				Cancel
			</button>
			<button
				type="submit"
				disabled={loading}
				class="rounded-lg bg-primary px-6 py-2 text-white transition-colors hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-50"
			>
				{loading ? 'Creating...' : 'Create Blog'}
			</button>
		</div>
	</form>
</div>
