import { error } from '@sveltejs/kit';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async ({ parent }) => {
	const { user } = await parent();

	if (!user || user.role !== 'admin') {
		throw error(404, 'Not Found');
	}

	return {};
};
