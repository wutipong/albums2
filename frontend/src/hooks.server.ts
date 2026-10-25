import { error, json } from '@sveltejs/kit';
import { building } from '$app/env';
import { auth } from '#lib/server/auth.js';
import { svelteKitHandler } from 'better-auth/svelte-kit';
import mongoose from 'mongoose';
import { DB_CONNECTION } from '$app/env/private';

import {
	sequence,
	type Handle,
	type HandleServerError,
	type ServerInit
} from '@sveltejs/kit/hooks';

import log from '#lib/log.js';
import redis from '#lib/server/redis.js';
import { Media } from '#lib/server/database.js';

export const init: ServerInit = async () => {
	await mongoose.connect(DB_CONNECTION ?? 'mongodb://localhost');
	await Media.init();
	await Media.createSearchIndexes();
	if (!redis.isReady) {
		await redis.connect();
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

const handleSession: Handle = async ({ event, resolve }) => {
	if (event.url.pathname == '/login' || event.url.pathname == '/') {
		return resolve(event);
	}

	const session = event.locals.session;
	if (!session) {
		error(403, 'forbidden.');
	}

	return resolve(event);
};

const handleAdmin: Handle = async ({ event, resolve }) => {
	if (!event.url.pathname.startsWith('/admin')) {
		return resolve(event);
	}

	const user = event.locals.user;
	if (!user || user.role !== 'admin') {
		error(403, 'forbidden.');
	}

	return resolve(event);
};

export const handle: Handle = sequence(handleBetterAuth, handleAdmin, handleSession);
