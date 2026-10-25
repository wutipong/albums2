import { IMGPROXY_URL, S3_BUCKET, IMGPROXY_SALT, IMGPROXY_KEY } from '$app/env/private';
import { Album, Collection } from '#lib/server/database.js';
import { generateImageUrl } from '@imgproxy/imgproxy-node';

export async function populateUrls(media: any): Promise<any> {
	const c = await Collection.findById(media.collectionId);
	const a = await Album.findById(media.albumId);

	let thumbnailUrl,
		previewUrl,
		viewUrl = '';
	if (media.processStatus === 'processed') {
		thumbnailUrl = generateImageUrl({
			endpoint: IMGPROXY_URL,
			url: `s3://${S3_BUCKET}/${media.thumbnail}`,
			options: { resizing_type: 'auto', height: 450, enlarge: 1 },
			salt: IMGPROXY_SALT,
			key: IMGPROXY_KEY
		});

		previewUrl = generateImageUrl({
			endpoint: IMGPROXY_URL,
			url: `s3://${S3_BUCKET}/${media.preview}`,
			options: { raw: true, resizing_type: 'auto', height: 450, enlarge: 1 },
			salt: IMGPROXY_SALT,
			key: IMGPROXY_KEY
		});

		viewUrl = generateImageUrl({
			endpoint: IMGPROXY_URL,
			url: `s3://${S3_BUCKET}/${media.view}`,
			options: { raw: true },
			salt: IMGPROXY_SALT,
			key: IMGPROXY_KEY
		});
	}
	return {
		...media,
		thumbnailUrl,
		previewUrl,
		viewUrl,
		collectionName: c.name,
		albumName: a.name
	};
}
