import { createClient } from 'redis';
import { env } from '$env/dynamic/private';

const client = createClient({
	url: env.REDIS_URL
});

export default client;
