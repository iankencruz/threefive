// frontend/src/routes/[slug]/+page.ts
import { error, redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { PUBLIC_API_URL } from "$env/static/public";

export const load: PageServerLoad = async ({ params, fetch }) => {
  // Block direct access to /home - redirect to /
  if (params.slug === "home") {
    throw redirect(308, "/");
  }

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

    // ✨ NEW: Pre-fetch all media for blocks
    const mediaMap = new Map();

    if (page.blocks && Array.isArray(page.blocks)) {
      // Collect all media IDs from blocks
      const mediaIds = new Set<string>();

      for (const block of page.blocks) {
        if (block.data?.image_id) {
          mediaIds.add(block.data.image_id);
        }
        // Add more media fields if needed (e.g., gallery images)
        if (block.data?.images && Array.isArray(block.data.images)) {
          block.data.images.forEach((img: any) => {
            if (img.media_id) mediaIds.add(img.media_id);
          });
        }
      }

      // Fetch all media in parallel
      const mediaPromises = Array.from(mediaIds).map(async (id) => {
        try {
          const res = await fetch(`${PUBLIC_API_URL}/api/v1/media/${id}`, {
            credentials: "include",
          });
          if (res.ok) {
            const media = await res.json();
            return [id, media];
          }
        } catch (err) {
          console.error(`Failed to load media ${id}:`, err);
        }
        return [id, null];
      });

      const mediaResults = await Promise.all(mediaPromises);
      mediaResults.forEach(([id, media]) => {
        if (media) mediaMap.set(id, media);
      });
    }

    return {
      page,
      mediaMap: Object.fromEntries(mediaMap), // Convert to plain object for serialization
    };
  } catch (err) {
    console.error("Error loading page:", err);

    if (err && typeof err === "object" && "status" in err) {
      throw err; // Re-throw SvelteKit errors
    }

    throw error(500, "Failed to load page");
  }
};
