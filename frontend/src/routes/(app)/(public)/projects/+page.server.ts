// frontend/src/routes/[slug]/+page.ts
import type { PageServerLoad } from "./$types";
import { PUBLIC_API_URL } from "$env/static/public";

export const load: PageServerLoad = async ({ fetch }) => {
  const response = await fetch(`${PUBLIC_API_URL}/api/v1/projects`);
  if (!response.ok) {
    throw new Error("Failed to fetch projects");
  }
  const data = await response.json();
  return data;
};
