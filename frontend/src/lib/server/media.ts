import { env } from '$env/dynamic/private';
import type { Media } from '$lib/server/database'
import { generateImageUrl } from '@imgproxy/imgproxy-node';

export function populateUrls(media: any): any {
    let thumbnailUrl, previewUrl, viewUrl = ""
    if (media.processStatus === 'processed') {
        thumbnailUrl = generateImageUrl({
            endpoint: env.IMGPROXY_URL,
            url: `s3://${env.S3_BUCKET}/${media.thumbnail}`,
            options: {
                resizing_type: 'auto',
                height: 450,
                enlarge: 1
            },
            salt: env.IMGPROXY_SALT,
            key: env.IMGPROXY_KEY
        });

        previewUrl = generateImageUrl({
            endpoint: env.IMGPROXY_URL,
            url: `s3://${env.S3_BUCKET}/${media.preview}`,
            options: {
                raw: true,
                resizing_type: 'auto',
                height: 450,
                enlarge: 1
            },
            salt: env.IMGPROXY_SALT,
            key: env.IMGPROXY_KEY
        });

        viewUrl = generateImageUrl({
            endpoint: env.IMGPROXY_URL,
            url: `s3://${env.S3_BUCKET}/${media.view}`,
            options: {
                raw: true
            },
            salt: env.IMGPROXY_SALT,
            key: env.IMGPROXY_KEY
        });
    }
    return {
        ...media,
        thumbnailUrl, 
        previewUrl, 
        viewUrl
    }
}