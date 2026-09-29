import database from '$lib/server/database';
import { toNamespacedPath } from 'node:path/posix';
import type { PageServerLoad } from './$types';
import { fail } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params }) => {
	const albums = await database.collection('collections').find({ status: 'active' }).toArray();

	return {
		albums: albums.map((album) => ({ ...album, _id: album._id.toString() }))
	};
};

export const actions = {
	// This is a named action called 'create'
	create: async ({ request }) => {
		const data = await request.formData();
		const name = data.get('name');

		if (name === 'test') {
			return fail(400, { error: 'invalid name' });
		}

		await database.collection('collections').insertOne({
			name: name,
			status: 'active'
		});

		return { success: true };
	}
};
