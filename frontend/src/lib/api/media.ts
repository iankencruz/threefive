
export async function fetchMedia(page = 1, limit = 20) {
	const res = await fetch(`/api/v1/admin/media?page=${page}&limit=${limit}`);
	if (!res.ok) throw new Error('Failed to fetch media');
	return res.json();
}

export async function updateMedia(id: string, data: { title: string; alt_text: string }) {
	const res = await fetch(`/api/v1/admin/media/${id}`, {
		method: 'PUT',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(data)
	});
	if (!res.ok) throw new Error('Failed to update media');
	return res.json();
}

export async function deleteMedia(id: string) {
	const res = await fetch(`/api/v1/admin/media/${id}`, {
		method: 'DELETE'
	});
	if (!res.ok) throw new Error('Failed to delete media');
}



// src/lib/api/media.ts

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

