<!-- frontend/src/routes/(app)/(admin)/admin/projects/new/+page.svelte -->
<script lang="ts">
	import { goto } from '$app/navigation';
	import { PUBLIC_API_URL } from '$env/static/public';
	import SEOFields from '$components/admin/shared/SEOField.svelte';
	import ProjectMediaGallery from '$components/projects/ProjectMediaGallery.svelte';
	import { toast } from 'svelte-sonner';
	import type { SEOData } from '$types/seo';
	import type { Media } from '$api/media';
	import { capitalize } from '$src/lib/utilities';

	let formData = $state<{
		title: string;
		slug: string;
		description: string;
		status: 'draft' | 'published' | 'archived';
		client_name: string;
		project_year: number;
		project_url: string;
		technologies: string[];
		project_status: 'completed' | 'ongoing' | 'archived';
		media_ids: string[];
		featured_image_id: string | null;
		seo: SEOData;
	}>({
		title: '',
		slug: '',
		description: '',
		status: 'draft',
		client_name: '',
		project_year: new Date().getFullYear(),
		project_url: '',
		technologies: [],
		project_status: 'completed',
		media_ids: [],
		featured_image_id: null,
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

	let projectMedia = $state<Media[]>([]);
	let errors = $state<Record<string, string>>({});
	let loading = $state(false);
	let currentTab = $state<'content' | 'seo' | 'project'>('content');
	let newTech = $state('');
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
				description: formData.description || null,
				status: formData.status,
				client_name: formData.client_name || null,
				project_year: formData.project_year,
				project_url: formData.project_url || null,
				technologies: formData.technologies,
				project_status: formData.project_status,
				media_ids: formData.media_ids,
				featured_image_id: formData.featured_image_id,
				seo: formData.seo.meta_title || formData.seo.meta_description ? formData.seo : undefined
			};

			const response = await fetch(`${PUBLIC_API_URL}/api/v1/admin/projects`, {
				method: 'POST',
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
					toast.error(result.message || 'Failed to create project');
				}
				return;
			}

			toast.success('Project created successfully!');
			goto(`/admin/projects/${result.id}/edit`);
		} catch (error) {
			console.error('Submit error:', error);
			toast.error('An unexpected error occurred');
		} finally {
			loading = false;
		}
	};
</script>

<svelte:head>
	<title>Admin: Create New Project</title>
</svelte:head>

<div class="mx-auto max-w-5xl">
	<div class="mb-8">
		<h1 class="text-3xl font-bold">Create New Project</h1>
		<p class="mt-2 text-gray-600">Add a new project to your portfolio</p>
	</div>

	<form
		onsubmit={(e) => {
			e.preventDefault();
			handleSubmit();
		}}
		class="space-y-8"
	>
		<div class="rounded-lg bg-surface shadow-sm">
			<!-- Tabs -->
			<div class="border-b border-gray-200">
				<nav class="flex space-x-8 px-6" aria-label="Tabs">
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
						class="border-b-2 px-1 py-4 text-sm font-medium transition-colors {currentTab === 'seo'
							? 'border-primary text-primary'
							: 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-200'}"
					>
						SEO
					</button>
					<button
						type="button"
						onclick={() => (currentTab = 'project')}
						class="border-b-2 px-1 py-4 text-sm font-medium transition-colors {currentTab ===
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
					<div class="space-y-6">
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
							<label class="mb-2 block text-sm font-medium">Description</label>
							<textarea
								bind:value={formData.description}
								rows="4"
								class="form-input"
								placeholder="Brief description of the project"
							></textarea>
							{#if errors.description}
								<p class="mt-1 text-sm text-red-600">{errors.description}</p>
							{/if}
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

						<!-- Project Media Gallery -->
						<div class="border-t border-gray-200 pt-8">
							<ProjectMediaGallery
								bind:media={projectMedia}
								bind:featuredImageId={formData.featured_image_id}
								onMediaChange={(mediaIds) => (formData.media_ids = mediaIds)}
								onFeaturedImageChange={(mediaId) => (formData.featured_image_id = mediaId)}
							/>
						</div>
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
										class="inline-flex items-center gap-1 rounded-full bg-primary/10 px-3 py-1 text-sm text-primary"
									>
										{tech}
										<button
											onclick={() => removeTechnology(tech)}
											type="button"
											class="hover:text-primary/70"
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
				class="rounded-lg border border-gray-300 px-6 py-2 transition-colors hover:bg-gray-50"
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
