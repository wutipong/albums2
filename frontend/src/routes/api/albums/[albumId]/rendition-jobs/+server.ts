import type { RequestHandler } from './$types';
import redis from '#lib/server/redis.js';
import { Media } from '#lib/server/database.js';
import log from '#lib/log.js';

export const POST: RequestHandler = async ({ request, params }) => {
	log.debug(request.url);
	const media = await Media.find({ albumId: params.albumId, deletedAt: null });

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

	return Response.json({
		success: true,
		count: count
	});
};
