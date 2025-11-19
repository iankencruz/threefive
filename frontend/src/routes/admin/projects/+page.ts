import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, parent, url }) => {
	const { user } = await parent();

	const page = parseInt(url.searchParams.get("page") || "1");
	const limit = 20;

	try {
		const params = new URLSearchParams({
			page: page.toString(),
			limit: limit.toString(),
		});

		const response = await fetch(`/api/v1/projects?${params.toString()}`, {
			credentials: "include",
		});

		const data = await response.json();

		return {
			user,
			projects: data.projects || [],
			pagination: data.pagination || {
				page: 1,
				limit: 20,
				total_pages: 0,
				total_count: 0,
			},
		};
	} catch (err) {
		console.error("Error constructing URLSearchParams:", err);
		throw error(500, "Internal Server Error");
	}
};
