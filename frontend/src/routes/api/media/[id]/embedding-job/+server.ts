import type { RequestHandler } from './$types';
import redis from '#lib/server/redis.js';

export const POST: RequestHandler = async ({ params }) => {
	const count = await redis.LPUSH(
		'tasks',
		JSON.stringify({
			type: 'embedding',
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
