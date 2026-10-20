import { type PageServerLoad } from './$types';
import { encodeText } from '$lib/server/grpc/mlservice';
import { error } from '@sveltejs/kit';
import { Media } from '$lib/server/database';
import { generateImageUrl } from '@imgproxy/imgproxy-node';
import { env } from '$env/dynamic/private';
import { ObjectId } from 'bson';

export const ssr = false

export const load: PageServerLoad = async ({ params, url }) => {
	const search = url.searchParams.get('search');
	if (!search) {
		error(400, 'search parameter is required.');
	}
	const embedding = await encodeText(search);
	const media = await searchByVector(params.id, embedding.embedding);

	for (let i = 0; i < media.length; i++) {
		let preview_url = ''
		let thumbnail_url = ''
		let view_url = ''

		let m = media[i]
		if (m.processStatus === 'processed') {
			thumbnail_url = generateImageUrl({
				endpoint: env.IMGPROXY_URL,
				url: `s3://${env.S3_BUCKET}/${m.thumbnail}`,
				options: {
					resizing_type: 'auto',
					height: 450,
					enlarge: 1
				},
				salt: env.IMGPROXY_SALT,
				key: env.IMGPROXY_KEY
			});

			preview_url = generateImageUrl({
				endpoint: env.IMGPROXY_URL,
				url: `s3://${env.S3_BUCKET}/${m.preview}`,
				options: {
					raw: true,
					resizing_type: 'auto',
					height: 450,
					enlarge: 1
				},
				salt: env.IMGPROXY_SALT,
				key: env.IMGPROXY_KEY
			});

			view_url = generateImageUrl({
				endpoint: env.IMGPROXY_URL,
				url: `s3://${env.S3_BUCKET}/${m.view}`,
				options: {
					raw: true
				},
				salt: env.IMGPROXY_SALT,
				key: env.IMGPROXY_KEY
			});
		}

		media[i] = {
			...m,
			preview_url,
			thumbnail_url,
			view_url,
		}
	}


	return {
		collectionId: params.id,
		search: search,
		media: media
	};
};

async function searchByVector(collectionId: string, queryEmbedding: number[]) {
	return await Media.aggregate([
		{
			$vectorSearch: {
				index: 'media_search_index',
				path: 'embedding',
				queryVector: queryEmbedding, 
				numCandidates: 100,
				limit: 30,
				filter: { "collectionId": new ObjectId(collectionId) }
			}
		},
		{
			$project: {
				imageEmbedding: 0,
				score: { $meta: 'vectorSearchScore' }
			}
		}
	]);
}
