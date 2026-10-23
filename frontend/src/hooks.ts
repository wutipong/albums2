// src/hooks.ts
import { ObjectId } from 'bson';
import type { Transport } from '@sveltejs/kit';

export const transport: Transport = {
	MongooseObjectId: {
		// Server-side check
		encode: (value) => {
			if (
				value &&
				typeof value === 'object' &&
				'_bsontype' in value &&
				value._bsontype === 'ObjectId'
			) {
				// SvelteKit transport requirements: Must wrap inside an array or plain object
				return value.toString();
			}
			return false;
		},
		// Client-side execution
		// Destructure the string back out of the array to instantiate it safely
		decode: (str) => new ObjectId(str)
	}
};
