import type { PageServerLoad } from './$types';
import { Settings } from '#lib/server/database.js';

export const load: PageServerLoad = async () => {
    return {
        settings: await Settings.findOneAndUpdate(
            {}, {
            $setOnInsert: { 
                album: {
                    excludedCoverNames: []
                }
            }
        },
            { 
                upsert: true, 
                returnDocument: 'after' 
            },
        ).lean()
    };
};

export const actions = {
    add: async ({ request }) => {
        const formData = await request.formData();
        const name = formData.get('name');

        if (typeof name !== 'string' || !name.trim()) {
            return Response.json(
                { error: 'Media name is required.' }, 
                { status: 400 });
        }

        const r = await Settings.findOneAndUpdate(
            {},
            { $addToSet: { 'album.excludedCoverNames': name.trim() } },
            { upsert: true }
        );

        if (r?.errors) {
            return Response.json(
                { error: 'Failed to add media name.' },
                { status: 500 }
            );
        }
        
        return { success: true };
    },
    delete: async ({ request }) => {
        const formData = await request.formData();
        const name = formData.get('name');

        if (typeof name !== 'string' || !name.trim()) {
            return Response.json(
                { error: 'Media name is required.' },
                { status: 400 });
        }

        const r = await Settings.findOneAndUpdate(
            {},
            { $pull: { 'album.excludedCoverNames': name.trim() } }
        );

        if (r?.errors) {
            return Response.json(
                { error: 'Failed to remove media name.' },
                { status: 500 }
            );
        }

        return { success: true };
    }
};