import type { SEOData } from "./seo";

type BlogStatus = "draft" | "published" | "archived";

export interface Blog {
  id: string;
  title: string;
  slug: string;
  status: BlogStatus | null;
  excerpt: string | null;
  reading_time: number | null;
  is_featured: boolean | null;
  featured_image_id: string | null;
  created_at: Date;
  updated_at: Date;
  published_at: Date | null;
  deleted_at: Date | null;
  seo: SEOData;
}
