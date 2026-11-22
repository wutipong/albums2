import type { RequestHandler } from './$types';
import redis from '#lib/server/redis.js';
import log from '#lib/log.js';
import { db } from '#lib/server/database.ts';
import { ObjectId } from 'mongodb';

export const POST: RequestHandler = async ({ params }) => {
	log.debug(params, "populate collection's albums cover");

	const albums = await db
		.collection('albums')
		.find({
			collectionId: new ObjectId(params.collection_id)
		})
		.toArray();

	const tasks = albums.map((album) =>
		JSON.stringify({
			type: 'album.cover',
			payload: {
				id: album._id.toString()
			}
		})
	);

	const count = await redis.RPUSH('tasks', tasks);

	return Response.json({
		success: true,
		count: count
	});
};
