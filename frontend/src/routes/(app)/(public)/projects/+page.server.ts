// frontend/src/routes/(app)/(public)/projects/+page.server.ts
import type { PageServerLoad } from "./$types";
import { PUBLIC_API_URL } from "$env/static/public";
import { error } from "@sveltejs/kit";

export const load: PageServerLoad = async ({ fetch }) => {
  console.log("Fetching from:", `${PUBLIC_API_URL}/api/v1/projects`);

  const response = await fetch(`${PUBLIC_API_URL}/api/v1/projects`);

  // console.log("Response status:", response.status);
  // console.log("Response ok:", response.ok);

  if (!response.ok) {
    throw error(response.status, "Failed to fetch projects");
  }

  const result = await response.json();
  // console.log("Full API result:", result);
  // console.log("Projects array:", result.data);

  return {
    projects: result.data,
    pagination: result.pagination,
  };
};
