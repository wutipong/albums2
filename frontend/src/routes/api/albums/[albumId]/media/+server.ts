import type { RequestHandler } from './$types';
import { error } from '@sveltejs/kit';
import { Album, Media } from '#lib/server/database.js';
import { s3 } from '#lib/server/s3.js';
import * as mime from 'mime-types';
import { randomUUIDv7 } from 'node:crypto';
import path from 'node:path';
import log from '#lib/log.js';

export const POST: RequestHandler = async ({ request, params }) => {
	const req = await request.json();
	const name = req.filename;
	const contentType = mime.contentType(path.basename(name));

	const { albumId } = params;

	log.info({ albumId, name, contentType }, `POST ${request.url}`);
	if (!contentType) {
		return error(400, 'unable to recognize filetype');
	}

	const extension = mime.extension(mime.lookup(name) || '');
	const key = `public/${randomUUIDv7()}.${extension}`;

	const type = contentType.substring(0, contentType.indexOf('/'));

	if (type != 'image' && type != 'video') {
		return Response.json({ success: false, error: 'Unsupported asset type.' }, { status: 400 });
	}

	const album = await Album.findById(albumId).where('deletedAt', null);

	if (!album) {
		return error(404, 'album not found');
	}

	const existing = await Media.findOne({ albumId: albumId, name: name });
	if (existing && existing.processStatus != 'uploading' && existing.processStatus != 'failed') {
		return error(409, 'duplicate asset');
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
			collectionId: album.collectionId,
			albumId: albumId,
			name: name,
			type: type,
			processStatus: 'uploading',
			original: key
		});
	}

	if (!asset) {
		return error(500, 'Failed to create asset');
	}
	if (!asset.original) {
		return error(500, 'asset.original is missing');
	}

	const url = s3.presign(asset.original, {
		method: 'PUT',
		expiresIn: 3600,
		type: contentType
	});

	return Response.json({ id: asset.id, url, success: true });
};

export const GET: RequestHandler = async ({ request, params }) => {
	log.debug(params, `GET ${request.url}`);
	const { albumId } = params;
	if (await Album.exists({ albumId: albumId, deletedAt: null })) {
		return error(404, 'album not found');
	}

	const media = await Media.find({ albumId: albumId, deletedAt: null }).lean();

	return Response.json({ media });
};
