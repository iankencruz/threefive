// frontend/src/lib/api/media.ts
import { PUBLIC_API_URL } from "$env/static/public";
import type { Media, MediaListParams, MediaListResponse } from "$types/media";

async function fetchAPI<T>(url: string, options: RequestInit = {}): Promise<T> {
  const response = await fetch(url, {
    credentials: "include",
    headers: {
      ...options.headers,
    },
    ...options,
  });

  const data = await response.json();

  if (!response.ok) {
    throw data;
  }

  return data as T;
}

// Helper function to get the correct media URL
export function getMediaUrl(media: Media): string {
  // If url exists, use it
  if (media.url) return media.url;

  // Otherwise, build URL from storage_path
  const baseUrl = PUBLIC_API_URL || "http://localhost:8080";
  return `${baseUrl}/uploads/${media.filename}`;
}

export const mediaApi = {
  async listMedia(params: MediaListParams = {}) {
    const queryParams = new URLSearchParams();

    // Add page and limit (with defaults)
    queryParams.append("page", (params.page || 1).toString());
    queryParams.append("limit", (params.limit || 20).toString());

    // Add optional filters
    if (params.search) queryParams.append("search", params.search);
    if (params.type) queryParams.append("type", params.type);
    if (params.sort) queryParams.append("sort", params.sort);
    if (params.order) queryParams.append("order", params.order);

    return fetchAPI<MediaListResponse>(
      `${PUBLIC_API_URL}/api/v1/media?${queryParams.toString()}`,
    );
  },

  async uploadMedia(file: File): Promise<Media> {
    const formData = new FormData();
    formData.append("file", file);

    const response = await fetch(`${PUBLIC_API_URL}/api/v1/media/upload`, {
      method: "POST",
      credentials: "include",
      body: formData,
    });

    const data = await response.json();
    if (!response.ok) throw data;
    return data;
  },

  async deleteMedia(id: string): Promise<{ message: string }> {
    return fetchAPI(`${PUBLIC_API_URL}/api/v1/media/${id}`, {
      method: "DELETE",
    });
  },
};

export type { Media, MediaListResponse, MediaListParams };
