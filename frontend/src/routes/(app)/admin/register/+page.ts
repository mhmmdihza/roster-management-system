import { listRoles } from '$lib/api';
import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
	try {
		const roles = await listRoles();
		const filteredRoles = roles.filter((role) => role.id !== 0);
		return {
			roles: filteredRoles
		};
	} catch (error) {
		console.error('Failed to fetch roles', error);
		return {
			roles: []
		};
	}
};
