import type { SEOData } from "./seo";

export interface PageContent {
  title: string;
  slug: string;
  status: "draft" | "published" | "archived";
  blocks: any[];
  seo: SEOData;
}

export interface Link {
  id: number;
  title: string;
  href: string;
}
