import type { RequestHandler } from './$types';
import { loadAlbumMediaPage } from '#lib/server/album.js';
import { json } from '@sveltejs/kit';

export const GET: RequestHandler = async ({ params }) => {
	const { albumId, page } = params;
	const resp = await loadAlbumMediaPage(albumId, page);

	return json(resp);
};
