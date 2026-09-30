import { env } from '$env/dynamic/private';
import { MongoClient, ObjectId } from 'mongodb';
import mongoose from 'mongoose';

export const client = new MongoClient(env.DB_CONNECTION);
export const db = client.db();

const collectionSchema = new mongoose.Schema({
	name: { type: [String], index: true, unique: true },
	status: String
});

const albumSchema = new mongoose.Schema({
	name: String,
	collectionSchema: ObjectId
});

const mediaSchema = new mongoose.Schema({
	name: String,
	albumId: ObjectId
});

export const Collection = mongoose.model('Collection', collectionSchema);
export const Album = mongoose.model('Album', albumSchema);
export const Media = mongoose.model('Media', mediaSchema);
