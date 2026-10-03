import type { RequestHandler } from './$types';
import { error, json } from '@sveltejs/kit';
import { Album, Media } from '$lib/server/database';
import { s3 } from '$lib/server/s3';
import * as mime from 'mime-types';
import { randomUUID } from 'node:crypto';
import path from 'node:path';

export const POST: RequestHandler = async ({ request }) => {
	const req = await request.json();
	const albumId = req.album_id;
	const name = req.filename;

	const contentType = mime.contentType(path.basename(name));

	if (!contentType) {
		return error(400, { message: 'unable to recognize filetype' });
	}

	const extension = mime.extension(mime.lookup(name) || '');
	const key = `public/${randomUUID()}.${extension}`;

	const type = contentType.substring(0, contentType.indexOf('/'));

	if (type != 'image' && type != 'video') {
		return json({ success: false, error: 'Unsupported asset type.' }, { status: 400 });
	}

	const album = await Album.findById(albumId).where('deletedAt', null);

	if (!album) {
		return error(404, { message: 'album not found' });
	}

	const existing = await Media.findOne({ albumId: albumId, name: name });
	if (existing && existing.processStatus != 'uploading' && existing.processStatus != 'failed') {
		return error(409, { message: 'duplicate asset' });
	}

	let asset = null;
	if (existing) {
		asset = existing;
		asset.original = key;
		asset.type = type;
		asset.processStatus = 'uploading';

		await asset.save();
	} else {
		asset = await Media.create({
			albumId: albumId,
			name: name,
			type: type,
			processStatus: 'uploading',
			original: key
		});
	}

	if (!asset) {
		return error(500, { message: 'Failed to create asset' });
	}
	if (!asset.original) {
		return error(500, { message: 'asset.original is missing' });
	}

	const url = s3.presign(asset.original, {
		method: 'PUT',
		expiresIn: 3600,
		type: contentType
	});

	return json({ id: asset.id, url, success: true });
};
