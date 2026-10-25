import type { RequestHandler } from '@sveltejs/kit';
import { Collection } from '#lib/server/database.js';
import log from '#lib/log.js';

export const GET: RequestHandler = async ({ request, params }) => {
	log.debug(params, 'GET /api/collections/by-name');

	const name = new URL(request.url).searchParams.get('name');
	if (!name) {
		return Response.json({ message: 'Missing name query parameter' }, { status: 400 });
	}

	const collection = await Collection.findOne({ name });
	if (!collection) {
		return Response.json({ existed: false });
	}

	return Response.json({ existed: true, collection });
};
