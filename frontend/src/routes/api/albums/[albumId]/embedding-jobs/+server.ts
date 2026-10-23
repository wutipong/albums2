import { Media } from '$lib/server/database';
import { json } from '@sveltejs/kit';
import { type RequestHandler } from './$types';
import redis from '$lib/server/redis';
import log from '$lib/log';
import mongoose from 'mongoose';

export const POST: RequestHandler = async ({ params }) => {
	const { albumId } = params;

	log.debug(`id: "${params.albumId}"`);

	const media = await Media.find({ albumId: albumId });
	log.info(media, 'media information');
	log.info(`media count: ${media.length}`);
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
