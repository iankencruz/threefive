// frontend/src/routes/admin/projects/+page.ts
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { PUBLIC_API_URL } from "$env/static/public";

export const load: PageServerLoad = async ({ fetch, url, parent }) => {
  const { user } = await parent();

  const page = parseInt(url.searchParams.get("page") || "1");
  const limit = 20;

  try {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString(),
    });

    const response = await fetch(
      `${PUBLIC_API_URL}/api/v1/admin/projects?${params.toString()}`,
      {
        credentials: "include",
      },
    );

    if (!response.ok) {
      throw error(response.status, "Failed to fetch projects");
    }

    const data = await response.json();

    return {
      user,
      projects: data || [],
      pagination: data.pagination || {
        page: 1,
        limit: 20,
        total_pages: 0,
        total_count: 0,
      },
    };
  } catch (err) {
    console.error("Error fetching projects:", err);
    throw error(500, "Failed to load projects");
  }
};
