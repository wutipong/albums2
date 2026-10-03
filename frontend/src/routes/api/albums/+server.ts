import type { RequestHandler } from '@sveltejs/kit';
import { json, error } from '@sveltejs/kit';
import { Album } from '$lib/server/database';

export const POST: RequestHandler = async ({ request }) => {
	const body = await request.json();
	const { name } = body;

	if (!name || typeof name !== 'string') {
		return error(400, { message: 'Missing or invalid name' });
	}

	const album = await Album.create({ name });

	return json(album, { status: 201 });
};

export const GET: RequestHandler = async () => {
	const albums = await Album.find({ deletedAt: null }).sort({ name: 1 }).lean();
    
	return json({ albums });
};