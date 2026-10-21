import { populateUrls } from '$lib/server/media'
import { Media } from '$lib/server/database';
import { ObjectId } from 'bson';

const PAGE_SIZE = 30

export async function loadAlbumMediaPage(id: string, pageno: number) {
    const page = await Media.aggregatePaginate(
        Media.aggregate([
            { $match: { albumId: new ObjectId(id) } }
        ]), {
        page: pageno, limit: 10, lean: true,
    })

    for (let i = 0; i < page.docs.length; i++) {
        let media = populateUrls(page.docs[i])

        page.docs[i] = media
    }

    return page
}