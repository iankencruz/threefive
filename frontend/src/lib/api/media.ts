import { PUBLIC_API_URL } from "$env/static/public";

interface Media {
  id: string;
  filename: string;
  original_filename: string;
  mime_type: string;
  size_bytes: number;
  width?: number;
  height?: number;
  url: string;
  thumbnail_url?: string;
  created_at: string;
}

interface MediaListResponse {
  data: Media[];
  pagination?: {
    page: number;
    limit: number;
    total?: number;
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
