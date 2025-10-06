// frontend/src/routes/admin/pages/+page.server.ts
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ fetch, url, parent }) => {
  // Get user from parent layout
  const { user } = await parent();

  const page = parseInt(url.searchParams.get("page") || "1");
  const limit = 20;

  try {
    const response = await fetch(`/api/v1/pages?page=${page}&limit=${limit}`);

    if (!response.ok) {
      throw error(response.status, "Failed to fetch pages");
    }

    const data = await response.json();

    return {
      user,
      pages: data.pages || [],
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
