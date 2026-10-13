import redis from '$lib/server/redis';
import log from '$lib/log';

export async function notifyProcessMedia(id: String) {
	log.debug(redis, 'redis');
	try {
		const taskRequest = {
			type: 'media',
			payload: {
				media: id
			}
		};

		redis.LPUSH('tasks', JSON.stringify(taskRequest));
	} catch (e) {
		log.error(e, 'send command fails');
	}
}
