import type { SEOData } from "./seo";

// Define custom TypeScript types that map to your PostgreSQL ENUMs
type PageStatus = "draft" | "published" | "archived";
type ProjectStatus = "completed" | "in_progress" | "on_hold" | string;

export interface Project {
  id: string;
  created_at: Date;
  updated_at: Date;

  title: string;
  slug: string;
  description: string | null;

  project_date: Date | null;
  client_name: string | null;
  project_year: number | null;
  project_url: string | null;

  status: PageStatus | null;
  project_status: ProjectStatus;

  technologies: string[] | { [key: string]: any } | null;

  featured_image_id: string | null;

  published_at: Date | null;
  deleted_at: Date | null;

  seo: SEOData;
}
