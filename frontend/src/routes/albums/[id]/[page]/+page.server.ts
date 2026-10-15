import { type PageServerLoad } from './$types';
import { Album, Media } from '$lib/server/database';
import type { ObjectId } from 'mongoose';
import log from '$lib/log'
import { env } from '$env/dynamic/private'
import { generateImageUrl } from '@imgproxy/imgproxy-node';

const PAGE_SIZE = 30

export const load: PageServerLoad = async ({ params }) => {
    const { id, page } = params;

    const album = await Album.findById(id).lean()

    const resp = await Media.paginate({
        albumId: id,
    }, {
        page: page,
        limit: PAGE_SIZE,
        lean: true,
    })

    for (let i = 0; i < resp.docs.length; i++) {
        let preview_url = ''
        let thumbnail_url = ''

        let media = resp.docs[i]
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
        }

        resp.docs[i] = {
            ...media,
            preview_url,
            thumbnail_url,
        }
    }

    return {
        album: album,
        media: resp
    }
};
