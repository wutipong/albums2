import type { RequestHandler } from '@sveltejs/kit';
import { json } from '@sveltejs/kit';
import { Collection } from '$lib/server/database';
import log from '$lib/log';
import mongoose from 'mongoose';

export const GET: RequestHandler = async ({ params }) => {
	log.debug(params, 'GET /api/collections');

	if (!params.collection_id) {
		return json({ error: 'Missing collection ID' }, { status: 400 });
	}

	if (mongoose.Types.ObjectId.isValid(params.collection_id) === false) {
		return json({ error: 'Invalid collection ID' }, { status: 400 });
	}
	const collection = await Collection.findById(params.collection_id)

	if (!collection) {
		return json({ error: 'Collection not found' }, { status: 404 });
	}

	return json( collection );
};
