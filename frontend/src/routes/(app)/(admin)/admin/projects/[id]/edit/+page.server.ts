import { PUBLIC_API_URL } from "$env/static/public";
import { error } from "console";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ params, fetch, parent }) => {
  const { user } = await parent();

  try {
    const response = await fetch(
      `${PUBLIC_API_URL}/api/v1/admin/projects/${params.id}`,
      {
        credentials: "include",
      },
    );

    if (!response.ok) {
      if (response.status === 404) {
        throw error(404, "Project not found");
      }
      throw error(response.status, "Failed to load project");
    }

    const json = await response.json();

    return {
      user,
      project: json.data || json,
    };
  } catch (err) {
    console.log("Error fetching projects:", err);
    throw error(500, "Failed to load project");
  }
};
