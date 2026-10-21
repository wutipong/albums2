import { env } from '$env/dynamic/private';
import { Album, Collection } from '$lib/server/database'
import { generateImageUrl } from '@imgproxy/imgproxy-node';

export async function populateUrls(media: any): Promise<any> {
    const c = await Collection.findById(media.collectionId)
    const a = await Album.findById(media.albumId)

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
        viewUrl,
        collectionName: c.name,
        albumName: a.name
    }
}