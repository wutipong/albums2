import type { RequestHandler } from './$types';
import redis from '#lib/server/redis.js';

export const POST: RequestHandler = async ({ params }) => {
	const count = await redis.RPUSH(
		'tasks',
		JSON.stringify({
			type: 'media.embedding',
			payload: {
				id: params.id
			}
		})
	);

	return Response.json({
		success: true,
		count: count
	});
};
