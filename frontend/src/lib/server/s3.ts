import { S3Client } from 'bun';

import {
	AWS_DEFAULT_REGION,
	AWS_ENDPOINT_URL,
	AWS_ACCESS_KEY_ID,
	AWS_SECRET_ACCESS_KEY
} from '$app/env/private';

export const s3 = new S3Client({
	region: AWS_DEFAULT_REGION || 'dummy-region',
	endpoint: AWS_ENDPOINT_URL || 'http://localhost',
	accessKeyId: AWS_ACCESS_KEY_ID || 'dummy-key',
	secretAccessKey: AWS_SECRET_ACCESS_KEY || 'dummy-secret'
});
