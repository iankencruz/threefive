

import type { User } from '$lib/types';

export class AuthStore {
	user = $state<User | null>(null);

	setUser(u: User) {
		this.user = u;
	}

	clear() {
		this.user = null;
	}

	get isAuthenticated(): boolean {
		return !!this.user;
	}

	async logout() {
		try {
			await fetch('/api/v1/auth/logout', { method: 'POST' });
		} catch (err) {
			console.error('Logout request failed:', err);
		} finally {
			this.clear();
		}
	}
}

export const auth = new AuthStore();

