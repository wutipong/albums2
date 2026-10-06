import { json, type Handle, type HandleServerError, type ServerInit } from '@sveltejs/kit';
import { building } from '$app/environment';
import { auth } from '$lib/server/auth';
import { svelteKitHandler } from 'better-auth/svelte-kit';
import mongoose from 'mongoose';
import { env } from '$env/dynamic/private';
import { sequence } from '@sveltejs/kit/hooks';
import log from '$lib/log';

export const init: ServerInit = async () => {
	await mongoose.connect(env.DB_CONNECTION);
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
