import { redirect } from '@sveltejs/kit';
import { type PageServerLoad } from './$types';

export const load: PageServerLoad = ({ params }) => {
	redirect(308, `/collections/${params.collectionId}/1`);
};
