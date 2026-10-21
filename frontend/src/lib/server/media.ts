import { env } from '$env/dynamic/private';
import type { Media } from '$lib/server/database'
import { generateImageUrl } from '@imgproxy/imgproxy-node';

export function populateUrls(media: any): any {
    let thumbnail_url, preview_url, view_url = ""
    if (media.processStatus === 'processed') {
        thumbnail_url = generateImageUrl({
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

        preview_url = generateImageUrl({
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

        view_url = generateImageUrl({
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
        thumbnail_url, 
        preview_url, 
        view_url
    }
}