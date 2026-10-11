import { z } from 'zod';
import mongoose from 'mongoose';
import { toMongooseSchema, genTimestampsSchema, withMongoose } from '@nullix/zod-mongoose'; // 👈 Fully typed natively

const CollectionZodSchema = withMongoose(
	z.object({
		name: z.string(),
		...genTimestampsSchema()
	}),
	{ timestamps: true }
);

export const CollectionSchema = toMongooseSchema(CollectionZodSchema);
CollectionSchema.index({ name: 1 }, { unique: true });

const AlbumZodSchema = withMongoose(
	z.object({
		name: z.string(),
		cover: z.string(),
		collectionID: withMongoose(z.unknown(), {
			type: mongoose.Schema.Types.ObjectId,
			ref: 'Collection',
			required: true
		}),
		...genTimestampsSchema()
	}),
	{ timestamps: true }
);
export const AlbumSchema = toMongooseSchema(AlbumZodSchema);
AlbumSchema.index({ name: 1 }, { unique: true });

const MediaZodSchema = withMongoose(
	z.object({
		name: z.string(),
        albumID: withMongoose(z.unknown(), {
			type: mongoose.Schema.Types.ObjectId,
			ref: 'Album',
			required: true
		}),
		original: z.string(),
		preview: z.string(),
		thumbnail: z.string(),
		view: z.string(),
		type: z.string(),
		viewWidth: z.int32(),
		viewHeight: z.int32(),
		duration: z.int32(),
		thumbnailWidth: z.int32(),
		thumbnailHeight: z.int32(),
		processStatus: z.string(),
        ...genTimestampsSchema()
	}),
	{ timestamps: true }
);

export const MediaSchema = toMongooseSchema(MediaZodSchema);
