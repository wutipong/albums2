import { authClient } from '$lib/auth';
import { redirect } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
	authClient.signOut();

	redirect(308, '/login');
};