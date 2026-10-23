import { json } from '@sveltejs/kit';
import { type RequestHandler } from './$types';
import redis from '$lib/server/redis';

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

	return json({
		success: true,
		count: count
	});
};
