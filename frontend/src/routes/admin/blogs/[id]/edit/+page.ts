// frontend/src/routes/admin/blogs/[id]/edit/+page.ts
import { PUBLIC_API_URL } from "$env/static/public";
import type { PageLoad } from "./$types";
import { error } from "@sveltejs/kit";

export const load: PageLoad = async ({ params, fetch, parent }) => {
  const { user } = await parent();

  const response = await fetch(`${PUBLIC_API_URL}/api/v1/blogs/${params.id}`, {
    credentials: "include",
  });

  if (!response.ok) {
    if (response.status === 404) {
      throw error(404, "Blog not found");
    }
    throw error(response.status, "Failed to load blog");
  }

  const json = await response.json();

  return {
    user,
    blog: json.data || json,
  };
};
