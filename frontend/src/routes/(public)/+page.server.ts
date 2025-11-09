// frontend/src/routes/(public)/+page.server.ts
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { PUBLIC_API_URL } from "$env/static/public";

export const load: PageServerLoad = async ({ fetch }) => {
	try {
		// Fetch the "home" page
		const response = await fetch(`${PUBLIC_API_URL}/api/v1/pages/home`);

		if (!response.ok) {
			throw error(response.status, "Failed to load homepage");
		}

		const page = await response.json();

		// Only show published pages
		if (page.status !== "published") {
			throw error(404, {
				message: "Page not found",
			});
		}

		// ✨ Pre-fetch all media for blocks (same logic as [slug])
		const mediaMap = new Map();

		if (page.blocks && Array.isArray(page.blocks)) {
			const mediaIds = new Set<string>();

			for (const block of page.blocks) {
				// Hero block image
				if (block.data?.image_id) {
					mediaIds.add(block.data.image_id);
				}
				// Gallery block media
				if (block.data?.media && Array.isArray(block.data.media)) {
					block.data.media.forEach((media: any) => {
						if (media.id) mediaIds.add(media.id);
					});
				}
				// Legacy images array
				if (block.data?.images && Array.isArray(block.data.images)) {
					block.data.images.forEach((img: any) => {
						if (img.media_id) mediaIds.add(img.media_id);
					});
				}
			}

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
			mediaMap: Object.fromEntries(mediaMap),
		};
	} catch (err) {
		console.error("Error loading homepage:", err);

		if (err && typeof err === "object" && "status" in err) {
			throw err;
		}

		throw error(500, "Failed to load homepage");
	}
};
