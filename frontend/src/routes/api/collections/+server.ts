import type { RequestHandler } from '@sveltejs/kit';
import { json, error } from '@sveltejs/kit';
import { Collection, Album } from '$lib/server/database';
import log from '$lib/log';

export const POST: RequestHandler = async ({ request }) => {
    const body = await request.json();
    log.debug(body, 'POST /api/collections');

    const { name } = body;

    if (!name || typeof name !== 'string') {
        return error(400, { message: 'Missing or invalid name' });
    }

    const collection = await Collection.create({ name });

    return json(collection, { status: 201 });
};

export const GET: RequestHandler = async () => {
    log.debug('GET /api/collections');
    const collections = await Collection.find({ deletedAt: null }).sort({ name: 1 }).lean();

    return json({ collections });
};
