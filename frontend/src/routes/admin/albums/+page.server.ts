import type { PageServerLoad } from './$types';
import { Settings } from '#lib/server/database.js';
import { fail } from '@sveltejs/kit';

export const load: PageServerLoad = async () => {
	return {
		settings: await Settings.findOneAndUpdate(
			{},
			{
				$setOnInsert: {
					album: {
						excludedCoverNames: []
					}
				}
			},
			{
				upsert: true,
				returnDocument: 'after'
			}
		).lean()
	};
};

export const actions = {
	add: async ({ request }) => {
		const formData = await request.formData();
		const name = formData.get('name');

		if (typeof name !== 'string' || !name.trim()) {
			return fail(400, { message: 'Media name is required.' });
		}

		const r = await Settings.findOneAndUpdate(
			{},
			{ $addToSet: { 'album.excludedCoverNames': name.trim() } },
			{ upsert: true }
		);

		if (r?.errors) {
			return fail(500, { message: 'Failed to add media name.' });
		}

		return { success: true };
	},
	delete: async ({ request }) => {
		const formData = await request.formData();
		const name = formData.get('name');

		if (typeof name !== 'string' || !name.trim()) {
			return fail(400, { message: 'Media name is required.' });
		}

		const r = await Settings.findOneAndUpdate(
			{},
			{ $pull: { 'album.excludedCoverNames': name.trim() } }
		);

		if (r?.errors) {
			return fail(500, { message: 'Failed to remove media name.' });
		}

		return { success: true };
	}
};
