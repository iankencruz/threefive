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
