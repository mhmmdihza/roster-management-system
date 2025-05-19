export interface UserPayload {
	email: string;
	employee_id: string;
	employee_name: string;
	primary_role: number;
	role: string;
}

import { writable } from 'svelte/store';

export const userStore = writable<UserPayload | null>(null);
