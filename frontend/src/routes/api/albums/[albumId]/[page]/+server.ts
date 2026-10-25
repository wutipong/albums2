import type { RequestHandler } from './$types';
import { loadAlbumMediaPage } from '#lib/server/album.js';

export const GET: RequestHandler = async ({ params }) => {
	const { albumId, page } = params;
	const resp = await loadAlbumMediaPage(albumId, page);

	return Response.json(resp);
};
