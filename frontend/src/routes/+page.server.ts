import { redirect } from '@sveltejs/kit';
import type { PageServerLoad, PageServerLoadEvent } from './$types';

export const load: PageServerLoad = ({ locals }: PageServerLoadEvent) => {
	if (locals.session) {
		redirect(302, '/collections');
	}

	redirect(302, '/login');
};
