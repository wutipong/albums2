import type { PageServerLoad, PageServerLoadEvent } from './$types';
import redis from '#lib/server/redis.ts';
import log from '#lib/log.ts';

export const load: PageServerLoad = async ({ params }: PageServerLoadEvent) => {
	const taskStrs = await redis.lRange('tasks', 1, 50);

	const tasks = [];
	for (const str of taskStrs) {
		tasks.push(JSON.parse(str));
	}

	const taskCount = await redis.lLen('tasks');

	const processingKeys = [];
	for await (const key of redis.scanIterator({
		MATCH: 'prefix:processing_',
		COUNT: 100 // Grab 100 keys per iteration batch
	})) {
		processingKeys.push(...key);
	}

	log.debug(processingKeys, 'processing list');

	const processings = [];
	for (const k of processingKeys) {
		const processingTaskStrs = await redis.lRange(k, 1, -1);
		for (const str of processingTaskStrs) {
			processings.push(JSON.parse(str));
		}
	}

	const failures = [];
	const failStrs = await redis.lRange('errors', 1, 50);
	for (const str of failStrs) {
		failures.push(JSON.parse(str));
	}

	return { tasks, taskCount, processings, failures };
};
