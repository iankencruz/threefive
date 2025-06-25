
export type User = {
	id: number;
	first_name: string;
	last_name: string;
	email: string;
	roles: string[];
}


export type Project = {
	id: string;
	title: string;
	slug: string;
	description: string | null;
	meta_description: string | null;
	canonical_url: string | null;
	cover_media_id: string | null;
	is_published: boolean;
	published_at: string | null;
	created_at: string;
	updated_at: string;
};



export type MediaItem = {
	id: string;
	title: string;
	url: string;
	thumbnail_url?: string;
	alt_text?: string;
	file_size?: string;
};

