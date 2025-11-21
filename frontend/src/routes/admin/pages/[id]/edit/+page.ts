export const ssr = false;
export const csr = true;

import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { PUBLIC_API_URL } from "$env/static/public";

export const load: PageLoad = async ({ params, fetch, parent }) => {
  const { user } = await parent();
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
