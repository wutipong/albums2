import { Media } from '#lib/server/database.js';
import type { RequestHandler } from './$types';
import redis from '#lib/server/redis.js';

export const POST: RequestHandler = async ({ params }) => {
	const media = await Media.find({ collectionId: params.collection_id });
	const tasks = [];
	for (const m of media) {
		tasks.push(
			JSON.stringify({
				type: 'media.embedding',
				payload: {
					id: m._id
				}
			})
		);
	}

	const count = await redis.RPUSH('tasks', tasks);

	return Response.json({
		success: true,
		count: count
	});
};
