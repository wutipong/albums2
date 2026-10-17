import { type PageServerLoad } from './$types';
import { Album } from '$lib/server/database';
import { loadAlbumMediaPage } from '$lib/server/album';

export const load: PageServerLoad = async ({ params }) => {
    const { id } = params;
    const album = await Album.findById(id).lean()
    const media = await loadAlbumMediaPage(id, 1)

    return {
        album: album,
        media: media
    }
};

export const ssr = false;