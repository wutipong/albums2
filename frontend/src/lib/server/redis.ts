import { createClient } from 'redis';
import { env } from '$env/dynamic/private';
import log from '$lib/log'
const client = createClient({
	url: env.REDIS_URL
});

log.info({url: env.REDIS_URL, redis: client}, "using Redis")

export default client;
