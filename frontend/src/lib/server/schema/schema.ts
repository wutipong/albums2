import mongoose from 'mongoose';
import { Int32, ObjectId } from 'mongodb';

export interface Collection {
	name?: string;
	description?: string;
	cover?: string;
	s3Prefix: string;
	createdAt?: Date;
	modifiedAt?: Date;
	updatedAt?: Date;
	deletedAt?: Date;
}

export interface Album {
	name?: string;
	collectionId?: ObjectId;
	cover?: string;
	createdAt?: Date;
	modifiedAt?: Date;
	updatedAt?: Date;
	deletedAt?: Date;
}

export interface Media {
	name?: string;
	collectionId?: ObjectId;
	albumId?: ObjectId;
	original?: string;
	type?: 'image' | 'video';
	processStatus?: 'uploading' | 'failed' | 'pending';
	view?: string;
	viewWidth?: number;
	viewHeight?: number;
	videoDuration?: number;
	thumbnail?: string;
	thumbnailWidth?: number;
	thumbnailHeight?: number;
	preview?: string;
	imageFrames?: number;
	embedding?: number[];
	createdAt?: Date;
	modifiedAt?: Date;
	updatedAt?: Date;
	deletedAt?: Date;
}

export const CollectionSchema = new mongoose.Schema<Collection>(
	{
		name: { type: String, index: true, unique: true },
		description: String,
		cover: String,
		s3Prefix: { type: String, required: true },
		createdAt: { type: Date, default: Date.now },
		modifiedAt: { type: Date, default: Date.now },
		deletedAt: Date
	},
	{
		timestamps: true
	}
);

export const AlbumSchema = new mongoose.Schema<Album>(
	{
		name: String,
		collectionId: ObjectId,
		cover: String,
		createdAt: { type: Date, default: Date.now },
		modifiedAt: { type: Date, default: Date.now },
		deletedAt: Date
	},
	{
		timestamps: true
	}
);

export const MediaSchema = new mongoose.Schema<Media>(
	{
		name: String,
		collectionId: ObjectId,
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
		timestamps: true
	}
);

MediaSchema.searchIndex({
	name: 'media_search_index',
	type: 'vectorSearch',
	definition: {
		fields: [
			{
				type: 'vector',
				path: 'embedding',
				numDimensions: 768,
				similarity: 'cosine'
			},
			{
				type: 'filter',
				path: 'collectionId'
			}
		]
	}
});
