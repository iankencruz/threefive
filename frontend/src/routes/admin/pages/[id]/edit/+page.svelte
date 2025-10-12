<!-- frontend/src/routes/admin/pages/[id]/edit/+page.svelte -->
<script lang="ts">
import { goto } from "$app/navigation";
import { page } from "$app/stores";
import BlockEditor from "$lib/components/blocks/BlockEditor.svelte";
import { PUBLIC_API_URL } from "$env/static/public";
import { toast } from "svelte-sonner";

type PageType = "generic" | "project" | "blog";
type PageStatus = "draft" | "published";

interface PageData {
	page: {
		id: string;
		title: string;
		slug: string;
		page_type: PageType;
		status: PageStatus;
		blocks: any[];
		seo?: {
			meta_title?: string;
			meta_description?: string;
			og_title?: string;
			og_description?: string;
			robots_index?: boolean;
			robots_follow?: boolean;
		};
		project_data?: {
			client_name?: string;
			project_year?: number;
			project_url?: string;
			technologies?: string[];
			project_status?: string;
		};
		blog_data?: {
			excerpt?: string;
			reading_time?: number;
		};
	};
}

let { data } = $props<{ data: PageData }>();

let formData = $state({
	title: data.page.title,
	slug: data.page.slug,
	page_type: data.page.page_type,
	status: data.page.status,
	blocks: data.page.blocks || [],
	seo: {
		meta_title: data.page.seo?.meta_title || "",
		meta_description: data.page.seo?.meta_description || "",
		og_title: data.page.seo?.og_title || "",
		og_description: data.page.seo?.og_description || "",
		robots_index: data.page.seo?.robots_index ?? true,
		robots_follow: data.page.seo?.robots_follow ?? true,
	},
	project_data: {
		client_name: data.page.project_data?.client_name || "",
		project_year:
			data.page.project_data?.project_year || new Date().getFullYear(),
		project_url: data.page.project_data?.project_url || "",
		technologies: data.page.project_data?.technologies || [],
		project_status: data.page.project_data?.project_status || "completed",
	},
	blog_data: {
		excerpt: data.page.blog_data?.excerpt || "",
		reading_time: data.page.blog_data?.reading_time || 0,
	},
});

let errors = $state<Record<string, string>>({});
let loading = $state(false);
let currentTab = $state<"content" | "seo" | "metadata">("content");
let newTech = $state("");

// Add these tracking flags for auto-generation
let slugManuallyEdited = $state(false);
let seoTitleManuallyEdited = $state(false);

// Auto-generate slug from title (only if not manually edited)
$effect(() => {
	if (formData.title && !slugManuallyEdited) {
		formData.slug = formData.title
			.toLowerCase()
			.replace(/[^a-z0-9]+/g, "-")
			.replace(/^-|-$/g, "");
	}
});

