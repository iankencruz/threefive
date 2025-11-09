// frontend/src/routes/preview/pages/[id]/+page.ts
import { error } from "@sveltejs/kit";
import type { PageLoad } from "../[slug]/$types";
import { PUBLIC_API_URL } from "$env/static/public";

interface MediaItem {
	media_id?: string;
}

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

		// ✨ Pre-fetch all media for blocks (same as public pages)
		const mediaMap = new Map();

		if (page.blocks && Array.isArray(page.blocks)) {
			// Collect all media IDs from blocks
			const mediaIds = new Set<string>();

			for (const block of page.blocks) {
				// Hero block image
				if (block.data?.image_id) {
					mediaIds.add(block.data.image_id);
				}

				// Gallery block media - backend returns media array with full objects
				if (block.data?.media && Array.isArray(block.data.media)) {
					block.data.media.forEach((media: any) => {
						if (media.id) mediaIds.add(media.id);
					});
				}

				// Legacy: images array (if you have other blocks using this)
				if (block.data?.images && Array.isArray(block.data.images)) {
					block.data.images.forEach((img: MediaItem) => {
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
