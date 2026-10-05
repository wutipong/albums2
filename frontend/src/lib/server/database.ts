import { env } from '$env/dynamic/private';
import { MongoClient, ObjectId } from 'mongodb';
import mongoose from 'mongoose';

export const client = new MongoClient(env.DB_CONNECTION);
export const db = client.db();

// Define a global plugin to map _id to id
mongoose.plugin((schema) => {
	schema.set('toJSON', {
		transform: function (doc, ret) {
			ret.id = ret._id;
			delete ret._id;
			Reflect.deleteProperty(ret, '__v');
			return ret;
		}
	});

	// Do the same for toObject if you use .toObject() manually
	schema.set('toObject', {
		transform: function (doc, ret) {
			ret.id = ret._id;
			delete ret._id;
			return ret;
		}
	});
});

const collectionSchema = new mongoose.Schema(
	{
		name: { type: String, index: true, unique: true },
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

const albumSchema = new mongoose.Schema(
	{
		name: String,
		collectionId: ObjectId,
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

const mediaSchema = new mongoose.Schema(
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
		deletedAt: Date
	},
	{
		timestamps: true,
		query: {
			timestamps: true
		}
	}
);

if (import.meta.hot) {
	import.meta.hot.accept((newModule) => {
		// Delete the model from mongoose so the next execution compiles it fresh
		delete mongoose.models.Collection;
		delete mongoose.models.Album;
		delete mongoose.models.Media;
	});
}

export const Collection =
	mongoose.models.Collection || mongoose.model('Collection', collectionSchema);
export const Album = mongoose.models.Album || mongoose.model('Album', albumSchema);
export const Media = mongoose.models.Media || mongoose.model('Media', mediaSchema);
