// frontend/src/routes/[slug]/+page.ts
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { PUBLIC_API_URL } from "$env/static/public";

export const load: PageLoad = async ({ params, fetch }) => {
  try {
    const response = await fetch(
      `${PUBLIC_API_URL}/api/v1/pages/${params.slug}`,
    );

    if (!response.ok) {
      if (response.status === 404) {
        throw error(404, {
          message: "Page not found",
        });
      }
      throw error(response.status, "Failed to load page");
    }

    const page = await response.json();

    // Only show published pages on the public site
    if (page.status !== "published") {
      throw error(404, {
        message: "Page not found",
      });
    }

    return {
      page,
    };
  } catch (err) {
    console.error("Error loading page:", err);

    if (err && typeof err === "object" && "status" in err) {
      throw err; // Re-throw SvelteKit errors
    }

    throw error(500, "Failed to load page");
  }
};
