import { type PageServerLoad } from './$types';
import { Album } from '$lib/server/database';
import { loadAlbumMediaPage } from '$lib/server/album';
import { redirect } from '@sveltejs/kit';

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


export const actions = {
    search: async ({ request }) => {

        const data = await request.formData();
        const text = data.get('text');
        const collectionId = data.get('collectionId')

        const query = new URLSearchParams({ search: text });

        redirect(303, `/collections/${collectionId}/search?${query.toString()}`)
    }
};
