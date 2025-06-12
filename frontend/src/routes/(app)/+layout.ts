import { browser } from '$app/environment';
import { redirect } from '@sveltejs/kit';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async ({ url, fetch }) => {
	// Only run on client-side for CSR
	if (!browser) {
		return {
			user: { id: 0 },
			isLoading: true
		};
	}

	let user = null;

	try {
		const res = await fetch('/api/v1/admin/me', {
			credentials: 'include'
		});
		const result = await res.json();

		if (res.ok && result.user?.id) {
			user = result.user;
		}
	} catch (error) {
		// User is not authenticated
		console.log('Failed to fetch user:', error);
	}

	// Handle authentication redirects for admin routes
	if (url.pathname.startsWith('/admin')) {
		// Block unauthenticated access to any /admin route except login
		if (!user || user.id === 0) {
			if (!url.pathname.startsWith('/admin/login')) {
				throw redirect(302, '/admin/login');
			}
		} else {
			// Prevent logged-in users from seeing login page again
			if (url.pathname === '/admin/login') {
				throw redirect(302, '/admin/dashboard');
			}
		}
	}

	return {
		user: user || { id: 0 },
		isLoading: false
	};
};
