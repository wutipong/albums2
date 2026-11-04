import type { RequestHandler } from './$types';
import redis from '#lib/server/redis.js';
import log from '#lib/log.js';

export const POST: RequestHandler = async ({ params }) => {
	log.debug(params, 'album synchronize');

	const count = await redis.RPUSH('tasks', JSON.stringify({
		type: 'album.sync',
		payload: {
			id: params.albumId
		}
	}));

	return Response.json({
		success: true,
		count: count
	});
};
