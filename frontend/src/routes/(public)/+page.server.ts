
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
	const res = await fetch('http://localhost:8080/api/v1/projects'); // <-- use actual backend API address
	const json = await res.json();

	return {
		projects: json.data
	};
};

