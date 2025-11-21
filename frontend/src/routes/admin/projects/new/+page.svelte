<!-- frontend/src/routes/admin/projects/new/+page.svelte -->
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
		client_name: string;
		project_year: number;
		project_url: string;
		technologies: string[];
		project_status: 'completed' | 'ongoing' | 'archived';
		blocks: any[];
		seo: SEOData;
	}>({
		title: '',
		slug: '',
		status: 'draft' as 'draft' | 'published' | 'archived',
		client_name: '',
		project_year: new Date().getFullYear(),
		project_url: '',
		technologies: [] as string[],
		project_status: 'completed' as 'completed' | 'ongoing' | 'archived',
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
	let currentTab = $state<'content' | 'seo' | 'project'>('content');
	let newTech = $state('');
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

	const addTechnology = () => {
		if (newTech.trim() && !formData.technologies.includes(newTech.trim())) {
			formData.technologies = [...formData.technologies, newTech.trim()];
			newTech = '';
		}
	};

	const removeTechnology = (tech: string) => {
		formData.technologies = formData.technologies.filter((t) => t !== tech);
	};

	const handleSubmit = async () => {
		loading = true;
		errors = {};

		try {
			const payload = {
				title: formData.title,
				slug: formData.slug,
				status: formData.status,
				client_name: formData.client_name,
				project_year: formData.project_year,
				project_url: formData.project_url,
				technologies: formData.technologies,
				project_status: formData.project_status,
				blocks: formData.blocks,
				seo: formData.seo.meta_title || formData.seo.meta_description ? formData.seo : undefined
			};

			const response = await fetch(`${PUBLIC_API_URL}/api/v1/projects`, {
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
					toast.error(errorData.message || 'Failed to create project');
				}
				return;
			}

			const result = await response.json();
			toast.success('Project created successfully!');

			// Redirect to projects list
			goto('/admin/projects');
		} catch (error) {
			console.error('Error creating project:', error);
			toast.error('An unexpected error occurred');
		} finally {
			loading = false;
		}
	};
</script>

<div class="mx-auto max-w-7xl">
	<div class="mb-8 flex items-center gap-4">
		<button
			onclick={() => goto('/admin/projects')}
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
		<h1 class="">Create New Project</h1>
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
						onclick={() => (currentTab = 'project')}
						class="ml-8 border-b-2 px-1 py-4 text-sm font-medium transition-colors {currentTab ===
						'project'
							? 'border-primary text-primary'
							: 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-200'}"
					>
						Project Details
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
									placeholder="Enter project title"
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
									placeholder="project-slug"
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
				{:else if currentTab === 'project'}
					<!-- Project-Specific Fields -->
					<div class="space-y-6">
						<div class="grid grid-cols-2 gap-6">
							<div>
								<label class="mb-2 block text-sm font-medium">Client Name</label>
								<input
									type="text"
									bind:value={formData.client_name}
									class="form-input"
									placeholder="Client or company name"
								/>
							</div>

							<div>
								<label class="mb-2 block text-sm font-medium">Project Year</label>
								<input
									type="number"
									bind:value={formData.project_year}
									min="1900"
									max="2100"
									class="form-input"
								/>
							</div>
						</div>

						<div>
							<label class="mb-2 block text-sm font-medium">Project URL</label>
							<input
								type="url"
								bind:value={formData.project_url}
								class="form-input"
								placeholder="https://example.com"
							/>
						</div>

						<div>
							<label class="mb-2 block text-sm font-medium">Technologies</label>
							<div class="mb-3 flex gap-2">
								<input
									type="text"
									bind:value={newTech}
									onkeydown={(e) => e.key === 'Enter' && (e.preventDefault(), addTechnology())}
									class="form-input flex-1"
									placeholder="Add a technology"
								/>
								<button
									onclick={addTechnology}
									type="button"
									class="rounded-lg bg-primary px-4 py-2 text-white transition-colors hover:bg-primary/90"
								>
									Add
								</button>
							</div>
							<div class="flex flex-wrap gap-2">
								{#each formData.technologies as tech}
									<span
										class="inline-flex items-center gap-1 rounded-full bg-blue-900/30 px-3 py-1 text-sm text-blue-300"
									>
										{tech}
										<button
											onclick={() => removeTechnology(tech)}
											type="button"
											class="hover:text-blue-100"
											aria-label="Remove {tech}"
										>
											<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path
													stroke-linecap="round"
													stroke-linejoin="round"
													stroke-width="2"
													d="M6 18L18 6M6 6l12 12"
												/>
											</svg>
										</button>
									</span>
								{/each}
							</div>
						</div>

						<div>
							<label class="mb-2 block text-sm font-medium">Project Status</label>
							<select bind:value={formData.project_status} class="form-input">
								<option value="completed">Completed</option>
								<option value="ongoing">Ongoing</option>
								<option value="archived">Archived</option>
							</select>
						</div>
					</div>
				{/if}
			</div>
		</div>

		<!-- Footer Actions -->
		<div class="flex justify-end gap-4">
			<button
				type="button"
				onclick={() => goto('/admin/projects')}
				class="rounded-lg border border-gray-600 px-6 py-2 transition-colors hover:bg-gray-700"
			>
				Cancel
			</button>
			<button
				type="submit"
				disabled={loading}
				class="rounded-lg bg-primary px-6 py-2 text-white transition-colors hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-50"
			>
				{loading ? 'Creating...' : 'Create Project'}
			</button>
		</div>
	</form>
</div>
