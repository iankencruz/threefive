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

// --- mock data ---

// Helper function to generate a date offset in the past
const daysAgo = (days: number): Date => {
  const date = new Date();
  date.setDate(date.getDate() - days);
  return date;
};

/**
 * --- MOCK BLOG DATA ---
 * Array containing mock Blog objects for testing purposes.
 */
export const Blogs: Blog[] = [
  {
    id: "blog-001",
    title: "The Future of Web Development: Signals and AI Integration",
    slug: "future-web-development-signals-ai",
    status: "published",
    excerpt:
      "Exploring the shift towards reactive programming paradigms and how AI tools are fundamentally changing the development workflow and time-to-market for complex applications.",
    reading_time: 7,
    is_featured: true,
    featured_image_id: "image-ai-signal-001",
    created_at: daysAgo(30),
    updated_at: daysAgo(1),
    published_at: daysAgo(5),
    deleted_at: null,
    seo: {
      meta_title: "Web Dev Future: Signals & AI | MockBlog",
      meta_description:
        "A deep dive into the latest trends shaping web development, including signal-based state management and integrating generative AI tools.",
      og_title: "The Future of Web Development: Signals and AI Integration",
      og_description:
        "Read our comprehensive analysis on Signals and AI in web dev.",
      canonical_url: "https://mockblog.com/future-web-development-signals-ai",
    },
  },
  {
    id: "blog-002",
    title: "Advanced TypeScript Utility Types for Enterprise Scale",
    slug: "typescript-utility-types-enterprise",
    status: "published",
    excerpt:
      "Beyond Partial and Readonly: A look at advanced techniques using conditional types, template literal types, and key remapping to build robust, scalable interfaces.",
    reading_time: 12,
    is_featured: false,
    featured_image_id: "image-ts-util-002",
    created_at: daysAgo(180),
    updated_at: daysAgo(100),
    published_at: daysAgo(150),
    deleted_at: null,
    seo: {
      meta_title: "Advanced TypeScript Types | MockBlog",
      meta_description:
        "Master TypeScript utility types for enterprise applications.",
      og_title: "Advanced TypeScript Utility Types",
      og_description:
        "Learn how to use complex TS features to build better software.",
      canonical_url: "https://mockblog.com/typescript-utility-types-enterprise",
    },
  },
  {
    id: "blog-003",
    title: "Initial Draft: Understanding CSS Grid Layout",
    slug: "initial-draft-css-grid-layout",
    status: "draft",
    excerpt:
      "This is the first draft of a tutorial on CSS Grid. It currently covers basic syntax and a few examples of responsive two-column layouts. Needs review and expansion on nesting.",
    reading_time: 4,
    is_featured: false,
    featured_image_id: null,
    created_at: daysAgo(3),
    updated_at: daysAgo(3),
    published_at: null, // Draft
    deleted_at: null,
    seo: {
      meta_title: "DRAFT: CSS Grid Basics",
      meta_description: "Draft content for an upcoming CSS Grid tutorial.",
      og_title: "Understanding CSS Grid Layout (WIP)",
      og_description: "Internal draft review.",
      canonical_url: "https://mockblog.com/initial-draft-css-grid-layout",
    },
  },
  {
    id: "blog-004",
    title: "The State of Functional Programming in 2022",
    slug: "fp-state-of-2022",
    status: "archived",
    excerpt:
      "An older analysis of the adoption of functional programming concepts in mainstream JavaScript frameworks, written back in late 2022. It still provides historical context.",
    reading_time: 6,
    is_featured: false,
    featured_image_id: "image-fp-2022-004",
    created_at: daysAgo(700),
    updated_at: daysAgo(600),
    published_at: daysAgo(650),
    deleted_at: null,
    seo: {
      meta_title: "Archived: Functional Programming 2022",
      meta_description: "Historical look at FP adoption in 2022.",
      og_title: "The State of Functional Programming in 2022",
      og_description: "A functional programming retrospective.",
      canonical_url: "https://mockblog.com/fp-state-of-2022",
    },
  },
  {
    id: "blog-005",
    title: "Optimizing Database Queries: A Deep Dive into Indexing",
    slug: "optimizing-db-queries-indexing",
    status: "published",
    excerpt:
      "Learn how to effectively use B-tree and Hash indexes to drastically reduce query execution time and improve application performance across various SQL and NoSQL databases.",
    reading_time: 15,
    is_featured: true,
    featured_image_id: "image-db-index-005",
    created_at: daysAgo(50),
    updated_at: daysAgo(2),
    published_at: daysAgo(10),
    deleted_at: null,
    seo: {
      meta_title: "DB Query Optimization with Indexing | MockBlog",
      meta_description:
        "A comprehensive guide to indexing strategies for performance optimization.",
      og_title: "Optimizing Database Queries: A Deep Dive into Indexing",
      og_description: "Boost your database performance today!",
      canonical_url: "https://mockblog.com/optimizing-db-queries-indexing",
    },
  },
];
