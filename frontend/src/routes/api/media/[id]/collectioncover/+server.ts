import type { RequestHandler } from './$types';
import redis from '#lib/server/redis.js';
import log from '#lib/log.js';

export const POST: RequestHandler = async ({ request, params }) => {
	log.debug(request.url);

	const count = await redis.RPUSH(
		'tasks',
		JSON.stringify({
			type: 'media.collection_cover',
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
