// src/routes/admin/+layout.server.ts
import { redirect } from "@sveltejs/kit";
import type { LayoutServerLoad } from "./$types";
import { PUBLIC_API_URL } from "$env/static/public";

export const load: LayoutServerLoad = async ({ parent }) => {
	const { user } = await parent();

	// Redirect to login if not authenticated
	if (!user) {
		redirect(303, `/auth/login`);
	}

	return { user };
};
