import { S3Client } from 'bun';
import { env } from '$env/dynamic/private';

export const s3 = new S3Client({
	region: env.AWS_DEFAULT_REGION || 'dummy-region',
	endpoint: env.AWS_ENDPOINT_URL || 'http://localhost',
	accessKeyId: env.AWS_ACCESS_KEY_ID || 'dummy-key',
	secretAccessKey: env.AWS_SECRET_ACCESS_KEY || 'dummy-secret'
});
