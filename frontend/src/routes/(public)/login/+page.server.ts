
import { redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { authenticateUser } from '$lib/server/auth'; // hypothetical function

export const actions: Actions = {
	default: async ({ cookies, request }) => {
		const data = await request.formData();
		const email = data.get('email')?.toString();
		const password = data.get('password')?.toString();

		if (!email || !password) {
			return { error: 'Missing credentials' };
		}

		const user = await authenticateUser(email, password);

		if (!user) {
			return { error: 'Invalid credentials' };
		}

		// Assume user.token or generate a token string
		const token = user.token;

		cookies.set('user_session', token, {
			path: '/',
			httpOnly: true,
			sameSite: 'lax',
			secure: false,
			maxAge: 60 * 60 * 24
		});

		throw redirect(303, '/admin/dashboard');
	}
};

