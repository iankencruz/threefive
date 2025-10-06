<!-- frontend/src/routes/admin/pages/new/+page.svelte -->
<script lang="ts">
import { goto } from "$app/navigation";
import BlockEditor from "$lib/components/blocks/BlockEditor.svelte";

type PageType = "generic" | "project" | "blog";
type PageStatus = "draft" | "published";

let formData = $state({
	title: "",
	slug: "",
	page_type: "generic" as PageType,
	status: "draft" as PageStatus,
	blocks: [],
	seo: {
		meta_title: "",
		meta_description: "",
		og_title: "",
		og_description: "",
		robots_index: true,
		robots_follow: true,
	},
	project_data: {
		client_name: "",
		project_year: new Date().getFullYear(),
		project_url: "",
		technologies: [] as string[],
		project_status: "completed",
	},
	blog_data: {
		excerpt: "",
		reading_time: 0,
	},
});

let errors = $state<Record<string, string>>({});
let loading = $state(false);
let currentTab = $state<"content" | "seo" | "metadata">("content");
let newTech = $state("");

// Auto-generate slug from title
$effect(() => {
	if (formData.title && !formData.slug) {
		formData.slug = formData.title
			.toLowerCase()
			.replace(/[^a-z0-9]+/g, "-")
			.replace(/^-|-$/g, "");
	}
});

const addTechnology = () => {
	if (
		newTech.trim() &&
		!formData.project_data.technologies.includes(newTech.trim())
	) {
		formData.project_data.technologies = [
			...formData.project_data.technologies,
			newTech.trim(),
		];
		newTech = "";
	}
};

const removeTechnology = (tech: string) => {
	formData.project_data.technologies =
		formData.project_data.technologies.filter((t) => t !== tech);
};

const handleSubmit = async () => {
	loading = true;
	errors = {};

	try {
		const payload = {
			title: formData.title,
			slug: formData.slug,
			page_type: formData.page_type,
			status: formData.status,
			blocks: formData.blocks,
			seo:
				formData.seo.meta_title || formData.seo.meta_description
					? formData.seo
					: undefined,
			project_data:
				formData.page_type === "project" ? formData.project_data : undefined,
			blog_data: formData.page_type === "blog" ? formData.blog_data : undefined,
		};

		const response = await fetch("/api/v1/pages", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(payload),
		});

		if (!response.ok) {
			const error = await response.json();
			if (error.errors) {
				error.errors.forEach((e: any) => {
					errors[e.field] = e.message;
				});
			} else {
				errors.general = error.message || "Failed to create page";
			}
			return;
		}

		goto("/admin/pages");
	} catch (error) {
		errors.general = "An unexpected error occurred";
	} finally {
		loading = false;
	}
};
</script>

