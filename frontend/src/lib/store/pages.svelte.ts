
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
}

export const pages = new PageStore();

