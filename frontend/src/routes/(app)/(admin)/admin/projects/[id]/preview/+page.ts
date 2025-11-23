// frontend/src/routes/admin/pages/[id]/preview/+page.ts
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { PUBLIC_API_URL } from "$env/static/public";

export const load: PageLoad = async ({ fetch, params }) => {
  const response = await fetch(
    `${PUBLIC_API_URL}/api/v1/admin/projects/${params.id}`,
    {
      credentials: "include",
    },
  );

  if (!response.ok) {
    throw error(response.status, "Failed to fetch project");
  }

  const data = await response.json();

  return {
    project: data.data || data,
  };
};
