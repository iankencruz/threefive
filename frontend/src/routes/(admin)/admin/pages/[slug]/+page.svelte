<script lang="ts">
	import { page } from '$app/state';
	import { toast } from 'svelte-sonner';
	import PageForm from '$lib/components/Forms/PageForm.svelte';
	import type { Page } from '$lib/types';
	import { goto } from '$app/navigation';
	import { pages } from '$src/lib/store/pages.svelte';

	let loadedPage = $state<Page | null>(null);
	let slug = $state<string | null>(null);

	$effect(() => {
		loadPage(page.params.slug);
	});
	$inspect('Page URL', page.params.slug);
	//
	$inspect('Slug: ', slug);

	async function loadPage(slug: string): Promise<void> {
		try {
			const res = await fetch(`/api/v1/admin/pages/${slug}`);
			const json = await res.json();

			if (!res.ok) {
				throw new Error(json.message || 'Unknown error');
			}

			loadedPage = json.data;
		} catch (err) {
			console.error('Failed to load page:', err);
			toast.error('Failed to load page');
		}
	}

	function handleDelete(data: Page): void {
		if (confirm('Are you sure you want to delete this page? This action cannot be undone.')) {
			console.log('Deleting page:', data.slug);
			//refresh page
		}
	}

	function handleUpdate(data: Page): void {
		fetch(`/api/v1/admin/pages/${page.params.slug}`, {
			method: 'PUT',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(data)
		})
			.then(async (res) => {
				const json = await res.json();
				if (!res.ok) {
					toast.error(json.message || 'Failed to save');
					return;
				}
				toast.success('Page updated');
				pages.updatePage(data); // ✅ sync store with new title
				goto(`/admin/pages/${data.slug}`);
			})
			.catch((err) => {
				console.error(err);
				toast.error('Update failed');
			});
	}
</script>

{#if loadedPage}
	<PageForm content={loadedPage} onsubmit={handleUpdate} ondelete={handleDelete} />
{:else}
	<p class="text-gray-500">Loading page...</p>
{/if}
