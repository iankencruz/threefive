export interface Media {
  id: string;
  filename: string;
  original_filename: string;
  mime_type: string;
  url: string;
  original_url?: string;
  large_url?: string;
  medium_url?: string;
  thumbnail_url?: string;
  width?: number;
  height?: number;
  size_bytes: number;
  created_at: string;
}

export interface MediaListResponse {
  data: Media[];
  pagination?: {
    page: number;
    limit: number;
    total?: number;
    total_pages?: number;
  };
  filters?: {
    search: string;
    type: string;
    sort: string;
    order: string;
  };
}

export interface MediaListParams {
  page?: number;
  limit?: number;
  search?: string;
  type?: string; // 'image' or 'video'
  sort?: "created_at" | "filename" | "size";
  order?: "asc" | "desc";
}
