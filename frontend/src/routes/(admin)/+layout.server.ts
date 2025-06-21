
// +layout.server.ts
export const load = ({ locals }) => {
	return {
		user: locals.user
	};
};
