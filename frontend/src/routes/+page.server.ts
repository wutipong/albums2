import type { PageServerLoad, PageServerLoadEvent } from './$types';
import { fail, redirect } from '@sveltejs/kit';
import { Collection } from '#lib/server/database.js';
import log from '#lib/log.js';
import { generateImageUrl } from '@imgproxy/imgproxy-node';
import { IMGPROXY_URL, S3_BUCKET, IMGPROXY_SALT, IMGPROXY_KEY } from '$app/env/private';
import { MongooseError } from 'mongoose';

export const load: PageServerLoad = async ({ locals }: PageServerLoadEvent) => {
	if (!locals.session) {
		redirect(302, '/login');
	}
	const collections = await Collection.find().sort({ name: 1 }).lean();

	log.debug(collections, 'available collections');

	for (let i = 0; i < collections.length; i++) {
		const c = collections[i];
		let cover_url = '';
		if (c.cover !== undefined && c.cover !== '') {
			cover_url = generateImageUrl({
				endpoint: IMGPROXY_URL,
				url: `s3://${S3_BUCKET}/${c.cover}`,
				options: { resizing_type: 'fill', width: 300, height: 450, enlarge: 1 },
				salt: IMGPROXY_SALT,
				key: IMGPROXY_KEY
			});
		}
		collections[i] = { ...c, cover_url };
	}
	return {
		collections
	};
};

export const actions = {
	create: async ({ request }) => {
		try {
			const data = await request.formData();
			const name = data.get('name');
			const description = data.get('description');
			const prefix = data.get('prefix');

			if (!prefix || typeof prefix !== 'string' || prefix.trim() === '') {
				return fail(400, { message: 'Collection prefix is required' });
			}

			if (!prefix.endsWith('/')) {
				return fail(400, { message: 'Collection path prefix must end with a slash (/)' });
			}
			const newCol = new Collection({ name, description, prefix });

			await newCol.save();

			return { success: true };
		} catch (e: any) {
			log.error(e, 'error creating collection');
			if (e instanceof MongooseError) {
				return fail(400, { message: e.message });
			} else {
				return fail(400, { message: 'unknown error' });
			}
		}
	}
};
