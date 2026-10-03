import type { RequestHandler } from '@sveltejs/kit';
import { json } from '@sveltejs/kit';
import { Collection } from '$lib/server/database';
import log from '$lib/log';

export const GET: RequestHandler = async ({ request, params }) => {
	log.debug(params, 'GET /api/collections/by-name');

	const name = new URL(request.url).searchParams.get('name');
	if (!name) {
		return json({ message: 'Missing name query parameter' }, { status: 400 });
	}

	const collection = await Collection.findOne({ name }).lean();
	if (!collection) {
		return json({ existed: false });
	}

	return json({ existed: true, collection });
};
