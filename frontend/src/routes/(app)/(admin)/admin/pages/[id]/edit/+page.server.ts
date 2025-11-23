import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { PUBLIC_API_URL } from "$env/static/public";

export const load: PageServerLoad = async ({ params, fetch, parent }) => {
  const { user } = await parent();
  try {
    const response = await fetch(
      `${PUBLIC_API_URL}/api/v1/admin/pages/${params.id}`,
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

    const json = await response.json();

    return {
      user,
      page: json.data || json,
    };
  } catch (err) {
    console.error("Error fetching page:", err);
    throw error(500, "Failed to load page");
  }
};
