// frontend/src/routes/(public)/projects/[slug]/+page.ts
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { PUBLIC_API_URL } from "$env/static/public";

export const load: PageLoad = async ({ params, fetch }) => {
	try {
		// Projects have prefixed slugs: "projects/slug-name"
		const fullSlug = `projects/${params.slug}`;

		console.log("🔍 Loading project with slug:", fullSlug);
		console.log("🌐 API URL:", `${PUBLIC_API_URL}/api/v1/pages/${fullSlug}`);

		const response = await fetch(`${PUBLIC_API_URL}/api/v1/pages/${fullSlug}`);

		console.log("📡 Response status:", response.status);

		if (!response.ok) {
			const errorBody = await response.json();
			console.error("❌ Error response:", errorBody);

			if (response.status === 404) {
				throw error(404, {
					message: "Project not found",
				});
			}
			throw error(response.status, "Failed to load project");
		}

		const page = await response.json();
		console.log("✅ Page loaded:", page.title, "Blocks:", page.blocks?.length);

		// Only show published pages on the public site
		if (page.status !== "published") {
			console.log("! Page not published, status:", page.status);
			throw error(404, {
				message: "Project not found",
			});
		}

		// ✨ Pre-fetch all media for blocks
		const mediaMap = new Map();

		if (page.blocks && Array.isArray(page.blocks)) {
			const mediaIds = new Set<string>();

			for (const block of page.blocks) {
				// Hero block image
				if (block.data?.image_id) {
					mediaIds.add(block.data.image_id);
					console.log("🖼 Found hero image:", block.data.image_id);
				}
				// Gallery block media
				if (block.data?.media && Array.isArray(block.data.media)) {
					block.data.media.forEach((media: any) => {
						if (media.id) {
							mediaIds.add(media.id);
							console.log("🖼 Found gallery media:", media.id);
						}
					});
				}
			}

			console.log("📦 Total media IDs to fetch:", mediaIds.size);

			// Fetch all media in parallel
			const mediaPromises = Array.from(mediaIds).map(async (id) => {
				try {
					const res = await fetch(`${PUBLIC_API_URL}/api/v1/media/${id}`, {
						credentials: "include",
					});
					if (res.ok) {
						const media = await res.json();
						console.log("✅ Fetched media:", id);
						return [id, media];
					} else {
						console.error("❌ Failed to fetch media:", id, res.status);
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

		console.log("🎉 Final mediaMap size:", mediaMap.size);

		return {
			page,
			mediaMap: Object.fromEntries(mediaMap),
		};
	} catch (err) {
		console.error("Error loading project:", err);

		if (err && typeof err === "object" && "status" in err) {
			throw err;
		}

		throw error(500, "Failed to load project");
	}
};
