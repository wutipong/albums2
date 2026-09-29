import { env } from '$env/dynamic/private';
import { betterAuth } from 'better-auth/minimal';
import { sveltekitCookies } from 'better-auth/svelte-kit';
import { getRequestEvent } from '$app/server';
import { mongodbAdapter } from 'better-auth/adapters/mongodb';
import { genericOAuth } from 'better-auth/plugins';
import { client } from '$lib/server/database';

const db = client.db();

export const auth = betterAuth({
	baseURL: env.ORIGIN,
	secret: env.BETTER_AUTH_SECRET,
	database: mongodbAdapter(db, {
		client: client
	}),

	plugins: [
		genericOAuth({
			config: [
				{
					providerId: env.OIDC_PROVIDER_ID || 'placeholder-provider',
					clientId: env.OIDC_CLIENT_ID || 'placeholder-client',
					clientSecret: env.OIDC_SECRET || 'placeholder-secret',
					discoveryUrl: env.OIDC_DISCOVERY_URL || 'http://placehodler-config',
					scopes: ['openid', 'email', 'profile']
				}
			]
		}),
		sveltekitCookies(getRequestEvent) // make sure this is the last plugin in the array
	]
});
