// frontend/src/lib/types/projects.ts
import type { Media } from "./media";
import type { SEOData } from "./seo";

// Define custom TypeScript types that map to your PostgreSQL ENUMs
type PageStatus = "draft" | "published" | "archived";
type ProjectStatus = "completed" | "ongoing" | "archived";

export interface Project {
  id: string;
  created_at: Date;
  updated_at: Date;

  title: string;
  slug: string;
  description: string | null;

  project_date: string | null;
  client_name: string | null;
  project_year: number | null;
  project_url: string | null;

  status: PageStatus;
  project_status: ProjectStatus;

  technologies: string[];

  featured_image_id: string | null;
  media: Media[]; // NEW: Array of project media
  featured_image?: Media | null; // NEW: Full featured image object

  published_at: Date | null;
  deleted_at: Date | null;

  seo?: SEOData | null;
}

// Project Mocks
export const Projects: Project[] = [
  {
    id: "a1b2c3d4-e5f6-7890-abcd-ef0123456789",
    created_at: new Date(),
    updated_at: new Date(),
    title: "Project Alpha",
    slug: "project-alpha",
    project_status: "completed",
    status: "published",
    description:
      "A powerful app for managing daily tasks and schedules efficiently.",
    project_date: "2024-05-15",
    client_name: "TechCorp Inc.",
    project_year: 2024,
    project_url: "https://project-alpha.com",
    technologies: ["SvelteKit", "PostgreSQL", "TailwindCSS"],
    featured_image_id: null,
    media: [],
    published_at: new Date(),
    deleted_at: null,
    seo: {
      meta_title: "Project Alpha",
      meta_description: "Task management app",
    },
  },
  {
    id: "b1c2d3e4-f5a6-7890-bcde-f01234567890",
    created_at: new Date(),
    updated_at: new Date(),
    title: "Project Beta",
    slug: "project-beta",
    project_status: "ongoing",
    status: "draft",
    description:
      "Developing a modern e-commerce platform with a focus on speed.",
    project_date: null,
    client_name: "RetailGiant LLC",
    project_year: 2025,
    project_url: "https://project-beta.com",
    technologies: ["React", "Node.js", "MongoDB"],
    featured_image_id: null,
    media: [],
    published_at: null,
    deleted_at: null,
    seo: {
      meta_title: "Project Beta",
      meta_description: "E-commerce platform",
    },
  },
  {
    id: "c1d2e3f4-a5b6-7890-cdef-012345678901",
    created_at: new Date(),
    updated_at: new Date(),
    title: "Project Gamma",
    slug: "project-gamma",
    project_status: "archived",
    status: "archived",
    description:
      "A cutting-edge solution for data visualization and analytics.",
    project_date: "2023-11-01",
    client_name: null,
    project_year: 2023,
    project_url: "https://project-gamma.com",
    technologies: ["Vue.js", "Python", "D3.js"],
    featured_image_id: "12345678-90ab-cdef-1234-567890abcdef",
    media: [],
    published_at: new Date("2023-12-01"),
    deleted_at: null,
    seo: {
      meta_title: "Project Gamma",
      meta_description: "Data visualization tool",
    },
  },
];
