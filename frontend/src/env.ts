import { defineEnvVars } from '@sveltejs/kit/env';

export const variables = defineEnvVars({
	AWS_ACCESS_KEY_ID: { schema: (input) => input ?? '<placeholder>' },
	AWS_DEFAULT_REGION: { schema: (input) => input ?? '<placeholder>' },
	AWS_ENDPOINT_URL: { schema: (input) => input ?? 'http://placeholder' },
	AWS_SECRET_ACCESS_KEY: { schema: (input) => input ?? '<placeholder>' },
	BETTER_AUTH_SECRET: { schema: (input) => input ?? '' },
	DB_CONNECTION: { schema: (input) => input ?? 'mongodb://localhost:5173' },
	IMGPROXY_KEY: { schema: (input) => input ?? '' },
	IMGPROXY_SALT: { schema: (input) => input ?? '' },
	IMGPROXY_URL: { schema: (input) => input ?? '' },
	MLSERVICE_ADDRESS: { schema: (input) => input ?? 'http://localhost:8173' },
	OIDC_CLIENT_ID: { schema: (input) => input ?? '<placeholder client id>' },
	OIDC_DISCOVERY_URL: { schema: (input) => input ?? 'http://placeholder' },
	OIDC_PROVIDER_ID: { schema: (input) => input ?? '<placeholder>' },
	OIDC_SECRET: { schema: (input) => input ?? '<placeholder>' },
	ORIGIN: { schema: (input) => input ?? '' },
	REDIS_URL: { schema: (input) => input ?? '' },
	S3_BUCKET: { schema: (input) => input ?? '' }
});
