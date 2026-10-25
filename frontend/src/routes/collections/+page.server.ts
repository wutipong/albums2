import type { PageServerLoad } from './$types';
import { error, fail } from '@sveltejs/kit';
import { Collection } from '#lib/server/database.js';
import log from '#lib/log.js';
import { generateImageUrl } from '@imgproxy/imgproxy-node';
import { IMGPROXY_URL, S3_BUCKET, IMGPROXY_SALT, IMGPROXY_KEY } from '$app/env/private';

export const load: PageServerLoad = async ({ params }) => {
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
			const newCol = new Collection({ name, description });

			await newCol.save();

			return { success: true };
		} catch (e) {
			console.log(e);
			return fail(400, { code: e.code, message: e.errmsg });
		}
	}
};
