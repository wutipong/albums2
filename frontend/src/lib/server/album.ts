import { populateUrls } from '#lib/server/media.js';
import { Media } from '#lib/server/database.js';
import { ObjectId } from 'bson';
import { PAGE_SIZE } from '#lib/const.js';
import log from '#lib/log.ts';

export async function loadAlbumMediaPage(id: string, pageno: number) {
	const page = await Media.aggregatePaginate(
		Media.aggregate([{ $match: { albumId: new ObjectId(id) } }]),
		{
			page: pageno,
			limit: PAGE_SIZE
		}
	);

	log.debug(page, 'loaded media page for album', id, 'page', pageno);

	for (let i = 0; i < page.docs.length; i++) {
		const media = await populateUrls(page.docs[i]);

		page.docs[i] = media;
	}

	return page;
}
