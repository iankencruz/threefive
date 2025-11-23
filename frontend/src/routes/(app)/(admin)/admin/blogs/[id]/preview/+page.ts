// frontend/src/routes/admin/blog/[id]/preview/+page.ts
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { PUBLIC_API_URL } from "$env/static/public";

export const load: PageLoad = async ({ fetch, params }) => {
  const response = await fetch(
    `${PUBLIC_API_URL}/api/v1/admin/blogs/${params.id}`,
    {
      credentials: "include",
    },
  );

  if (!response.ok) {
    throw error(response.status, "Failed to fetch blogs");
  }

  const data = await response.json();

  return {
    blog: data.data || data,
  };
};
