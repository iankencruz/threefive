
export async function uploadMedia({
	file,
	title = '',
	alt = '',
	sort = 0
}: {
	file: File;
	title?: string;
	alt?: string;
	sort?: number;
}) {
	const formData = new FormData();
	formData.append('file', file);
	formData.append('title', title);
	formData.append('alt', alt);
	formData.append('sort', String(sort));

	const res = await fetch('/api/admin/media/upload', {
		method: 'POST',
		body: formData
	});

	if (!res.ok) {
		const error = await res.json().catch(() => ({}));
		throw new Error(error.message || 'Upload failed');
	}

	return res.json();
}
