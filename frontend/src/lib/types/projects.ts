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

// Project Mocks
export const Projects: Project[] = [
  {
    // Minimal required fields
    id: "a1b2c3d4-e5f6-7890-abcd-ef0123456789",
    created_at: new Date(),
    updated_at: new Date(),
    title: "Project Alpha",
    slug: "project-alpha",
    project_status: "completed",
    seo: { meta_title: "Project Alpha" },

    // Fields used in your Svelte component/required fields
    description:
      "A powerful app for managing daily tasks and schedules efficiently.",

    // Remaining optional/nullable fields
    project_date: new Date("2024-05-15"),
    client_name: "TechCorp Inc.",
    project_year: 2024,
    project_url: "/projects/project-alpha", // Matches the original 'href' intent
    status: "published",
    technologies: ["SvelteKit", "PostgreSQL"],
    featured_image_id: null,
    published_at: new Date(),
    deleted_at: null,
  },
  {
    id: "b1c2d3e4-f5a6-7890-bcde-f01234567890",
    created_at: new Date(),
    updated_at: new Date(),
    title: "Project Beta",
    slug: "project-beta",
    project_status: "in_progress",
    seo: { meta_title: "Project Beta" },

    description:
      "Developing a modern e-commerce platform with a focus on speed.",

    project_date: null,
    client_name: "RetailGiant LLC",
    project_year: 2025,
    project_url: "/projects/project-beta",
    status: "draft",
    technologies: ["React", "Node.js"],
    featured_image_id: null,
    published_at: null,
    deleted_at: null,
  },
  {
    id: "c1d2e3f4-a5b6-7890-cdef-012345678901",
    created_at: new Date(),
    updated_at: new Date(),
    title: "Project Gamma",
    slug: "project-gamma",
    project_status: "on_hold",
    seo: { meta_title: "Project Gamma" },

    description:
      "A cutting-edge solution for data visualization and analytics.",

    project_date: new Date("2023-11-01"),
    client_name: null,
    project_year: 2023,
    project_url: "/projects/project-gamma",
    status: "archived",
    technologies: ["Vue.js", "Python"],
    featured_image_id: "12345678-90ab-cdef-1234-567890abcdef",
    published_at: new Date("2023-12-01"),
    deleted_at: null,
  },
  {
    id: "d1e2f3a4-b5c6-7890-def0-1234567890ab",
    created_at: new Date(),
    updated_at: new Date(),
    title: "Project Delta",
    slug: "project-delta",
    project_status: "completed",
    seo: { meta_title: "Project Delta" },
    description:
      "An innovative mobile app designed to enhance user productivity.",
    project_date: new Date("2022-08-20"),
    client_name: "AppMakers Co.",
    project_year: 2022,
    project_url: "/projects/project-delta",
    status: "published",
    technologies: { frontend: "Flutter", backend: "Firebase" },
    featured_image_id: null,
    published_at: new Date("2022-09-15"),
    deleted_at: null,
  },
];
