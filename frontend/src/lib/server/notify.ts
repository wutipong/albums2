import redis from '#lib/server/redis.js';
import log from '#lib/log.js';

export async function notifyProcessMedia(id: string) {
	log.debug(redis, 'redis');
	try {
		const taskRequest = {
			type: 'media',
			payload: {
				id: id
			}
		};

		redis.LPUSH('tasks', JSON.stringify(taskRequest));
	} catch (e) {
		log.error(e, 'send command fails');
	}
}
