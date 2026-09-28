import { env } from '$env/dynamic/private';
import { betterAuth } from 'better-auth/minimal';
import { sveltekitCookies } from 'better-auth/svelte-kit';
import { getRequestEvent } from '$app/server';
import { MongoClient } from 'mongodb';
import { mongodbAdapter } from 'better-auth/adapters/mongodb';
import { genericOAuth } from 'better-auth/plugins';

const client = new MongoClient(env.DB_CONNECTION);
const db = client.db();

export const auth = betterAuth({
	baseURL: env.ORIGIN,
	secret: env.BETTER_AUTH_SECRET,
	database: mongodbAdapter(db, {
		client
	}),

	plugins: [
		genericOAuth({
			config: [
				{
					providerId: env.OIDC_PROVIDER_ID,
					clientId: env.OIDC_CLIENT_ID,
					clientSecret: env.OIDC_SECRET,
					discoveryUrl: env.OIDC_DISCOVERY_URL
				}
			]
		}),
		sveltekitCookies(getRequestEvent) // make sure this is the last plugin in the array
	]
});
