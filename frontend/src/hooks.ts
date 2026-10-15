import { Types } from 'mongoose';
import type { Transport } from '@sveltejs/kit';

export const transport: Transport = {
    ObjectId: {
        encode: (value) => value instanceof Types.ObjectId && [value.toString()],
        decode: ([hexString]) => new Types.ObjectId(hexString)
    }
};
