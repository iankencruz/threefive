import type { MediaItem } from "../types";


export async function fetchMedia(page = 1, limit = 20) {
	const res = await fetch(`/api/v1/admin/media?page=${page}&limit=${limit}`);
	if (!res.ok) throw new Error('Failed to fetch media');
	const result = await res.json();
	return result.data;

}


export async function updateMedia(id: string, payload: { title: string; alt_text: string }) {
	const res = await fetch(`/api/v1/admin/media/${id}`, {
		method: 'PUT',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({
			title: payload.title,
			alt: payload.alt_text // 👈 Must be `alt` to match the Go handler
		})
	});

	if (!res.ok) throw new Error('Failed to update media');
}


export async function deleteMedia(id: string) {
	const res = await fetch(`/api/v1/admin/media/${id}`, {
		method: 'DELETE'
	});
	if (!res.ok) throw new Error('Failed to delete media');
}


export async function getMediaById(id: string): Promise<MediaItem> {
	const res = await fetch(`/api/v1/admin/media/${id}`);
	if (!res.ok) throw new Error('Failed to fetch media');
	const json = await res.json();
	return json.data
}




export async function uploadMedia(
	file: File,
	onProgress?: (percent: number) => void
): Promise<any> {
	return new Promise((resolve, reject) => {
		const formData = new FormData();
		formData.append('file', file);
		formData.append('title', file.name); // Optional: add more fields like alt text or sort order if needed

		const xhr = new XMLHttpRequest();
		xhr.open('POST', '/api/v1/admin/media', true);

		xhr.upload.onprogress = (event) => {
			if (event.lengthComputable && onProgress) {
				const percent = Math.round((event.loaded / event.total) * 100);
				onProgress(percent);
			}
		};

		xhr.onload = () => {
			if (xhr.status >= 200 && xhr.status < 300) {
				try {
					const response = JSON.parse(xhr.responseText);
					resolve(response);
				} catch (e) {
					reject(new Error('Failed to parse response'));
				}
			} else {
				reject(new Error(`Upload failed: ${xhr.statusText}`));
			}
		};

		xhr.onerror = () => reject(new Error('Network error during upload'));
		xhr.send(formData);
	});
}


export async function uploadMediaAndLinkToContext(file: File, context: { type: 'project' | 'block'; id: string }): Promise<MediaItem> {
	const formData = new FormData();
	formData.append('file', file);

	const res = await fetch(`/api/v1/admin/media/upload?type=${context.type}&id=${context.id}`, {
		method: 'POST',
		body: formData
	});

	if (!res.ok) throw new Error('Upload failed');

	const json = await res.json();
	return json.data as MediaItem;
}

