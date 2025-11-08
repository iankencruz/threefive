<!-- frontend/src/routes/admin/pages/[id]/edit/+page.svelte -->
<script lang="ts">
	import { goto } from "$app/navigation";
	import BlockEditor from "$components/blocks/BlockEditor.svelte";
	import type { Block } from "$components/blocks/BlockRenderer.svelte";
	import { PUBLIC_API_URL } from "$env/static/public";
	import { ArrowLeftIcon, Eye, Trash2Icon } from "lucide-svelte";
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
			reading_time: data.page.blog_data?.reading_time || null,
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
		const title = formData.title; // Track only title
		if (title && !slugManuallyEdited) {
			formData.slug = title
				.toLowerCase()
				.replace(/[^a-z0-9]+/g, "-")
				.replace(/^-|-$/g, "");
		}
	});

	// Auto-fill SEO meta title from page title (only if not manually edited)
	$effect(() => {
		const title = formData.title; // Track only title
		if (title && !seoTitleManuallyEdited) {
			formData.seo.meta_title = title;
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
			// Transform blocks to extract media_ids from media objects
			const transformedBlocks = formData.blocks.map((block: Block) => {
				if (block.type === "gallery" && block.data?.media) {
					return {
						...block,
						data: {
							title: block.data.title,
							media_ids: block.data.media.map((m: any) => m.id),
							// Don't send the full media objects to backend
						},
					};
				}
				return block;
			});

			const payload = {
				title: formData.title,
				slug: formData.slug,
				page_type: formData.page_type,
				status: formData.status,
				blocks: transformedBlocks, // ✅ Use transformed blocks
				seo:
					formData.seo.meta_title || formData.seo.meta_description
						? formData.seo
						: undefined,
				project_data:
					formData.page_type === "project" ? formData.project_data : undefined,
				blog_data:
					formData.page_type === "blog"
						? {
								excerpt: formData.blog_data.excerpt || undefined,
								reading_time:
									formData.blog_data.reading_time &&
									formData.blog_data.reading_time > 0
										? formData.blog_data.reading_time
										: undefined,
							}
						: undefined,
			};

			console.log("=== DEBUGGING GALLERY BLOCK SAVE ===");
			console.log("Full payload:", JSON.stringify(payload, null, 2));

			const galleryBlocks = payload.blocks?.filter(
				(b: any) => b.type === "gallery",
			);
			console.log("Gallery blocks:", galleryBlocks);

			if (galleryBlocks && galleryBlocks.length > 0) {
				galleryBlocks.forEach((block: any, idx: number) => {
					console.log(`\nGallery Block #${idx + 1}:`);
					console.log("  - data.media_ids:", block.data?.media_ids);
					console.log("  - media_ids length:", block.data?.media_ids?.length);
				});
			}
			console.log("=== END DEBUG ===\n");

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

<div class="min-h-screen bg-background py-8 px-4">
  <form onsubmit={handleSubmit} class="max-w-5xl mx-auto">
    <div class="bg-surface rounded-lg shadow">
      <div class="border-b border-gray-200 px-8 py-6">
        <div class="flex items-center justify-between">
          <div>
            <button
              type="button"
              class=" hover:text-gray-100 mb-4 flex items-center gap-2 transition-colors"
              onclick={() => goto("/admin/pages")}
            >
              <ArrowLeftIcon size={16}/>
              Back to Pages
            </button>
            <h1 class="">Edit Page</h1>
            <p class=" mt-2">
              Update your page content and settings
            </p>
          </div>
          <div class="flex items-center gap-3">
            <button
              class="px-4 py-2 text-white hover:text-primary hover:bg-blue-50 rounded-lg font-medium transition-colors flex items-center gap-2"
              onclick={() =>
                window.open(`/preview/pages/${data.page.id}`, "_blank")}
            >
              <Eye size={18} />
              Preview
            </button>
            <button
              class="px-4 py-2 text-red-600 hover:text-red-700 hover:bg-red-50 rounded-lg font-medium transition-colors flex items-center gap-2"
              onclick={handleDelete}
              disabled={loading}
            >
              <Trash2Icon size={16} />
              Delete Page
            </button>
          </div>
        </div>
        {#if errors.submit}
          <div class="mt-4 p-4 bg-red-50 border border-red-200 rounded-lg">
            <p class="text-sm text-red-800">{errors.submit}</p>
          </div>
        {/if}

        <nav class="flex gap-8 mt-6 ">
          <button
            type="button"
            class="pb-4 px-1 border-b-2 font-medium transition-colors {currentTab ===
            'content'
              ? 'border-primary text-primary'
              : 'border-transparent text-gray-700 hover:text-gray-500'}"
            onclick={() => (currentTab = "content")}
          >
            Content
          </button>
          <button
            type="button"
            class="pb-4 px-1 border-b-2 font-medium transition-colors {currentTab ===
            'seo'
              ? 'border-primary text-primary'
              : 'border-transparent text-gray-700 hover:text-gray-500'}"
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
                <label for={formData.title} class="block text-sm font-medium  mb-2"
                  >Title <span class="text-red-500">*</span></label
                >
                <input
                  type="text"
                  name={formData.title}
                  bind:value={formData.title}
                  required
                  class="form-input"
                  placeholder="Page title"
                />
              </div>

              <div>
                <label for={formData.slug} class="block text-sm font-medium  mb-2"
                  >Slug <span class="text-red-500">*</span></label
                >
                <input
                  type="text"
                  name={formData.slug}
                  bind:value={formData.slug}
                  oninput={() => (slugManuallyEdited = true)}
                  required
                  pattern="[a-z0-9\-]+"
                  class="form-input"
                  placeholder="page-slug"
                />
                <p class="mt-1 text-sm text-gray-500">
                  URL-friendly version (lowercase, hyphens only)
                </p>
              </div>
            </div>

            <div class="grid grid-cols-2 gap-6">
              <div>
                <label for={formData.page_type} class="mb-2"
                  >Page Type</label
                >
                <select
                  name={formData.page_type}
                  bind:value={formData.page_type}
                  class="form-input"
                >
                  <option value="generic">Generic</option>
                  <option value="project">Project</option>
                  <option value="blog">Blog</option>
                </select>
              </div>

              <div>
                <label for={formData.status} class="mb-2"
                  >Status</label
                >
                <select
                  name={formData.status}
                  bind:value={formData.status}
                  class="form-input"
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
              <label for={formData.seo.meta_title} class="mb-2"
                >Meta Title</label
              >
              <input
                type="text"
                name={formData.seo.meta_title}
                bind:value={formData.seo.meta_title}
                oninput={() => (seoTitleManuallyEdited = true)}
                maxlength="60"
                class="form-input"
                placeholder="Page title for search engines (60 chars max)"
              />
              <p class="mt-1 text-sm text-gray-500">
                {formData.seo.meta_title.length}/60 characters
              </p>
            </div>

            <div>
              <label for={formData.seo.meta_description} class="mb-2"
                >Meta Description</label
              >
              <textarea
                name={formData.seo.meta_description}
                bind:value={formData.seo.meta_description}
                maxlength="160"
                rows="3"
                class="form-input"
                placeholder="Brief description for search results (160 chars max)"
              ></textarea>
              <p class="mt-1 text-sm text-gray-500">
                {formData.seo.meta_description.length}/160 characters
              </p>
            </div>

            <div>
              <label for={formData.seo.og_title} class="mb-2"
                >OG Title</label
              >
              <input
                type="text"
                name={formData.seo.og_title}
                bind:value={formData.seo.og_title}
                maxlength="60"
                class="form-input"
                placeholder="Title when shared on social media"
              />
            </div>

            <div>
              <label for={formData.seo.og_description} class="mb-2"
                >OG Description</label
              >
              <textarea
                name={formData.seo.og_description}
                bind:value={formData.seo.og_description}
                maxlength="160"
                rows="3"
                class="form-input"
                placeholder="Description when shared on social media"
              ></textarea>
            </div>

            <div class="flex gap-6">
              <label for={formData.seo.robots_index} class="flex items-center gap-2">
                <input
                  type="checkbox"
                  name={formData.seo.robots_index}
                  bind:checked={formData.seo.robots_index}
                  class="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                />
                <span class=" "
                  >Allow search engines to index</span
                >
              </label>

              <label for={formData.seo.robots_follow} class="flex items-center gap-2">
                <input
                  type="checkbox"
                  name={formData.seo.robots_follow}
                  bind:checked={formData.seo.robots_follow}
                  class="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                />
                <span class=""
                  >Allow search engines to follow links</span
                >
              </label>
            </div>
          </div>
        {:else if currentTab === "metadata"}
          {#if formData.page_type === "project"}
            <div class="space-y-6">
              <div>
                <label for={formData.project_data.client_name} class="mb-2"
                  >Client Name</label
                >
                <input
                  type="text"
                  name={formData.project_data.client_name}
                  bind:value={formData.project_data.client_name}
                  class="form-input"
                  placeholder="Client or company name"
                />
              </div>

              <div class="grid grid-cols-2 gap-6">
                <div>
                  <label for={formData.project_data.project_year} class="mb-2"
                    >Project Year</label
                  >
                  <input
                    type="number"
                    name={formData.project_data.project_year}
                    bind:value={formData.project_data.project_year}
                    min="1900"
                    max="2100"
                    class="form-input"
                  />
                </div>

                <div>
                  <label for={formData.project_data.project_url} class="mb-2"
                    >Project URL</label
                  >
                  <input
                    type="url"
                    name={formData.project_data.project_url}
                    bind:value={formData.project_data.project_url}
                    class="form-input"
                    placeholder="https://project-url.com"
                  />
                </div>
              </div>

              <div>
                <label for="technologies" class="mb-2"
                  >Technologies</label
                >
                <div class="flex gap-2 mb-2">
                  <input
                    type="text"
                    name="technologies"
                    bind:value={newTech}
                    onkeydown={(e) =>
                      e.key === "Enter" &&
                      (e.preventDefault(), addTechnology())}
                    class="form-input"
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
                <label for={formData.project_data.project_status} class="mb-2"
                  >Project Status</label
                >
                <select
                  name={formData.project_data.project_status}
                  bind:value={formData.project_data.project_status}
                  class="form-input"
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
                <label for={formData.blog_data.excerpt} class="mb-2"
                  >Excerpt</label
                >
                <textarea
                  name={formData.blog_data.excerpt}
                  bind:value={formData.blog_data.excerpt}
                  rows="4"
                  maxlength="300"
                  class="form-input"
                  placeholder="Brief excerpt or summary of the blog post"
                ></textarea>
              </div>

              <div>
                <label for={formData.blog_data.reading_time} class="mb-2"
                  >Reading Time (minutes)</label
                >
                <input
                  type="number"
                  name={formData.blog_data.reading_time}
                  bind:value={formData.blog_data.reading_time}
                  min="0"
                  class="form-input"
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
        class="px-6 py-2 bg-primary hover:bg-primary/70 text-white rounded-lg font-medium disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
      >
        {loading ? "Saving..." : "Update Page"}
      </button>
    </div>
  </form>
</div>
