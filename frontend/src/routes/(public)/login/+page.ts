export const ssr = false
export const csr = true;


import { redirect } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ parent }) => {
	const { user } = await parent();

	if (user) {
		throw redirect(302, '/admin');
	}

	return {};
};

