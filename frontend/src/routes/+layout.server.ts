import type { PageServerLoad, PageServerLoadEvent } from './$types';

export const load: PageServerLoad = async ({ locals }: PageServerLoadEvent) => {
	return { session: locals.session, user: locals.user };
};

export const ssr = false;
