export interface SEOData {
  meta_title: string;
  meta_description: string;
  og_title: string;
  og_description: string;
  robots_index: boolean;
  robots_follow: boolean;
  canonical_url?: string; // ← Optional property
}
