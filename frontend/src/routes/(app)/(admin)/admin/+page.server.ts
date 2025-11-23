// src/routes/admin/+page.server.ts

// This file is empty because all it needs to do is exist.
// Its existence makes /admin a valid route, allowing the redirect
// in +layout.server.ts to execute cleanly.
import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "../$types";

export const load: PageServerLoad = async ({ parent, url }) => {
	// Note: The redirect is already in +layout.server.ts, but
	// for redundancy/safety, you could also place it here:
	// redirect(303, "/admin/dashboard");

	const { user } = await parent();

	// Redirect to login if not authenticated
	if (!user) {
		redirect(303, `/auth/login`);
	}

	// 2. Redirect the root /admin path to /admin/dashboard
	// Check if the current pathname is EXACTLY /admin (or /admin/)
	if (url.pathname === "/admin" || url.pathname === "/admin/") {
		// You would typically redirect to a specific dashboard page,
		// e.g., /admin/dashboard
		redirect(303, `/admin/dashboard`);
	}
};
