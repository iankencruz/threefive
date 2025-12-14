// frontend/src/routes/admin/pages/[id]/preview/+page.ts
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { PUBLIC_API_URL } from "$env/static/public";

export const load: PageServerLoad = async ({ fetch, params }) => {
  const response = await fetch(
    `${PUBLIC_API_URL}/api/v1/admin/pages/${params.id}`,
    {
      credentials: "include",
    },
  );

  if (!response.ok) {
    throw error(response.status, "Failed to fetch page");
  }

  const data = await response.json();

  return {
    page: data.data || data,
  };
};
