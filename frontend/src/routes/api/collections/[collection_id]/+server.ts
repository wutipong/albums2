import type { RequestHandler } from '@sveltejs/kit';
import { json, error } from '@sveltejs/kit';
import { Collection, Album } from '$lib/server/database';
import log from '$lib/log';

export const GET: RequestHandler = async ({ params }) => {
	log.debug(params, 'GET /api/collections');
	const collection = await Collection.findById(params.collection_id).lean();
	const albums = await Album.find({ collectionId: params.collection_id, deletedAt: null })
		.sort({ name: 1 })
		.lean();

	return json({ collection, albums });
};
