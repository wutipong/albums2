import { json, type Handle, type HandleServerError, type ServerInit } from '@sveltejs/kit';
import { building } from '$app/environment';
import { auth } from '$lib/server/auth';
import { svelteKitHandler } from 'better-auth/svelte-kit';
import mongoose from 'mongoose';
import { env } from '$env/dynamic/private';
import { sequence } from '@sveltejs/kit/hooks';
import { logger } from 'better-auth';

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

const handleError: Handle = async ({ event, resolve }) => {
	// Only target routes under /api
	if (event.url.pathname.startsWith('/api')) {
		try {
			const response = await resolve(event);

			if (response.status >= 400) {
				// Attempt to extract the existing error message if one was provided
				let message = `API Error: Received status ${response.status}`;
				try {
					const contentType = response.headers.get('content-type');
					if (contentType?.includes('application/json')) {
						const body = await response.json();
						message = body.error || message;
					} else {
						message = (await response.text()) || message;
					}
				} catch (err){
					// Fallback if the body can't be read or parsed
					logger.error('Error occurred while parsing API error response:', err);
				}

				return json(
					{
						message,
						status: response.status
					},
					{ status: response.status }
				);
			}

			return response;
		} catch (error) {
			// Intercept any unhandled/unexpected errors thrown inside your +server.ts files
			console.error('API Error Intercepted:', error);

			return json(
				{
					message: error instanceof Error ? error.message : 'An unexpected error occurred',
					status: 500
				},
				{ status: 500 }
			);
		}
	}

	// Fallback behavior for all standard non-API frontend pages
	return await resolve(event);
};

export const handle: Handle = sequence(handleBetterAuth, handleError);
