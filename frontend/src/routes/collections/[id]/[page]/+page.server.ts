import { type PageServerLoad } from './$types';
import { Collection, Album } from '$lib/server/database';
import { env } from '$env/dynamic/private';
import log from '$lib/log';
import { generateImageUrl } from '@imgproxy/imgproxy-node';
import { ObjectId } from 'bson';

const PAGE_SIZE = 30;

export const load: PageServerLoad = async ({ params }) => {
	const { id, page } = params;
	const collection = await Collection.findById(id).lean();

	const albums = await Album.aggregatePaginate(
		Album.aggregate([{ $match: { collectionId: new ObjectId(id) } }]),
		{
			page: page,
			limit: PAGE_SIZE,
			lean: true
		}
	);

	log.debug(albums, 'albums response');

	for (let i = 0; i < albums.docs.length; i++) {
		const album = albums.docs[i];
		let cover_url = '';
		if (album.cover !== undefined && album.cover !== '') {
			cover_url = generateImageUrl({
				endpoint: env.IMGPROXY_URL,
				url: `s3://${env.S3_BUCKET}/${album.cover}`,
				options: {
					resizing_type: 'fill',
					width: 400,
					height: 600,
					enlarge: 1
				},
				salt: env.IMGPROXY_SALT,
				key: env.IMGPROXY_KEY
			});
		}
		albums.docs[i] = { ...album, cover_url };
	}

	return {
		collection,
		albums
	};
};
