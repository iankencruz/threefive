
// src/lib/api/project.ts
export async function createProject({
	title,
	slug,
	description
}: {
	title: string;
	slug: string;
	description: string;
}) {
	const res = await fetch('/api/v1/admin/projects', {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ title, slug, description })
	});

	const data = await res.json();

	if (!res.ok) {
		throw new Error(data?.message || 'Failed to create project');
	}

	return data;
}

// Get all projects
export async function getProjects(): Promise<any[]> {
	const res = await fetch('/api/v1/admin/projects');

	if (!res.ok) {
		throw new Error('Failed to load projects');
	}

	let result = await res.json();
	return result.data
}


// Get a single project by ID
export async function getProjectById(id: string): Promise<any> {
	const res = await fetch(`/api/v1/admin/projects/${id}`);
	if (!res.ok) throw new Error('Failed to load project');
	const result = await res.json();
	return result.data;
}



export async function deleteProjectByID(id: string): Promise<void> {
	const res = await fetch(`/api/v1/admin/projects/${id}`, {
		method: 'DELETE'
	});

	if (!res.ok) {
		throw new Error('Failed to delete project');
	}
}
