
// src/lib/stores/pages.svelte.ts

import type { Page } from '$lib/types';


export class PageStore {
	pagesStore = $state<Page[]>([]);

	setPages(p: Page[]) {
		this.pagesStore.splice(0, this.pagesStore.length, ...p);
	}

	updatePage(updated: Page) {
		const index = this.pagesStore.findIndex((p) => p.id === updated.id);
		if (index !== -1) {
			this.pagesStore[index] = updated;
		}
	}

	clear() {
		this.pagesStore.length = 0;
	}

	async fetchPages() {
		try {
			const res = await fetch('/api/v1/admin/pages');
			if (!res.ok) throw new Error('Failed to fetch pages');
			const json = await res.json();
			this.setPages(json.data);
		} catch (err) {
			console.error('❌ Failed to fetch pages:', err);
		}
	}
}

export const pages = new PageStore();

