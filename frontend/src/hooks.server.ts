import { json, type Handle, type HandleServerError, type ServerInit } from '@sveltejs/kit';
import { building } from '$app/environment';
import { auth } from '$lib/server/auth';
import { svelteKitHandler } from 'better-auth/svelte-kit';
import mongoose from 'mongoose';
import { env } from '$env/dynamic/private';
import { sequence } from '@sveltejs/kit/hooks';
import log from '$lib/log';
import redis from '$lib/server/redis'
import { Media } from '$lib/server/database';

export const init: ServerInit = async () => {
	await mongoose.connect(env.DB_CONNECTION);
	await Media.init();
	await Media.createSearchIndexes();
	if (!redis.isReady) {
		await redis.connect()
	}
};

const handleBetterAuth: Handle = async ({ event, resolve }) => {
	const session = await auth.api.getSession({ headers: event.request.headers });

	if (session) {
		event.locals.session = session.session;
		event.locals.user = session.user;
	}

	return svelteKitHandler({ event, resolve, auth, building });
};

export const handle: Handle = sequence(handleBetterAuth /*, handleError*/);
