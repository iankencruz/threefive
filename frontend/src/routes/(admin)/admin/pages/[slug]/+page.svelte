<script lang="ts">
	import { page } from '$app/state';
	import { toast } from 'svelte-sonner';
	import PageForm from '$lib/components/Forms/PageForm.svelte';
	import type { Page, Block } from '$lib/types';
	import { goto } from '$app/navigation';
	import { pages } from '$lib/store/pages.svelte';
	import { sortBlocks } from '$lib/api/pages';

	let loadedPage = $state<Page | null>(null);
	let blocks = $state<Block[]>([]);
	let slug = $state<string | null>(null);

	// Load page when slug changes
	$effect(() => {
		const incomingSlug = page.params.slug;
		slug = incomingSlug;
		loadPage(incomingSlug);
	});

	async function loadPage(slug: string): Promise<void> {
		try {
			const res = await fetch(`/api/v1/admin/pages/${slug}`);
			const json = await res.json();

			if (!res.ok) throw new Error(json.message || 'Unknown error');

			loadedPage = json.data.page;
			blocks = json.data.blocks ?? [];
		} catch (err) {
			console.error('Failed to load page:', err);
			toast.error('Failed to load page');
		}
	}

	function handleDelete(data: Page): void {
		if (confirm('Are you sure you want to delete this page? This action cannot be undone.')) {
			console.log('Deleting page:', data.slug);
			// TODO: implement deletion handler
		}
	}

	async function handleUpdate({ page, blocks }: { page: Page; blocks: Block[] }) {
		const serializableBlocks = blocks.map((b) => ({
			...b,
			props: JSON.parse(JSON.stringify(b.props))
		}));

		const sortedBlocks = serializableBlocks
			.sort((a, b) => a.sort_order - b.sort_order)
			.map((b, i) => ({ ...b, sort_order: i }));

		try {
			// ✅ 2. Then save the full page + content
			const payload = {
				page,
				blocks: sortedBlocks
			};

			console.log('🧾 Payload:', payload);

			const res = await fetch(`/api/v1/admin/pages/${page.slug}`, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(payload)
			});

			const json = await res.json();
			if (!res.ok) {
				toast.error(json.message || 'Failed to save');
				return;
			}

			toast.success('Page updated');
			pages.updatePage(page);
			goto(`/admin/pages/${page.slug}`);
		} catch (err) {
			console.error(err);
			toast.error('Update failed');
		}
	}
</script>

{#if loadedPage}
	<PageForm content={loadedPage} {blocks} onsubmit={handleUpdate} ondelete={handleDelete} />
{:else}
	<p class="text-gray-500">Loading page...</p>
{/if}
