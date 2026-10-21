import { env } from '$env/dynamic/private';
import { MongoClient } from 'mongodb';
import mongoose from 'mongoose';
import { CollectionSchema, AlbumSchema, MediaSchema } from '$lib/server/schema/schema';
import aggregatePaginate from 'mongoose-aggregate-paginate-v2'

export const client = new MongoClient(env.DB_CONNECTION);
export const db = client.db();

if (import.meta.hot) {
	import.meta.hot.accept((newModule) => {
		// Delete the model from mongoose so the next execution compiles it fresh
		delete mongoose.models.Collection;
		delete mongoose.models.Album;
		delete mongoose.models.Media;
	});
}

CollectionSchema.plugin(aggregatePaginate)
AlbumSchema.plugin(aggregatePaginate)
MediaSchema.plugin(aggregatePaginate)

export const Collection =
	mongoose.models.Collection || mongoose.model('Collection', CollectionSchema);
export const Album = mongoose.models.Album || mongoose.model('Album', AlbumSchema);
export const Media = mongoose.models.Media || mongoose.model('Media', MediaSchema);
