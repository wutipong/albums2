import { env } from '$env/dynamic/private';
import { MongoClient } from 'mongodb';

export const client = new MongoClient(env.DB_CONNECTION);

export default client.db()