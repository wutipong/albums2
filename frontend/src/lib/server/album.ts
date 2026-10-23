import { populateUrls } from '$lib/server/media';
import { Media } from '$lib/server/database';
import { ObjectId } from 'bson';
import { PAGE_SIZE } from '$lib/const';

export async function loadAlbumMediaPage(id: string, pageno: number) {
	const page = await Media.aggregatePaginate(
		Media.aggregate([{ $match: { albumId: new ObjectId(id) } }]),
		{
			page: pageno,
			limit: PAGE_SIZE,
			lean: true
		}
	);

	for (let i = 0; i < page.docs.length; i++) {
		let media = await populateUrls(page.docs[i]);

		page.docs[i] = media;
	}

	return page;
}
