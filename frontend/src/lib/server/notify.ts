import redis from '$lib/server/redis';

export async function notifyProcessMedia(id: String) {
	const taskRequest = {
		type: 'media',
		payload: {
			media: id
		}
	};

	redis.LPUSH('tasks', JSON.stringify(taskRequest));
}
