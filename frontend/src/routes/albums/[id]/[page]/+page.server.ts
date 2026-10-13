import { type PageServerLoad } from './$types';
import { Media } from '$lib/server/database';
import log from '$lib/log'

const PAGE_SIZE = 30

export const load: PageServerLoad = async({ params }) => {
	const { id, page } = params;

    const resp = await Media.paginate({
        albumId: id, 
    }, {
        page: page,
        limit: PAGE_SIZE,
        lean: true,
    })

    log.debug(resp, "albums response")

    return JSON.parse(JSON.stringify(resp))
};
