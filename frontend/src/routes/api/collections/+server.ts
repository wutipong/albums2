import type { RequestHandler } from '@sveltejs/kit';
import { error } from '@sveltejs/kit';
import { Collection, Album } from '#lib/server/database.js';
import log from '#lib/log.js';

export const POST: RequestHandler = async ({ request }) => {
	const body = await request.json();
	log.debug(body, 'POST /api/collections');

	const { name } = body;

	if (!name || typeof name !== 'string') {
		return error(400, 'Missing or invalid name');
	}

	const collection = await Collection.create({ name });

	return Response.json(collection, { status: 201 });
};

export const GET: RequestHandler = async () => {
	log.debug('GET /api/collections');
	const collections = await Collection.find({ deletedAt: null }).sort({ name: 1 });

	log.debug(collections, 'GET /api/collections result');

	return Response.json({ collections });
};
