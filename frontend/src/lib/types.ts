import type { Component } from 'svelte';
export interface SEOData {
  title?: string;
  description?: string;
  keywords?: string[];
  ogImage?: string;
  canonical?: string;
}

export interface LayoutProps {
  title: string;
  path: string;
  isAdmin?: boolean;
  seo?: SEOData | null;
}

export interface NavItem {
  title: string;
  href: string;
  icon?: Component<any>;
}
