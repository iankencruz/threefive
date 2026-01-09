// frontend/src/routes/(app)/(public)/projects/[slug]/+page.server.ts
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { PUBLIC_API_URL } from "$env/static/public";

export const load: PageServerLoad = async ({ params, fetch }) => {
  try {
    const response = await fetch(
      `${PUBLIC_API_URL}/api/v1/projects/${params.slug}`,
    );

    if (!response.ok) {
      if (response.status === 404) {
        throw error(404, {
          message: "Project not found",
        });
      }
      throw error(response.status, "Failed to load project");
    }

    const project = await response.json();

    // Only show published projects on the public site
    if (project.status !== "published") {
      throw error(404, {
        message: "Project not found",
      });
    }

    // ✨ Fetch media gallery for this project
    // Projects have media linked via media_relations with entity_type='project'
    let projectMedia: any[] = [];

    try {
      const mediaResponse = await fetch(
        `${PUBLIC_API_URL}/api/v1/media/entity/project/${project.id}`,
      );

      if (mediaResponse.ok) {
        const mediaData = await mediaResponse.json();
        projectMedia = Array.isArray(mediaData)
          ? mediaData
          : mediaData.media || [];
      }
    } catch (err) {
      console.error("Failed to load project media:", err);
    }

    return {
      project,
      projectMedia,
    };
  } catch (err) {
    console.error("Error loading project:", err);

    if (err && typeof err === "object" && "status" in err) {
      throw err;
    }

    throw error(500, "Failed to load project");
  }
};
