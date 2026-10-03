import type { RequestHandler } from '@sveltejs/kit';
import { error, json } from '@sveltejs/kit';
import { Album, Media } from '$lib/server/database';
import log from '$lib/log';
// import { deleteAlbum } from '$lib/server/grpc/worker';

export const GET: RequestHandler = async ({ params }) => {
	log.debug(params, 'GET album params');
	const { album_id } = params;

	if (!album_id) {
		return error(400, { message: 'album id is required.' });
	}

	const album = await Album.findById(album_id).where('deletedAt', null).lean();
	if (!album) {
		return error(404, { message: 'album not found' });
	}
	return json({ album });
};

export const DELETE: RequestHandler = async ({ params }) => {
	const { album_id } = params;

	if (!album_id) {
		return error(400, { message: 'album id is required.' });
	}

	// await deleteAlbum(album_id);

	return json({ success: true });
};
