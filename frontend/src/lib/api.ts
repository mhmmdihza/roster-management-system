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

export async function registerUser(
	email: string,
	primaryRole: number | null,
	roleAdmin: boolean
): Promise<Response> {
	const body: Record<string, unknown> = {
		email,
		roleAdmin
	};

	if (!roleAdmin) {
		body.primaryRole = primaryRole;
	}

	return fetch(`${BASE_URL}/admin/register`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(body),
		credentials: 'include'
	});
}

export interface RoleResponse {
	id: number;
	roleName: string;
}
export async function listRoles(): Promise<RoleResponse[]> {
	const response = await fetch(`${BASE_URL}/admin/list-role`, {
		method: 'GET',
		credentials: 'include'
	});

	if (!response.ok) {
		throw new Error(`Failed to fetch roles: ${response.statusText}`);
	}

	return await response.json();
}

export interface ActivateAccountRequest {
	id: string;
	name: string;
	password: string;
}

export async function activateAccount(req: ActivateAccountRequest): Promise<Response> {
	return fetch(`${BASE_URL}/activate`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(req),
		credentials: 'include'
	});
}

export interface CreateNewShiftScheduleRequest {
	roleId: number;
	startTime: string;
	endTime: string;
}

export async function createNewShiftSchedule(
	req: CreateNewShiftScheduleRequest
): Promise<Response> {
	return fetch(`${BASE_URL}/admin/schedules`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(req),
		credentials: 'include'
	});
}
