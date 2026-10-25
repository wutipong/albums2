import pino from 'pino';
import { dev } from '$app/env';

export default pino({
	level: dev ? 'trace' : 'info',
	transport: dev ? { target: 'pino-pretty' } : undefined
});
