import mongoose from 'mongoose';
import { Int32, ObjectId } from 'mongodb';

export const CollectionSchema = new mongoose.Schema(
	{
		name: { type: String, index: true, unique: true },
		description: String,
		cover: String,
		createdAt: { type: Date, default: Date.now },
		modifiedAt: { type: Date, default: Date.now },
		deletedAt: Date
	},
	{
		timestamps: true,
		query: {
			timestamps: true
		}
	}
);

export const AlbumSchema = new mongoose.Schema(
	{
		name: String,
		collectionId: ObjectId,
		cover: String,
		createdAt: { type: Date, default: Date.now },
		modifiedAt: { type: Date, default: Date.now },
		deletedAt: Date
	},
	{
		timestamps: true,
		query: {
			timestamps: true
		}
	}
);

export const MediaSchema = new mongoose.Schema(
	{
		name: String,
		albumId: ObjectId,
		original: String,
		type: {
			type: String,
			enum: ['image', 'video']
		},
		processStatus: {
			type: String,
			enum: ['uploading', 'failed', 'pending']
		},
		view: String,
		viewWidth: Int32,
		viewHeight: Int32,
		videoDuration: Int32,
		thumbnail: String,
		thumbnailWidth: Int32,
		thumbnailHeight: Int32,
		preview: String,
		imageFrames: Int32,
		embedding: {
			type: [Number]
		},
		createdAt: { type: Date, default: Date.now },
		modifiedAt: { type: Date, default: Date.now },
		deletedAt: Date
	},
	{
		timestamps: true,
		query: {
			timestamps: true
		}
	}
);

MediaSchema.searchIndex({
	name: 'media_search_index',
	type: 'vectorSearch',
	definition: {
		fields: [
			{
				type: 'knnVector',
				path: 'embedding',
				dimensions: 768,       // ⚠️ Change to 1152 if using SO400M, or 1024/1536 depending on your model!
				similarity: 'dotProduct' // SigLIP outputs normalized vectors; dotProduct is highly performant
			}
		]
	}
})