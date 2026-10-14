import type { PageServerLoad } from './$types';
import { error, fail } from '@sveltejs/kit';
import { Collection } from '$lib/server/database';
import log from '$lib/log'

export const load: PageServerLoad = async ({ params }) => {
	const rawCollections = await Collection.find().sort({name: 1}).lean();
	const collections = JSON.parse(JSON.stringify(rawCollections));

	log.debug(collections, "available collections");

	return {
		collections,
	};
};

export const actions = {
	create: async ({ request }) => {
		try {
			const data = await request.formData();
			const name = data.get('name');
			const description = data.get('description')

			const newCol = new Collection({
				name: name,
				description: description,
			})
			await newCol.save()

			return { success: true };
		} catch (e) {
			console.log(e)
			return fail(400, {code: e.code, message: e.errmsg})
		}
	}
};
