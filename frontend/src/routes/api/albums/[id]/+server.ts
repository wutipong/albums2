import type { RequestHandler } from '@sveltejs/kit';
import { error, json } from '@sveltejs/kit';
import { Album, Media } from '$lib/server/database';
// import { deleteAlbum } from '$lib/server/grpc/worker';

export const GET: RequestHandler = async ({ params }) => {
	const { id } = params;

	if (!id) {
		return error(400, { message: 'album id is required.' });
	}

	const album = await Album.findById(id).where('deletedAt', null).lean();
	if (!album) {
		return error(404, { message: 'album not found' });
	}

	const media = await Media.find({ albumId: id, deletedAt: null }).lean();
	return json({ ...album, media });
};

export const DELETE: RequestHandler = async ({ params }) => {
	const { id } = params;

	if (!id) {
		return error( 400, { message: 'album id is required.' });
	}

	// await deleteAlbum(id);

	return json({ success: true });
};