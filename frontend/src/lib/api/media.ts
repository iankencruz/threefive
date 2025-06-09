
export async function fetchMedia(page = 1, limit = 20) {
	const res = await fetch(`/api/admin/media?page=${page}&limit=${limit}`);
	if (!res.ok) throw new Error('Failed to fetch media');
	return res.json();
}

export async function updateMedia(id: string, data: { title: string; alt_text: string }) {
	const res = await fetch(`/api/admin/media/${id}`, {
		method: 'PUT',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(data)
	});
	if (!res.ok) throw new Error('Failed to update media');
	return res.json();
}

export async function deleteMedia(id: string) {
	const res = await fetch(`/api/admin/media/${id}`, {
		method: 'DELETE'
	});
	if (!res.ok) throw new Error('Failed to delete media');
}
