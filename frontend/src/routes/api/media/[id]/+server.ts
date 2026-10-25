import type { RequestHandler } from './$types';
import { error } from '@sveltejs/kit';
import { Media } from '#lib/server/database.js';
import * as mime from 'mime-types';
import { notifyProcessMedia } from '#lib/server/notify.js';

export const PATCH: RequestHandler = async ({ request, params }) => {
	const req = await request.json();
	const id = params.id;
	const success = req.success;

	const asset = await Media.findById(id).where({ deletedAt: null });
	if (!asset) {
		return error(500, 'media data does not existed.');
	}

	if (success) {
		asset.processStatus = 'pending';

		const mimetype = mime.lookup(asset.name ?? '');
		if (!mimetype) {
			return Response.json({ success: false, error: 'invalid content type' }, { status: 400 });
		}

		if (mimetype.startsWith('image/')) {
			asset.type = 'image';
		} else if (mimetype.startsWith('video')) {
			asset.type = 'video';
		}
	} else {
		asset.processStatus = 'failed';
		asset.original = '';
	}
	await asset.save();

	if (success) {
		try {
			await notifyProcessMedia(asset.id);
		} catch (error) {
			return Response.json(
				{ success: false, error: 'Failed to notify asset processing' },
				{ status: 500 }
			);
		}
	}

	return Response.json({ asset: asset, success: true });
};
