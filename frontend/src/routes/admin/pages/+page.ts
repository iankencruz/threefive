// frontend/src/routes/admin/pages/+page.server.ts
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, url, parent }) => {
	// Get user from parent layout
	const { user } = await parent();

	const page = parseInt(url.searchParams.get("page") || "1");
	const limit = 20;
	const pageType = url.searchParams.get("page_type");

	try {
		// Build query parameters
		const params = new URLSearchParams({
			page: page.toString(),
			limit: limit.toString(),
		});

		// Add page_type filter if it exists
		if (pageType && pageType !== "all") {
			params.set("page_type", pageType);
		}

		const response = await fetch(`/api/v1/pages?${params.toString()}`);

		if (!response.ok) {
			throw error(response.status, "Failed to fetch pages");
		}

		const data = await response.json();
		console.log("ts - data: ", data);

		return {
			user,
			pages: data || [],
			pagination: data.pagination || {
				page: 1,
				limit: 20,
				total_pages: 0,
				total_count: 0,
			},
		};
	} catch (err) {
		console.error("Error fetching pages:", err);
		throw error(500, "Failed to load pages");
	}
};
