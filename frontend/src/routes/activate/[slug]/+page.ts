import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
	console.log(params);

	const id = params.slug;

	return {
		id
	};
};
