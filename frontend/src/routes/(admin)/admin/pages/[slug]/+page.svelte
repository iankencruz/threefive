<script lang="ts">
	import { page } from '$app/state';
	import { toast } from 'svelte-sonner';
	import PageForm from '$lib/components/Forms/PageForm.svelte';
	import type { Page } from '$lib/types';
	import { goto } from '$app/navigation';
	import { updatePage } from '$lib/api/pages';

	async function loadPage(slug: string): Promise<Page> {
		const res = await fetch(`/api/v1/admin/pages/${slug}`);
		const data = await res.json();

		return data.data;
	}

	function handleDelete(data: Page): void {
		if (confirm('Are you sure you want to delete this page? This action cannot be undone.')) {
			console.log('Deleting page:', data.slug);
			// TODO: implement deletion handler
		}
	}

	console.log('slug:', page.params.slug);
	async function handleUpdate(data: Page): Promise<void> {
		try {
			const res = updatePage(data, page.params.slug);

			toast.success('Page updated');
			goto('/admin/pages');
		} catch (err) {
			console.error(err);
			toast.error('Update failed');
		}
	}
</script>

{#await loadPage(page.params.slug ?? '')}
	<p>fetching</p>
{:then data}
	<PageForm content={data} onsubmit={handleUpdate} ondelete={handleDelete} />
{:catch error}
	<p>Something went wrong: {error}</p>
{/await}
