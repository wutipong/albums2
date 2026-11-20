import type { PageServerLoad } from './$types';
import { Settings } from '#lib/server/database.js';

export const load: PageServerLoad = async () => {
    return {
        settings: await Settings.findOneAndUpdate(
            {}, {
            $setOnInsert: { 
                album: {
                    banedNames: []
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

        await Settings.findOneAndUpdate(
            {},
            { $addToSet: { 'album.banedNames': name.trim() } },
            { upsert: true }
        );
        
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

        await Settings.findOneAndUpdate(
            {},
            { $pull: { 'album.banedNames': name.trim() } }
        );

        return { success: true };
    }
};