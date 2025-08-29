import type { Gallery } from "$lib/types";

export async function getGalleries(): Promise<any[]> {
  const res = await fetch('/api/v1/admin/galleries');

  if (!res.ok) {
    throw new Error('Failed to load projects');
  }

  let result = await res.json();
  return result.data
}


export async function createGallery({
  title,
  slug,
  description
}: {
  title: string;
  slug: string;
  description: string;
}) {
  const res = await fetch('/api/v1/admin/galleries', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ title, slug, description })
  });

  const data = await res.json();

  if (!res.ok) {
    throw new Error(data?.message || 'Failed to create gallery');
  }

  return data;
}


// Get a single gallery by slug
export async function getGalleryBySlug(slug: string): Promise<any> {
  const res = await fetch(`/api/v1/admin/galleries/${slug}`);
  if (!res.ok) throw new Error('Failed to load gallery');
  const result = await res.json();
  return result.data;
}





export async function updateGalleryBySlug(slug: string, payload: Partial<Gallery>): Promise<Gallery> {
  const res = await fetch(`/api/v1/admin/galleries/${slug}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload)
  });
  if (!res.ok) throw new Error('Failed to update gallery');
  const result = await res.json();
  return result.data;
}



export async function linkMediaToGallery(slug: string, mediaId: string): Promise<void> {
  const res = await fetch(`/api/v1/admin/galleries/${slug}/media`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ media_id: mediaId, sort_order: 0 })
  });

  if (!res.ok) throw new Error('Failed to link media to project');
}



export async function unlinkMediaFromGallery(slug: string, mediaId: string): Promise<void> {
  const res = await fetch(`/api/v1/admin/galleries/${slug}/media`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ media_id: mediaId })
  });

  if (!res.ok) throw new Error('Failed to unlink media from project');
}
