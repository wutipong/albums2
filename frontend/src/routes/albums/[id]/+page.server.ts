import type { PageServerLoad } from './$types';
import { Album, Collection } from '#lib/server/database.js';
import { loadAlbumMediaPage } from '#lib/server/album.js';

export const load: PageServerLoad = async ({ params, locals, url }) => {
	const { id } = params;
	const album = await Album.findById(id).lean();
	const media = await loadAlbumMediaPage(id, 1);
	const collection = await Collection.findById(album.collectionId).lean();

	const fromPage = url.searchParams.get('fromPage');

	return {
		collection: collection,
		album: album,
		media: media,
		fromPage: fromPage == null ? 1 : Number(fromPage),
		user: locals.user
	};
};

export const ssr = false;
