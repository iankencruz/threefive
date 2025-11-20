<!-- frontend/src/routes/admin/pages/new/+page.svelte -->
<script lang="ts">
	import { goto } from "$app/navigation";
	import { PUBLIC_API_URL } from "$env/static/public";
	import BlockEditor from "$components/blocks/BlockEditor.svelte";
	import SEOFields from "$components/admin/shared/SEOField.svelte";
	import { toast } from "svelte-sonner";

	let formData = $state({
		title: "",
		slug: "",
		status: "draft" as "draft" | "published" | "archived",
		blocks: [],
		seo: {
			meta_title: "",
			meta_description: "",
			og_title: "",
			og_description: "",
			robots_index: true,
			robots_follow: true,
			canonical_url: "",
		},
	});

	let errors = $state<Record<string, string>>({});
	let loading = $state(false);
	let currentTab = $state<"content" | "seo">("content");
	let slugManuallyEdited = $state(false);
	let seoTitleManuallyEdited = $state(false);

	// Auto-generate slug from title
	$effect(() => {
		if (formData.title && !slugManuallyEdited) {
			formData.slug = formData.title
				.toLowerCase()
				.replace(/[^a-z0-9]+/g, "-")
				.replace(/^-|-$/g, "");
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
				page_type: "generic",
				status: formData.status,
				blocks: formData.blocks,
				seo:
					formData.seo.meta_title || formData.seo.meta_description
						? formData.seo
						: undefined,
			};

			const response = await fetch(`${PUBLIC_API_URL}/api/v1/pages`, {
				method: "POST",
				headers: {
					"Content-Type": "application/json",
				},
				credentials: "include",
				body: JSON.stringify(payload),
			});

			if (!response.ok) {
				const errorData = await response.json();
				if (errorData.errors) {
					errors = errorData.errors;
					toast.error("Please fix the validation errors");
				} else {
					toast.error(errorData.message || "Failed to create page");
				}
				return;
			}

			const result = await response.json();
			toast.success("Page created successfully!");

			// Redirect to pages list
			goto("/admin/pages");
		} catch (error) {
			console.error("Error creating page:", error);
			toast.error("An unexpected error occurred");
		} finally {
			loading = false;
		}
	};
</script>

<div class="max-w-7xl mx-auto">
	<div class="flex items-center gap-4 mb-8">
		<button
			onclick={() => goto('/admin/pages')}
			class="p-2 hover:bg-gray-700 rounded-lg transition-colors"
		>
			<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M10 19l-7-7m0 0l7-7m-7 7h18"
				/>
			</svg>
		</button>
		<h1 class="">Create New Page</h1>
	</div>

	<form
		onsubmit={(e) => {
			e.preventDefault();
			handleSubmit();
		}}
		class="space-y-6"
	>
		<!-- Main Content Card -->
		<div class="bg-surface rounded-lg shadow-lg overflow-hidden">
			<!-- Tabs Navigation -->
			<div class="border-b border-gray-700">
				<nav class="flex px-6" aria-label="Tabs">
					<button
						type="button"
						onclick={() => (currentTab = "content")}
						class="py-4 px-1 border-b-2 font-medium text-sm transition-colors {currentTab ===
						'content'
							? 'border-primary text-primary'
							: 'border-transparent text-gray-500 hover:text-gray-200 hover:border-gray-300'}"
					>
						Content
					</button>
					<button
						type="button"
						onclick={() => (currentTab = "seo")}
						class="ml-8 py-4 px-1 border-b-2 font-medium text-sm transition-colors {currentTab ===
						'seo'
							? 'border-primary text-primary'
							: 'border-transparent text-gray-500 hover:text-gray-200 hover:border-gray-300'}"
					>
						SEO
					</button>
				</nav>
			</div>

			<div class="p-6">
				{#if currentTab === "content"}
					<!-- Basic Info -->
					<div class="space-y-6 mb-8">
						<div class="grid grid-cols-2 gap-6">
							<div>
								<label class="block font-medium mb-2">
									Title <span class="text-red-500">*</span>
								</label>
								<input
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
								<label class="block text-sm font-medium mb-2">
									Slug <span class="text-red-500">*</span>
								</label>
								<input
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
							<label class="block text-sm font-medium mb-2">
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
				{:else if currentTab === "seo"}
					<!-- SEO Fields -->
					<SEOFields
						bind:seo={formData.seo}
						onchange={(updated) => (formData.seo = updated)}
					/>
				{/if}
			</div>
		</div>

		<!-- Footer Actions -->
		<div class="flex justify-end gap-4">
			<button
				type="button"
				onclick={() => goto('/admin/pages')}
				class="px-6 py-2 border border-gray-600 rounded-lg hover:bg-gray-700 transition-colors"
			>
				Cancel
			</button>
			<button
				type="submit"
				disabled={loading}
				class="px-6 py-2 bg-primary text-white rounded-lg hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
			>
				{loading ? "Creating..." : "Create Page"}
			</button>
		</div>
	</form>
</div>
