import type { Project } from "../types";

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
export async function getProjectBySlug(slug: string): Promise<any> {
	const res = await fetch(`/api/v1/admin/projects/${slug}`);
	if (!res.ok) throw new Error('Failed to load project');
	const result = await res.json();
	return result.data;
}



export async function deleteProjectBySlug(slug: string): Promise<void> {
	const res = await fetch(`/api/v1/admin/projects/${slug}`, {
		method: 'DELETE'
	});

	if (!res.ok) {
		throw new Error('Failed to delete project');
	}
}


export async function updateProjectBySlug(slug: string, payload: Partial<Project>): Promise<Project> {
	const res = await fetch(`/api/v1/admin/projects/${slug}`, {
		method: 'PUT',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(payload)
	});
	if (!res.ok) throw new Error('Failed to update project');
	const result = await res.json();
	return result.data;
}
// Link media to a project
export async function linkMediaToProject(slug: string, mediaId: string): Promise<void> {
	const res = await fetch(`/api/v1/admin/projects/${slug}/media`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ media_id: mediaId, sort_order: 0 })
	});

	if (!res.ok) throw new Error('Failed to link media to project');
}

// Unlink media from a project
export async function unlinkMediaFromProject(slug: string, mediaId: string): Promise<void> {
	const res = await fetch(`/api/v1/admin/projects/${slug}/media`, {
		method: 'DELETE',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ media_id: mediaId })
	});

	if (!res.ok) throw new Error('Failed to unlink media from project');
}



export async function updateProjectMediaOrder(slug: string, mediaIds: string[]): Promise<void> {
	const res = await fetch(`/api/v1/admin/projects/${slug}/media/sort`, {
		method: 'PUT',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ media_ids: mediaIds })
	});
	if (!res.ok) {
		const error = await res.text();
		throw new Error(`Failed to update media order: ${error}`);

	}

	return await res.json();
}

