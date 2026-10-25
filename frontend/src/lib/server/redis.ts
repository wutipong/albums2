import { createClient } from 'redis';
import { REDIS_URL } from '$app/env/private';
import log from '#lib/log.js';

const client = createClient({ url: REDIS_URL });

log.info({ url: REDIS_URL, redis: client }, 'using Redis');

export default client;
