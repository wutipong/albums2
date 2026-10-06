import type { RequestHandler } from '@sveltejs/kit';
import { error, json } from '@sveltejs/kit';
import { Collection } from '$lib/server/database';
import log from '$lib/log';
import mongoose from 'mongoose';

export const GET: RequestHandler = async ({ params }) => {
	log.debug(params, 'GET /api/collections');

	if (!params.collection_id) {
		return error(400, 'missing collection id');
	}

	if (mongoose.Types.ObjectId.isValid(params.collection_id) === false) {
		return error(400, 'invalid collection id');
	}
	const collection = await Collection.findById(params.collection_id);

	if (!collection) {
		return error(404, 'collection not found');
	}

	return json(collection);
};
