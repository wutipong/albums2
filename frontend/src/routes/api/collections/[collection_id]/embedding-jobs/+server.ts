import { Media } from '$lib/server/database';
import { json } from '@sveltejs/kit';
import { type RequestHandler } from './$types';
import redis from '$lib/server/redis';

export const POST: RequestHandler = async ({ params }) => {
	let media = await Media.find({ collectionId: params.collection_id });
	const tasks = [];
	for (const m of media) {
		tasks.push(
			JSON.stringify({
				type: 'embedding',
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
