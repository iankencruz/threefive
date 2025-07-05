<script lang="ts">
	import { toast } from 'svelte-sonner';
	import PageForm from '$lib/components/Forms/PageForm.svelte';
	import type { Page } from '$lib/types';

	let newPage = $state<Page>({
		id: crypto.randomUUID() as Page['id'],
		slug: '',
		title: '',
		banner_image_id: null,
		seo_title: '',
		seo_description: '',
		seo_canonical: '',
		content: [],
		is_draft: true,
		is_published: false,
		created_at: new Date().toISOString(),
		updated_at: new Date().toISOString()
	});

	function handleCreate(data: Page): void {
		fetch('/api/v1/admin/pages', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(data)
		})
			.then(async (res) => {
				const json = await res.json();
				if (!res.ok) {
					toast.error(json.message || '❌ Failed to save');
					return;
				}
				toast.success('✅ Page created');
			})
			.catch((err) => {
				console.error(err);
				toast.error('❌ Create failed');
			});
	}
</script>

<PageForm content={newPage} onsubmit={handleCreate} />
