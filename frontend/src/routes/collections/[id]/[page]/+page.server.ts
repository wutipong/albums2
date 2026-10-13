import { type PageServerLoad } from './$types';
import { Collection, Album } from '$lib/server/database';
import log from '$lib/log'

const PAGE_SIZE = 30

export const load: PageServerLoad = async({ params }) => {
	const { id, page } = params;

    const albumsResp = await Album.paginate({
        collectionId: id, 
    }, {
        page: page,
        limit: PAGE_SIZE,
        lean: true,
    })

    log.debug(albumsResp, "albums response")

    return JSON.parse(JSON.stringify(albumsResp))
};
