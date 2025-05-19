const BASE_URL = import.meta.env.VITE_API_BASE_URL;

export async function login(username: string, password: string): Promise<Response> {
	return fetch(`${BASE_URL}/login`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ username, password }),
		credentials: 'include'
	});
}

// TODO on backend
export async function logout() {
	await fetch(`${BASE_URL}/logout`, {
		method: 'POST',
		credentials: 'include'
	});
}
