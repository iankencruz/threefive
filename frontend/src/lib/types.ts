import type { UUID } from 'crypto';

export type User = {
  id: UUID;
  first_name: string;
  last_name: string;
  email: string;
  roles: string[];
}


export type Project = {
  id: UUID;
  title: string;
  slug: string;
  description: string | null;
  meta_description: string | null;
  canonical_url: string | null;
  cover_media_id: string | null;
  is_published: boolean;
  published_at: string | null;
  created_at: string;
  updated_at: string;
  media?: MediaItem[]; // or `media?: Media[]` if it's optional

};



export type Gallery = {
  id: UUID;
  title: string;
  description: string | null;
  slug: string;
  is_published: boolean;
  published_at: string | null;
  created_at: string;
  updated_at: string;
  media?: MediaItem[]; // or `media?: Media[]` if it's optional
}



export type MediaItem = {
  id: UUID;
  title: string;
  url: string;
  thumbnail_url?: string;
  medium_url?: string;
  alt_text?: string;
  file_size?: string;
  mime_type?: string;
};




export interface Page {
  id: UUID;
  slug: string;
  title: string;
  cover_image_id: UUID | null;
  seo_title?: string;
  seo_description?: string;
  seo_canonical?: string;
  content: string;
  is_draft?: boolean;
  is_published?: boolean;
  created_at: string;
  updated_at: string;
}



export interface BaseBlock {
  id?: string;
  sort_order: number;
}


export type HeadingBlock = {
  id: string;
  type: 'heading';
  sort_order: number;
  props: {
    title: string;
    description: string;
  };
};

export type RichTextBlock = {
  id: string;
  type: 'richtext';
  sort_order: number;
  props: {
    html: string;
  };
};

export type ImageBlock = {
  id: string;
  type: 'image';
  sort_order: number;
  props: {
    media_id: string;
    alt_text?: string;
    align: 'left' | 'center' | 'right';
    object_fit?: 'cover' | 'contain' | 'fill' | 'scale-down' | 'none';
  };
};

export type Block = HeadingBlock | RichTextBlock | ImageBlock;

