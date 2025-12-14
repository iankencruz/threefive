<!-- frontend/src/lib/components/admin/shared/SEOFields.svelte -->
<script lang="ts">
	interface SEOData {
		meta_title: string;
		meta_description: string;
		og_title: string;
		og_description: string;
		robots_index: boolean;
		robots_follow: boolean;
		canonical_url?: string;
	}

	interface Props {
		seo: SEOData;
		onchange?: (seo: SEOData) => void;
	}

	let { seo = $bindable(), onchange }: Props = $props();

	// Notify parent of changes
	$effect(() => {
		if (onchange) {
			onchange(seo);
		}
	});
</script>

<div class="space-y-6">
	<div>
		<label for="meta_title" class="mb-2 block text-sm font-medium"> Meta Title </label>
		<input
			name="meta_title"
			type="text"
			bind:value={seo.meta_title}
			class="form-input"
			placeholder="Page title for search engines"
			maxlength="60"
		/>
		<p class="mt-1 text-xs text-gray-500">
			{seo.meta_title.length}/60 characters
		</p>
	</div>

	<div>
		<label for="meta_description" class="mb-2 block text-sm font-medium"> Meta Description </label>
		<textarea
			name="meta_description"
			bind:value={seo.meta_description}
			rows="3"
			class="form-input"
			placeholder="Brief description for search engines"
			maxlength="160"
		></textarea>
		<p class="mt-1 text-xs text-gray-500">
			{seo.meta_description.length}/160 characters
		</p>
	</div>

	<div>
		<label for="open_graph_title" class="mb-2 block text-sm font-medium"> Open Graph Title </label>
		<input
			name="open_graph_title"
			type="text"
			bind:value={seo.og_title}
			class="form-input"
			placeholder="Title when shared on social media"
		/>
	</div>

	<div>
		<label for="open_graph_description" class="mb-2 block text-sm font-medium">
			Open Graph Description
		</label>
		<textarea
			name="open_graph_description"
			bind:value={seo.og_description}
			rows="3"
			class="form-input"
			placeholder="Description when shared on social media"
		></textarea>
	</div>

	<div>
		<label for="canonical_url" class="mb-2 block text-sm font-medium"> Canonical URL </label>
		<input
			name="canonical_url"
			type="url"
			bind:value={seo.canonical_url}
			class="form-input"
			placeholder="https://example.com/canonical-page (optional)"
		/>
		<p class="mt-1 text-xs text-gray-500">Specify the preferred URL for this content</p>
	</div>

	<div class="grid grid-cols-2 gap-6">
		<label class="flex cursor-pointer items-center gap-3">
			<input
				type="checkbox"
				bind:checked={seo.robots_index}
				class="h-4 w-4 rounded border-gray-600 text-primary focus:ring-primary"
			/>
			<span class="text-sm font-medium">Allow search engines to index</span>
		</label>

		<label class="flex cursor-pointer items-center gap-3">
			<input
				type="checkbox"
				bind:checked={seo.robots_follow}
				class="h-4 w-4 rounded border-gray-600 text-primary focus:ring-primary"
			/>
			<span class="text-sm font-medium">Allow search engines to follow links</span>
		</label>
	</div>

	<div class="rounded-lg border border-blue-200 bg-blue-50 p-4">
		<h4 class="mb-2 text-sm font-medium text-blue-900">SEO Tips</h4>
		<ul class="space-y-1 text-xs text-blue-800">
			<li>• Meta title should be unique and descriptive (50-60 characters)</li>
			<li>• Meta description should summarize the page content (150-160 characters)</li>
			<li>• Open Graph tags control how content appears when shared on social media</li>
			<li>• Use canonical URLs to prevent duplicate content issues</li>
		</ul>
	</div>
</div>
