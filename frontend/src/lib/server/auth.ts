import {
	ORIGIN,
	BETTER_AUTH_SECRET,
	OIDC_PROVIDER_ID,
	OIDC_CLIENT_ID,
	OIDC_SECRET,
	OIDC_DISCOVERY_URL
} from '$app/env/private';

import { betterAuth } from 'better-auth/minimal';
import { sveltekitCookies } from 'better-auth/svelte-kit';
import { getRequestEvent } from '$app/server';
import { mongodbAdapter } from 'better-auth/adapters/mongodb';
import { admin, genericOAuth } from 'better-auth/plugins';
import { apiKey } from '@better-auth/api-key';
import { client } from '#lib/server/database.js';

const db = client.db();

export const auth = betterAuth({
	baseURL: ORIGIN,
	secret: BETTER_AUTH_SECRET,
	database: mongodbAdapter(db, { client }),
	plugins: [
		admin(),
		apiKey({
			rateLimit: { enabled: false },
			enableSessionForAPIKeys: true
		}),
		genericOAuth({
			config: [
				{
					providerId: OIDC_PROVIDER_ID || 'placeholder-provider',
					clientId: OIDC_CLIENT_ID || 'placeholder-client',
					clientSecret: OIDC_SECRET || 'placeholder-secret',
					discoveryUrl: OIDC_DISCOVERY_URL || 'http://placehodler-config',
					scopes: ['openid', 'email', 'profile']
				}
			]
		}),
		sveltekitCookies(getRequestEvent) // make sure this is the last plugin in the array
	]
});
