
// src/lib/server/auth.ts

export interface User {
	id: string;
	name: string;
	email: string;
	role: string;
	token: string;
}

const fakeUsers: User[] = [
	{
		id: '1',
		name: 'Admin User',
		email: 'admin@example.com',
		role: 'admin',
		token: 'admin-token-123'
	},
	{
		id: '2',
		name: 'Regular User',
		email: 'user@example.com',
		role: 'user',
		token: 'user-token-456'
	}
];

/**
 * Simulates checking user credentials.
 */
export async function authenticateUser(email: string, password: string): Promise<User | null> {
	// For now, we ignore the password. You can improve this with bcrypt later.
	const user = fakeUsers.find((u) => u.email === email);
	return user ?? null;
}



export async function getUserBySession(token: string): Promise<User | null> {
	const user = fakeUsers.find((u) => u.token === token);
	return user ?? null;
}

