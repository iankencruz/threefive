// frontend/src/lib/api/media.ts
import { PUBLIC_API_URL } from "$env/static/public";

interface Media {
  id: string;
  filename: string;
  original_filename: string;
  mime_type: string;
  size_bytes: number;
  width?: number;
  height?: number;
  storage_type: string;
  storage_path: string; // Added this field
  url?: string; // Made optional since it might be null
  thumbnail_url?: string;
  created_at: string;
}

interface MediaListResponse {
  data: Media[];
  pagination?: {
    page: number;
    limit: number;
    total?: number;
    total_pages?: number; // Added for pagination
  };
}

interface ErrorResponse {
  message?: string;
  code?: string;
}

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
  return `${baseUrl}/uploads/${media.storage_path || media.filename}`;
}

export const mediaApi = {
  async listMedia(page: number = 1, limit: number = 20) {
    return fetchAPI<MediaListResponse>(
      `${PUBLIC_API_URL}/api/v1/media?page=${page}&limit=${limit}`,
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

export type { Media, MediaListResponse, ErrorResponse };
