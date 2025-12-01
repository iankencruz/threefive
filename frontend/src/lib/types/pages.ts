import type { SEOData } from "./seo";
// Define custom TypeScript types that map to your PostgreSQL ENUM type
type PageStatus = "draft" | "published" | "archived";

export interface Page {
  id: string;
  title: string;
  slug: string;
  status: PageStatus | null;
  featured_image_id: string | null;
  created_at: Date;
  updated_at: Date;
  published_at: Date | null;
  deleted_at: Date | null;
  seo: SEOData;
}

export interface Link {
  id: number;
  title: string;
  href: string;
}
