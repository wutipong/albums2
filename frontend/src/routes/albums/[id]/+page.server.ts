import type { PageServerLoad } from './$types';
import { Album, Collection } from '#lib/server/database.js';
import { loadAlbumMediaPage } from '#lib/server/album.js';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params, locals }) => {
	const { id } = params;
	const album = await Album.findById(id).lean();
	const media = await loadAlbumMediaPage(id, 1);
	const collection = await Collection.findById(album.collectionId).lean();
	return {
		collection: collection,
		album: album,
		media: media,
		user: locals.user
	};
};

export const ssr = false;
