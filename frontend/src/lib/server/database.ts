import { DB_CONNECTION } from '$app/env/private';
import { MongoClient } from 'mongodb';
import mongoose from 'mongoose';
import { Int32, ObjectId } from 'mongodb';
import aggregatePaginate from 'mongoose-aggregate-paginate-v2';

export const client = new MongoClient(DB_CONNECTION);
export const db = client.db();

export interface Collection {
	name?: string;
	description?: string;
	cover?: string;
	prefix: string;
	createdAt?: Date;
	updatedAt?: Date;
	deletedAt?: Date;
}

export interface Album {
	name?: string;
	collectionId?: ObjectId;
	objectKey: string;
	cover?: string;
	createdAt?: Date;
	updatedAt?: Date;
	deletedAt?: Date;
}

export interface Media {
	name?: string;
	collectionId?: ObjectId;
	albumId?: ObjectId;
	original?: string;
	type?: 'image' | 'video';
	processStatus?: 'processed' | 'failed' | 'pending';
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
	updatedAt?: Date;
	deletedAt?: Date;
}

export const CollectionSchema = new mongoose.Schema<Collection>(
	{
		name: { type: String, index: true, unique: true },
		description: String,
		cover: String,
		prefix: { type: String, required: true },
		createdAt: { type: Date, default: Date.now },
		updatedAt: { type: Date, default: Date.now },
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
		updatedAt: { type: Date, default: Date.now },
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
			enum: ['processed', 'failed', 'pending']
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
		updatedAt: { type: Date, default: Date.now },
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


if (import.meta.hot) {
	import.meta.hot.accept((newModule) => {
		// Delete the model from mongoose so the next execution compiles it fresh
		delete mongoose.models.Collection;
		delete mongoose.models.Album;
		delete mongoose.models.Media;
	});
}

CollectionSchema.plugin(aggregatePaginate);
AlbumSchema.plugin(aggregatePaginate);
MediaSchema.plugin(aggregatePaginate);

export const Collection =
	(mongoose.models.Collection as mongoose.AggregatePaginateModel<any> | undefined) ??
	mongoose.model<any, mongoose.AggregatePaginateModel<any>>('Collection', CollectionSchema);
export const Album =
	(mongoose.models.Album as mongoose.AggregatePaginateModel<any> | undefined) ??
	mongoose.model<any, mongoose.AggregatePaginateModel<any>>('Album', AlbumSchema);
export const Media =
	(mongoose.models.Media as mongoose.AggregatePaginateModel<any> | undefined) ??
	mongoose.model<any, mongoose.AggregatePaginateModel<any>>('Media', MediaSchema);
