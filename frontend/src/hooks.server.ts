
// src/hooks.server.ts
import type { Handle } from '@sveltejs/kit';
import { getUserBySession } from '$lib/server/auth'; // adjust the import path as needed


export const handle: Handle = async ({ event, resolve }) => {

	const session = event.cookies.get('user_session'); // not 'session'
	console.log('Session from cookie:', session); // Add this

	if (session) {
		const user = await getUserBySession(session);
		console.log('User from session:', user); // Add this
		if (user) {
			event.locals.user = user;
		}
	}

	return resolve(event);
};




