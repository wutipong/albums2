import type { PageServerLoad } from './$types';
import { fail } from '@sveltejs/kit';
import { Collection } from '$lib/server/database';
import log from '$lib/log'

export const load: PageServerLoad = async ({ params }) => {
	const rawCollections = await Collection.find().lean();
	const collections = JSON.parse(JSON.stringify(rawCollections));

	log.debug(collections, "available collections");

	return {
		collections,
	};
};

export const actions = {
	create: async ({ request }) => {
		const data = await request.formData();
		const name = data.get('name');

		if (name === 'test') {
			return fail(400, { error: 'invalid name' });
		}

        const newCol = new Collection({
            name: name,
            status: 'active',
        })
        await newCol.save()

		return { success: true };
	}
};
