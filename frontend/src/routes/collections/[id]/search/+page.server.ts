import type { PageServerLoad } from './$types';
import { encodeText } from '#lib/server/grpc/mlservice.js';
import { error, redirect } from '@sveltejs/kit';
import log from '#lib/log.js';
import { searchMedia } from '#lib/server/collection.js';

export const ssr = false;

export const load: PageServerLoad = async ({ params, url }) => {
	const search = url.searchParams.get('search');
	if (!search) {
		error(400, 'search parameter is required.');
	}

	log.debug(`search: ${search}`);

	const embedding = await encodeText(search);
	const media = await searchMedia(params.id, embedding.embedding, 1);

	return {
		collectionId: params.id,
		search: search,
		media: media
	};
};

export const actions = {
	search: async ({ request }) => {
		const data = await request.formData();
		const text = data.get('text');
		const collectionId = data.get('collectionId');

		const query = new URLSearchParams({ search: text?.toString() ?? '' });

		redirect(303, `/collections/${collectionId}/search?${query.toString()}`);
	}
};
