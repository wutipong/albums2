import type { RequestHandler } from '@sveltejs/kit';
import { json } from '@sveltejs/kit';
import { Collection } from '$lib/server/database';
import log from '$lib/log';

export const GET: RequestHandler = async ({ params }) => {
	log.debug(params, 'GET /api/collections');
	const collection = await Collection.findById(params.collection_id)

	return json( collection );
};
