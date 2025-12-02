// // frontend/src/routes/+page.ts
// import { error, type LoadEvent } from "@sveltejs/kit";
// import type { PageLoad } from "./$types";
// import { PUBLIC_API_URL } from "$env/static/public";
//
// interface MediaItem {
// 	media_id?: string;
// }
//
// export const load: PageLoad = async ({ fetch }: LoadEvent) => {
// 	try {
// 		// Fetch the page with slug "home"
// 		const response = await fetch(`${PUBLIC_API_URL}/api/v1/pages/home`);
//
// 		if (!response.ok) {
// 			if (response.status === 404) {
// 				throw error(404, {
// 					message: "Homepage not found",
// 				});
// 			}
// 			throw error(response.status, "Failed to load homepage");
// 		}
//
// 		const page = await response.json();
//
// 		// Only show published pages
// 		if (page.status !== "published") {
// 			throw error(404, {
// 				message: "Homepage not found",
// 			});
// 		}
//
// 		// ✨ Pre-fetch all media for blocks
// 		const mediaMap = new Map();
//
// 		if (page.blocks && Array.isArray(page.blocks)) {
// 			// Collect all media IDs from blocks
// 			const mediaIds = new Set<string>();
//
// 			for (const block of page.blocks) {
// 				if (block.data?.image_id) {
// 					mediaIds.add(block.data.image_id);
// 				}
// 				// Add gallery images if present
// 				if (block.data?.images && Array.isArray(block.data.images)) {
// 					block.data.images.forEach((img: MediaItem) => {
// 						if (img.media_id) mediaIds.add(img.media_id);
// 					});
// 				}
// 			}
//
// 			// Fetch all media in parallel
// 			const mediaPromises = Array.from(mediaIds).map(async (id) => {
// 				try {
// 					const res = await fetch(`${PUBLIC_API_URL}/api/v1/media/${id}`, {
// 						credentials: "include",
// 					});
// 					if (res.ok) {
// 						const media = await res.json();
// 						return [id, media];
// 					}
// 				} catch (err) {
// 					console.error(`Failed to load media ${id}:`, err);
// 				}
// 				return [id, null];
// 			});
//
// 			const mediaResults = await Promise.all(mediaPromises);
// 			mediaResults.forEach(([id, media]) => {
// 				if (media) mediaMap.set(id, media);
// 			});
// 		}
//
// 		return {
// 			page,
// 			mediaMap: Object.fromEntries(mediaMap),
// 		};
// 	} catch (err) {
// 		console.error("Error loading homepage:", err);
//
// 		if (err && typeof err === "object" && "status" in err) {
// 			throw err;
// 		}
//
// 		throw error(500, "Failed to load homepage");
// 	}
// };
