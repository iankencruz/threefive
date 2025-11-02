// Add this at the top of +page.ts to disable caching
export const ssr = false;
export const csr = true;

import { error } from "@sveltejs/kit";
import type { PageLoad } from "../../[id]/edit/$types";
import { PUBLIC_API_URL } from "$env/static/public";

export const load: PageLoad = async ({ params, fetch }) => {
	try {
		const response = await fetch(
			`${PUBLIC_API_URL}/api/v1/pages/${params.id}`,
			{
				credentials: "include",
			},
		);

		if (!response.ok) {
			if (response.status === 404) {
				throw error(404, "Page not found");
			}
			throw error(response.status, "Failed to load page");
		}

		const page = await response.json();
		// console.log("✅ Loaded page:", page.id, page.title);

		return {
			page,
			isPreview: true,
		};
	} catch (err) {
		console.error("Error loading page:", err);
		throw error(500, "Failed to load page");
	}
};
