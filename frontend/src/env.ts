import { defineEnvVars } from '@sveltejs/kit/env';

// @migration-task Review usage of dynamic environment variables. They fall back to the empty string if not present, which may not be what you want.
export const variables = defineEnvVars({
	DB_CONNECTION: { schema: (input) => input ?? '' },
	OIDC_PROVIDER_ID: { schema: (input) => input ?? '' },
	IMGPROXY_URL: { schema: (input) => input ?? '' },
	S3_BUCKET: { schema: (input) => input ?? '' },
	IMGPROXY_SALT: { schema: (input) => input ?? '' },
	IMGPROXY_KEY: { schema: (input) => input ?? '' },
	ORIGIN: { schema: (input) => input ?? '' },
	BETTER_AUTH_SECRET: { schema: (input) => input ?? '' },
	OIDC_CLIENT_ID: { schema: (input) => input ?? '' },
	OIDC_SECRET: { schema: (input) => input ?? '' },
	OIDC_DISCOVERY_URL: { schema: (input) => input ?? '' },
	REDIS_URL: { schema: (input) => input ?? '' },
	AWS_DEFAULT_REGION: { schema: (input) => input ?? '' },
	AWS_ENDPOINT_URL: { schema: (input) => input ?? '' },
	AWS_ACCESS_KEY_ID: { schema: (input) => input ?? '' },
	AWS_SECRET_ACCESS_KEY: { schema: (input) => input ?? '' },
	MLSERVICE_ADDRESS: { schema: (input) => input ?? '' }
});
