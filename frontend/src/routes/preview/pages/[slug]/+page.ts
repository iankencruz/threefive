// frontend/src/routes/preview/pages/[slug]/+page.ts
import { error } from "@sveltejs/kit";
import { PUBLIC_API_URL } from "$env/static/public";

interface MediaItem {
	media_id?: string;
}

interface Media {
	id: string;
	url: string;
	mime_type: string;
	thumbnail_url?: string;
	large_url?: string;
}

export const load = async ({
	params,
	fetch,
}: {
	params: { slug: string };
	fetch: typeof globalThis.fetch;
}) => {
	try {
		// Fetch page by slug - this gets ANY page regardless of status
		const response = await fetch(
			`${PUBLIC_API_URL}/api/v1/pages/${params.slug}`,
			{
				credentials: "include",
			},
		);

		if (!response.ok) {
			if (response.status === 404) {
				throw error(404, "Page not found");
			}
			throw error(response.status, "Failed to load page preview");
		}

		const page = await response.json();

		// ✨ Pre-fetch all media for blocks
		const mediaMap = new Map<string, Media>();

		if (page.blocks && Array.isArray(page.blocks)) {
			const mediaIds = new Set<string>();

			for (const block of page.blocks) {
				if (block.data?.image_id) {
					mediaIds.add(block.data.image_id);
				}
				if (block.data?.images && Array.isArray(block.data.images)) {
					block.data.images.forEach((img: MediaItem) => {
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
						return [id, media] as const;
					}
				} catch (err) {
					console.error(`Failed to load media ${id}:`, err);
				}
				return [id, null] as const;
			});

			const mediaResults = await Promise.all(mediaPromises);
			mediaResults.forEach(([id, media]) => {
				if (media) mediaMap.set(id, media);
			});
		}

		return {
			page,
			mediaMap: Object.fromEntries(mediaMap),
			isPreview: true,
		};
	} catch (err) {
		console.error("Error loading page preview:", err);

		if (err && typeof err === "object" && "status" in err) {
			throw err;
		}

		throw error(500, "Failed to load page preview");
	}
};
