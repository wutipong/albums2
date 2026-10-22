import pino from 'pino';
import { dev } from '$app/environment';

export default pino({
	level: dev ? 'trace' : 'info',
	transport: dev ? { target: 'pino-pretty' } : undefined
});
