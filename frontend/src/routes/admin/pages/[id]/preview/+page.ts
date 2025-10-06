// frontend/src/routes/admin/pages/[id]/preview/+page.ts
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { PUBLIC_API_URL } from "$env/static/public";

export const load: PageLoad = async ({ params, fetch }) => {
  try {
    // Fetch page by ID (not slug) - this gets ANY page regardless of status
    const response = await fetch(
      `${PUBLIC_API_URL}/api/v1/pages/${params.id}`,
      {
        credentials: "include", // Include auth cookies
      },
    );

    if (!response.ok) {
      if (response.status === 404) {
        throw error(404, "Page not found");
      }
      throw error(response.status, "Failed to load page preview");
    }

    const page = await response.json();

    return {
      page,
      isPreview: true, // Flag to show preview banner
    };
  } catch (err) {
    console.error("Error loading page preview:", err);

    if (err && typeof err === "object" && "status" in err) {
      throw err;
    }

    throw error(500, "Failed to load page preview");
  }
};
