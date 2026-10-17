import { env } from '$env/dynamic/private'
import { generateImageUrl } from '@imgproxy/imgproxy-node';
import { Media } from "./database";

const PAGE_SIZE = 30

export async function loadAlbumMediaPage(id: string, pageno: number) {
    const page = await Media.paginate({
        albumId: id,
    }, {
        page: pageno,
        limit: PAGE_SIZE,
        lean: true,
    })

    for (let i = 0; i < page.docs.length; i++) {
        let preview_url = ''
        let thumbnail_url = ''
        let view_url = ''

        let media = page.docs[i]
        if (media.processStatus === 'processed') {
            thumbnail_url = generateImageUrl({
                endpoint: env.IMGPROXY_URL,
                url: `s3://${env.S3_BUCKET}/${media.thumbnail}`,
                options: {
                    resizing_type: 'auto',
                    height: 300,
                    enlarge: 1
                },
                salt: env.IMGPROXY_SALT,
                key: env.IMGPROXY_KEY
            });

            preview_url = generateImageUrl({
                endpoint: env.IMGPROXY_URL,
                url: `s3://${env.S3_BUCKET}/${media.preview}`,
                options: {
                    raw: media.type !== 'image',
                    resizing_type: 'auto',
                    height: 300,
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

        page.docs[i] = {
            ...media,
            preview_url,
            thumbnail_url,
            view_url,
        }
    }

    return page
}