// Auto-fill SEO meta title from page title (only if not manually edited)
$effect(() => {
	if (formData.title && !seoTitleManuallyEdited) {
		formData.seo.meta_title = formData.title;
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

		const response = await fetch(
			`${PUBLIC_API_URL}/api/v1/pages/${data.page.id}`,
			{
				method: "PUT",
				credentials: "include",
				headers: {
					"Content-Type": "application/json",
				},
				body: JSON.stringify(payload),
			},
		);

		if (!response.ok) {
			const error = await response.json();
			throw new Error(error.message || "Failed to update page");
		}

		toast.success("Page Updated");
		// goto("/admin/pages");
	} catch (err) {
		if (err instanceof Error) {
			errors.submit = err.message;
		}
	} finally {
		loading = false;
	}
};

const handleDelete = async () => {
	if (!confirm("Are you sure you want to delete this page?")) return;

	try {
		const response = await fetch(
			`${PUBLIC_API_URL}/api/v1/pages/${data.page.id}`,
			{
				method: "DELETE",
				credentials: "include",
			},
		);

		if (!response.ok) {
			throw new Error("Failed to delete page");
		}

		toast.success("Page successfully deleted");
		goto("/admin/pages");
	} catch (err) {
		alert(err instanceof Error ? err.message : "Failed to delete page");
	}
};
</script>

<div class="min-h-screen bg-gray-50 py-8 px-4">
  <form onsubmit={handleSubmit} class="max-w-5xl mx-auto">
    <div class="bg-white rounded-lg shadow">
      <div class="border-b border-gray-200 px-8 py-6">
        <div class="flex items-center justify-between">
          <div>
            <button
              type="button"
              class="text-gray-600 hover:text-gray-900 mb-4 flex items-center gap-2 transition-colors"
              onclick={() => goto("/admin/pages")}
            >
              <svg
                class="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M15 19l-7-7 7-7"
                />
              </svg>
              Back to Pages
            </button>
            <h1 class="text-3xl font-bold text-gray-900">Edit Page</h1>
            <p class="text-gray-600 mt-2">
              Update your page content and settings
            </p>
          </div>
          <div class="flex items-center gap-3">
            <button
              class="px-4 py-2 text-blue-600 hover:text-blue-700 hover:bg-blue-50 rounded-lg font-medium transition-colors flex items-center gap-2"
              onclick={() =>
                window.open(`/admin/pages/${data.page.id}/preview`, "_blank")}
            >
              <svg
                class="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                />
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                />
              </svg>
              Preview
            </button>
            <button
              class="px-4 py-2 text-red-600 hover:text-red-700 hover:bg-red-50 rounded-lg font-medium transition-colors"
              onclick={handleDelete}
              disabled={loading}
            >
              Delete Page
            </button>
          </div>
        </div>
        {#if errors.submit}
          <div class="mt-4 p-4 bg-red-50 border border-red-200 rounded-lg">
            <p class="text-sm text-red-800">{errors.submit}</p>
          </div>
        {/if}

        <nav class="flex gap-8 mt-6 border-b border-gray-200">
          <button
            type="button"
            class="pb-4 px-1 border-b-2 font-medium transition-colors {currentTab ===
            'content'
              ? 'border-blue-600 text-blue-600'
              : 'border-transparent text-gray-600 hover:text-gray-900'}"
            onclick={() => (currentTab = "content")}
          >
            Content
          </button>
          <button
            type="button"
            class="pb-4 px-1 border-b-2 font-medium transition-colors {currentTab ===
            'seo'
              ? 'border-blue-600 text-blue-600'
              : 'border-transparent text-gray-600 hover:text-gray-900'}"
            onclick={() => (currentTab = "seo")}
          >
            SEO
          </button>
          {#if formData.page_type !== "generic"}
            <button
              type="button"
              class="pb-4 px-1 border-b-2 font-medium transition-colors {currentTab ===
              'metadata'
                ? 'border-blue-600 text-blue-600'
                : 'border-transparent text-gray-600 hover:text-gray-900'}"
              onclick={() => (currentTab = "metadata")}
            >
              {formData.page_type === "project" ? "Project Data" : "Blog Data"}
            </button>
          {/if}
        </nav>
      </div>

      <div class="p-8">
        {#if currentTab === "content"}
          <div class="space-y-6">
            <div class="grid grid-cols-2 gap-6">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2"
                  >Title <span class="text-red-500">*</span></label
                >
                <input
                  type="text"
                  bind:value={formData.title}
                  required
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Page title"
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2"
                  >Slug <span class="text-red-500">*</span></label
                >
                <input
                  type="text"
                  bind:value={formData.slug}
                  oninput={() => (slugManuallyEdited = true)}
                  required
                  pattern="[a-z0-9\-]+"
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="page-slug"
                />
                <p class="mt-1 text-sm text-gray-500">
                  URL-friendly version (lowercase, hyphens only)
                </p>
              </div>
            </div>

            <div class="grid grid-cols-2 gap-6">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2"
                  >Page Type</label
                >
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
                <label class="block text-sm font-medium text-gray-700 mb-2"
                  >Status</label
                >
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

          <div class="border-t border-gray-200 pt-8 mt-8">
            <BlockEditor bind:blocks={formData.blocks} />
          </div>
        {:else if currentTab === "seo"}
          <div class="space-y-6">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2"
                >Meta Title</label
              >
              <input
                type="text"
                bind:value={formData.seo.meta_title}
                oninput={() => (seoTitleManuallyEdited = true)}
                maxlength="60"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Page title for search engines (60 chars max)"
              />
              <p class="mt-1 text-sm text-gray-500">
                {formData.seo.meta_title.length}/60 characters
              </p>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2"
                >Meta Description</label
              >
              <textarea
                bind:value={formData.seo.meta_description}
                maxlength="160"
                rows="3"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Brief description for search results (160 chars max)"
              ></textarea>
              <p class="mt-1 text-sm text-gray-500">
                {formData.seo.meta_description.length}/160 characters
              </p>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2"
                >OG Title</label
              >
              <input
                type="text"
                bind:value={formData.seo.og_title}
                maxlength="60"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Title when shared on social media"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2"
                >OG Description</label
              >
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
                <input
                  type="checkbox"
                  bind:checked={formData.seo.robots_index}
                  class="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                />
                <span class="text-sm text-gray-700"
                  >Allow search engines to index</span
                >
              </label>

              <label class="flex items-center gap-2">
                <input
                  type="checkbox"
                  bind:checked={formData.seo.robots_follow}
                  class="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                />
                <span class="text-sm text-gray-700"
                  >Allow search engines to follow links</span
                >
              </label>
            </div>
          </div>
        {:else if currentTab === "metadata"}
          {#if formData.page_type === "project"}
            <div class="space-y-6">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2"
                  >Client Name</label
                >
                <input
                  type="text"
                  bind:value={formData.project_data.client_name}
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Client or company name"
                />
              </div>

              <div class="grid grid-cols-2 gap-6">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2"
                    >Project Year</label
                  >
                  <input
                    type="number"
                    bind:value={formData.project_data.project_year}
                    min="1900"
                    max="2100"
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2"
                    >Project URL</label
                  >
                  <input
                    type="url"
                    bind:value={formData.project_data.project_url}
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="https://project-url.com"
                  />
                </div>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2"
                  >Technologies</label
                >
                <div class="flex gap-2 mb-2">
                  <input
                    type="text"
                    bind:value={newTech}
                    onkeydown={(e) =>
                      e.key === "Enter" &&
                      (e.preventDefault(), addTechnology())}
                    class="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="Add technology (e.g., React, Node.js)"
                  />
                  <button
                    type="button"
                    onclick={addTechnology}
                    class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
                  >
                    Add
                  </button>
                </div>
                <div class="flex flex-wrap gap-2">
                  {#each formData.project_data.technologies as tech}
                    <span
                      class="inline-flex items-center gap-2 px-3 py-1 bg-blue-100 text-blue-800 rounded-full text-sm"
                    >
                      {tech}
                      <button
                        type="button"
                        onclick={() => removeTechnology(tech)}
                        class="hover:text-blue-900"
                      >
                        <svg
                          class="w-4 h-4"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
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
                <label class="block text-sm font-medium text-gray-700 mb-2"
                  >Project Status</label
                >
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
          {:else if formData.page_type === "blog"}
            <div class="space-y-6">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2"
                  >Excerpt</label
                >
                <textarea
                  bind:value={formData.blog_data.excerpt}
                  rows="4"
                  maxlength="300"
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Brief excerpt or summary of the blog post"
                ></textarea>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2"
                  >Reading Time (minutes)</label
                >
                <input
                  type="number"
                  bind:value={formData.blog_data.reading_time}
                  min="0"
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Estimated reading time"
                />
              </div>
            </div>
          {/if}
        {/if}
      </div>
    </div>

    <div class="flex justify-end gap-4 mt-6">
      <button
        type="button"
        class="px-6 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 font-medium transition-colors"
        onclick={() => goto("/admin/pages")}
      >
        Cancel
      </button>
      <button
        type="submit"
        disabled={loading}
        class="px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
      >
        {loading ? "Saving..." : "Update Page"}
      </button>
    </div>
  </form>
</div>
