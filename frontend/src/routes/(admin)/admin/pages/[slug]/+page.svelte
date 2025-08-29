<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';
	import PageForm from '$lib/components/Forms/PageForm.svelte';
	import type { Page } from '$lib/types';
	import { updatePage } from '$lib/api/pages';

	let data = $state<Page | null>(null);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let saving = $state(false);

	async function fetchPage(slug: string, signal: AbortSignal) {
		loading = true;
		error = null;
		data = null;

		const res = await fetch(`/api/v1/admin/pages/${slug}`, { signal });
		if (!res.ok) throw new Error(`Failed to load page (${res.status})`);
		const json = await res.json();
		data = json.data as Page;
		loading = false;
	}

	// React only to slug changes
	$effect(() => {
		const slug = page.params.slug;
		if (!slug) {
			error = 'Missing slug';
			loading = false;
			return;
		}

		const ctrl = new AbortController();
		fetchPage(slug, ctrl.signal).catch((e) => {
			if (e.name !== 'AbortError') {
				console.error(e);
				error = e.message ?? 'Failed to load page';
				loading = false;
			}
		});

		// cancel in-flight request if slug changes / component reinitialises
		return () => ctrl.abort();
	});

	async function handleUpdate(next: Page): Promise<void> {
		try {
			saving = true;
			await updatePage(next, page.params.slug); // IMPORTANT: await this
			toast.success('Page updated');
			await goto('/admin/pages');
		} catch (e) {
			console.error(e);
			toast.error('Update failed');
		} finally {
			saving = false;
		}
	}

	function handleDelete(next: Page): void {
		if (confirm('Are you sure you want to delete this page? This action cannot be undone.')) {
			console.log('Deleting page:', next.slug);
			// TODO: call delete endpoint then navigate
		}
	}
</script>

{#if loading}
	<p>Fetching…</p>
{:else if error}
	<p class="text-red-600">Something went wrong: {error}</p>
{:else if data}
	<!-- Consider wiring `saving` into PageForm to disable submit buttons while saving -->
	<div class="max-w-4xl">
		<PageForm content={data} onsubmit={handleUpdate} ondelete={handleDelete} />
	</div>
{/if}
