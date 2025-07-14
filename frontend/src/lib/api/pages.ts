
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
