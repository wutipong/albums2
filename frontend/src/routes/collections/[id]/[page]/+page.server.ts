import { type PageServerLoad } from './$types';
import { Collection, Album } from '$lib/server/database';
import { env } from '$env/dynamic/private'
import log from '$lib/log'
import { generateImageUrl } from '@imgproxy/imgproxy-node';

const PAGE_SIZE = 30

export const load: PageServerLoad = async ({ params }) => {
    const { id, page } = params;

    const resp = await Album.paginate({
        collectionId: id,
    }, {
        page: page,
        limit: PAGE_SIZE,
        lean: true,
    })

    log.debug(resp, "albums response")

    for (let i = 0; i < resp.docs.length; i++) {
        const album = resp.docs[i];
        let cover_url = ''
        if (album.cover !== undefined && album.cover !== '') {
            cover_url = generateImageUrl({
                endpoint: env.IMGPROXY_URL,
                url: `s3://${env.S3_BUCKET}/${album.cover}`,
                options: {
                    resizing_type: 'fill',
                    height: 400,
                    width: 600,
                    enlarge: 1
                },
                salt: env.IMGPROXY_SALT,
                key: env.IMGPROXY_KEY
            });
        }
        resp.docs[i] = { ...album, cover_url }
    }

    return JSON.parse(JSON.stringify(resp))
};
