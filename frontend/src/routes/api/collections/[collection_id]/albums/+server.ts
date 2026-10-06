import type { RequestHandler } from '@sveltejs/kit';
import { json, error } from '@sveltejs/kit';
import { Collection, Album } from '$lib/server/database';
import log from '$lib/log';

export const POST: RequestHandler = async ({ request, params }) => {
	const body = await request.json();
	log.debug(body, `POST ${request.url} params`);

	if (!Collection.exists({ _id: params.collection_id })) {
		return error(404, { message: 'collection not found' });
	}

	const { name } = body;

	if (!name || typeof name !== 'string') {
		return error(400, { message: 'Missing or invalid name' });
	}

	const album = await Album.create({ name, collectionId: params.collection_id });

	return json(album, { status: 201 });
};

export const GET: RequestHandler = async ({ request, params }) => {
	log.debug(params, `GET ${request.url} params`);
	if (!Collection.exists({ _id: params.collection_id })) {
		return error(404, { message: 'collection not found' });
	}

	const albums = await Album.find({ collectionId: params.collection_id, deletedAt: null })
		.sort({ name: 1 });

	return json({ albums });
};
