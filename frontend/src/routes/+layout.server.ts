
// src/routes/+layout.server.ts
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals }) => {
	console.log('[layout.server.ts] locals.user:', locals.user);

	return {
		user: locals.user ?? null
	};
};

