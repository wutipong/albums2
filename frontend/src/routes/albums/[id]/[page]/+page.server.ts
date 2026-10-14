import { type PageServerLoad } from './$types';
import { Media } from '$lib/server/database';
import log from '$lib/log'
import {env} from '$env/dynamic/private'
import { generateImageUrl } from '@imgproxy/imgproxy-node';

const PAGE_SIZE = 30

export const load: PageServerLoad = async ({ params }) => {
    const { id, page } = params;

    const resp = await Media.paginate({
        albumId: id,
    }, {
        page: page,
        limit: PAGE_SIZE,
        lean: true,
    })

    log.debug(resp, "albums response")

    for (let i = 0; i < resp.docs; i++) {
        let preview_url = ''
        let thumbnail_url = ''

        let media = resp.docs[i]
        if (media.processStatus == 'processed') {
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

    return JSON.parse(JSON.stringify(resp))
};
