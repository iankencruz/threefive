// frontend/src/routes/(public)/blog/[slug]/+page.ts
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { PUBLIC_API_URL } from "$env/static/public";

export const load: PageLoad = async ({ params, fetch }) => {
	try {
		// Blog posts have prefixed slugs: "blog/slug-name"
		const fullSlug = `blog/${params.slug}`;

		const response = await fetch(`${PUBLIC_API_URL}/api/v1/pages/${fullSlug}`);

		if (!response.ok) {
			if (response.status === 404) {
				throw error(404, {
					message: "Blog post not found",
				});
			}
			throw error(response.status, "Failed to load blog post");
		}

		const page = await response.json();

		// Only show published pages on the public site
		if (page.status !== "published") {
			throw error(404, {
				message: "Blog post not found",
			});
		}

		// ✨ Pre-fetch all media for blocks
		const mediaMap = new Map();

		if (page.blocks && Array.isArray(page.blocks)) {
			const mediaIds = new Set<string>();

			for (const block of page.blocks) {
				if (block.data?.image_id) {
					mediaIds.add(block.data.image_id);
				}
				if (block.data?.media && Array.isArray(block.data.media)) {
					block.data.media.forEach((media: any) => {
						if (media.id) mediaIds.add(media.id);
					});
				}
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
		console.error("Error loading blog post:", err);

		if (err && typeof err === "object" && "status" in err) {
			throw err;
		}

		throw error(500, "Failed to load blog post");
	}
};