<div class="max-w-5xl mx-auto">
	<!-- Header -->
	<div class="mb-8">
		<div class="flex items-center gap-4 mb-4">
			<button
				onclick={() => goto('/admin/pages')}
				class="p-2 hover:bg-gray-100 rounded-lg transition-colors"
				aria-label="Go back"
			>
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
				</svg>
			</button>
			<h1 class="text-3xl font-bold text-gray-900">Create New Page</h1>
		</div>
	</div>

	{#if errors.general}
		<div class="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">
			{errors.general}
		</div>
	{/if}

	<form onsubmit={(e) => (e.preventDefault(), handleSubmit())}>
		<div class="bg-white rounded-lg shadow">
			<!-- Tabs -->
			<div class="border-b border-gray-200">
				<nav class="flex gap-8 px-6" aria-label="Tabs">
					<button
						type="button"
						onclick={() => currentTab = 'content'}
						class="py-4 px-1 border-b-2 font-medium text-sm transition-colors {currentTab === 'content' ? 'border-blue-500 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
					>
						Content
					</button>
					<button
						type="button"
						onclick={() => currentTab = 'seo'}
						class="py-4 px-1 border-b-2 font-medium text-sm transition-colors {currentTab === 'seo' ? 'border-blue-500 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
					>
						SEO
					</button>
					<button
						type="button"
						onclick={() => currentTab = 'metadata'}
						class="py-4 px-1 border-b-2 font-medium text-sm transition-colors {currentTab === 'metadata' ? 'border-blue-500 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
					>
						{formData.page_type === 'project' ? 'Project Data' : formData.page_type === 'blog' ? 'Blog Data' : 'Metadata'}
					</button>
				</nav>
			</div>

			<div class="p-6">
				{#if currentTab === 'content'}
					<!-- Basic Info -->
					<div class="space-y-6 mb-8">
						<div class="grid grid-cols-2 gap-6">
							<div>
								<label class="block text-sm font-medium text-gray-700 mb-2">
									Title <span class="text-red-500">*</span>
								</label>
								<input
									type="text"
									bind:value={formData.title}
									required
									class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
									placeholder="Enter page title"
								/>
								{#if errors.title}
									<p class="mt-1 text-sm text-red-600">{errors.title}</p>
								{/if}
							</div>

							<div>
								<label class="block text-sm font-medium text-gray-700 mb-2">
									Slug <span class="text-red-500">*</span>
								</label>
								<input
									type="text"
									bind:value={formData.slug}
									required
									class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
									placeholder="page-slug"
								/>
								{#if errors.slug}
									<p class="mt-1 text-sm text-red-600">{errors.slug}</p>
								{/if}
							</div>
						</div>

						<div class="grid grid-cols-2 gap-6">
							<div>
								<label class="block text-sm font-medium text-gray-700 mb-2">
									Page Type <span class="text-red-500">*</span>
								</label>
								<select
									bind:value={formData.page_type}
									class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
								>
									<option value="generic">Generic</option>
									<option value="project">Project</option>
									<option value="blog">Blog</option>
								</select>
							</div>

							<div>
								<label class="block text-sm font-medium text-gray-700 mb-2">
									Status <span class="text-red-500">*</span>
								</label>
								<select
									bind:value={formData.status}
									class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
								>
									<option value="draft">Draft</option>
									<option value="published">Published</option>
								</select>
							</div>
						</div>
					</div>

					<!-- Blocks Section using BlockEditor Component -->
					<div class="border-t border-gray-200 pt-8">
						<BlockEditor bind:blocks={formData.blocks} />
					</div>

				{:else if currentTab === 'seo'}
					<div class="space-y-6">
						<div>
							<label class="block text-sm font-medium text-gray-700 mb-2">Meta Title</label>
							<input
								type="text"
								bind:value={formData.seo.meta_title}
								maxlength="60"
								class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
								placeholder="Page title for search engines (60 chars max)"
							/>
							<p class="mt-1 text-sm text-gray-500">{formData.seo.meta_title.length}/60 characters</p>
						</div>

						<div>
							<label class="block text-sm font-medium text-gray-700 mb-2">Meta Description</label>
							<textarea
								bind:value={formData.seo.meta_description}
								maxlength="160"
								rows="3"
								class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
								placeholder="Brief description for search results (160 chars max)"
							></textarea>
							<p class="mt-1 text-sm text-gray-500">{formData.seo.meta_description.length}/160 characters</p>
						</div>

						<div>
							<label class="block text-sm font-medium text-gray-700 mb-2">OG Title</label>
							<input
								type="text"
								bind:value={formData.seo.og_title}
								maxlength="60"
								class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
								placeholder="Title when shared on social media"
							/>
						</div>

						<div>
							<label class="block text-sm font-medium text-gray-700 mb-2">OG Description</label>
							<textarea
								bind:value={formData.seo.og_description}
								maxlength="160"
								rows="3"
								class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
								placeholder="Description when shared on social media"
							></textarea>
						</div>

						<div class="flex gap-6">
							<label class="flex items-center gap-2">
								<input type="checkbox" bind:checked={formData.seo.robots_index} class="w-4 h-4 text-blue-600 rounded" />
								<span class="text-sm text-gray-700">Allow search engines to index</span>
							</label>
							<label class="flex items-center gap-2">
								<input type="checkbox" bind:checked={formData.seo.robots_follow} class="w-4 h-4 text-blue-600 rounded" />
								<span class="text-sm text-gray-700">Allow search engines to follow links</span>
							</label>
						</div>
					</div>

				{:else if currentTab === 'metadata'}
					{#if formData.page_type === 'project'}
						<div class="space-y-6">
							<div class="grid grid-cols-2 gap-6">
								<div>
									<label class="block text-sm font-medium text-gray-700 mb-2">Client Name</label>
									<input
										type="text"
										bind:value={formData.project_data.client_name}
										class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
										placeholder="Client or company name"
									/>
								</div>

								<div>
									<label class="block text-sm font-medium text-gray-700 mb-2">Project Year</label>
									<input
										type="number"
										bind:value={formData.project_data.project_year}
										min="1900"
										max="2100"
										class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
									/>
								</div>
							</div>

							<div>
								<label class="block text-sm font-medium text-gray-700 mb-2">Project URL</label>
								<input
									type="url"
									bind:value={formData.project_data.project_url}
									class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
									placeholder="https://example.com"
								/>
							</div>

							<div>
								<label class="block text-sm font-medium text-gray-700 mb-2">Technologies</label>
								<div class="flex gap-2 mb-3">
									<input
										type="text"
										bind:value={newTech}
										onkeydown={(e) => e.key === 'Enter' && (e.preventDefault(), addTechnology())}
										class="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
										placeholder="Add a technology"
									/>
									<button
										onclick={addTechnology}
										type="button"
										class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
									>
										Add
									</button>
								</div>
								<div class="flex flex-wrap gap-2">
									{#each formData.project_data.technologies as tech}
										<span class="inline-flex items-center gap-1 px-3 py-1 bg-blue-100 text-blue-800 rounded-full text-sm">
											{tech}
											<button
												onclick={() => removeTechnology(tech)}
												type="button"
												class="hover:text-blue-900"
												aria-label="Remove {tech}"
											>
												<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
													<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
												</svg>
											</button>
										</span>
									{/each}
								</div>
							</div>

							<div>
								<label class="block text-sm font-medium text-gray-700 mb-2">Project Status</label>
								<select
									bind:value={formData.project_data.project_status}
									class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
								>
									<option value="completed">Completed</option>
									<option value="ongoing">Ongoing</option>
									<option value="archived">Archived</option>
								</select>
							</div>
						</div>

					{:else if formData.page_type === 'blog'}
						<div class="space-y-6">
							<div>
								<label class="block text-sm font-medium text-gray-700 mb-2">Excerpt</label>
								<textarea
									bind:value={formData.blog_data.excerpt}
									maxlength="500"
									rows="4"
									class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
									placeholder="Brief summary of the blog post (500 chars max)"
								></textarea>
								<p class="mt-1 text-sm text-gray-500">{formData.blog_data.excerpt.length}/500 characters</p>
							</div>

							<div>
								<label class="block text-sm font-medium text-gray-700 mb-2">Reading Time (minutes)</label>
								<input
									type="number"
									bind:value={formData.blog_data.reading_time}
									min="0"
									class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
									placeholder="Estimated reading time in minutes"
								/>
							</div>
						</div>

					{:else}
						<div class="text-center py-12 text-gray-500">
							<p>No additional metadata required for generic pages.</p>
						</div>
					{/if}
				{/if}
			</div>
		</div>

		<!-- Footer Actions -->
		<div class="mt-6 flex justify-end gap-4">
			<button
				type="button"
				onclick={() => goto('/admin/pages')}
				class="px-6 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors"
			>
				Cancel
			</button>
			<button
				type="submit"
				disabled={loading}
				class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
			>
				{loading ? 'Creating...' : 'Create Page'}
			</button>
		</div>
	</form>
</div>
