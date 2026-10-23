import { type RequestHandler } from './$types';
import redis from '$lib/server/redis';
import { Media } from '$lib/server/database';
import { json } from '@sveltejs/kit';
import log from '$lib/log';

export const POST: RequestHandler = async ({ request, params }) => {
	log.debug(request.url);
	const media = await Media.find({ collectionId: params.collection_id, deletedAt: null });

	const tasks = [];
	for (const m of media) {
		tasks.push(
			JSON.stringify({
				type: 'media',
				payload: {
					id: m._id
				}
			})
		);
	}
	const count = await redis.LPUSH('tasks', tasks);

	return json({
		success: true,
		count: count
	});
};
