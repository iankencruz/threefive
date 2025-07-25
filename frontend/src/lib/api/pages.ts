import type { Block } from "$lib/types";

export async function saveToBackend(content: string) {
  try {
    const res = await fetch('/api/v1/pages/page-id', {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ content }),
    });

    if (!res.ok) {
      throw new Error('Failed to save content');
    }

    console.log('Saved successfully');
  } catch (err) {
    console.error(err);
  }
}



export async function sortBlocks(slug: string, blocks: Block[]) {
  const res = await fetch(`/api/v1/admin/pages/${slug}/blocks/sort`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(
      blocks.map((b, i) => ({
        id: b.id,
        sort_order: i
      }))
    )
  });

  if (!res.ok) {
    throw new Error('Failed to sort blocks');
  }
}
