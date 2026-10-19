import { type PageServerLoad } from './$types'
import { encodeText } from '$lib/server/grpc/mlservice'
import { error } from '@sveltejs/kit'
import { Media } from '$lib/server/database'

export const load: PageServerLoad = async ({ params, url }) => {
    const search = url.searchParams.get('search')
    if (!search){
        error(400, "search parameter is required.")
    }
    const embedding = await encodeText(search)
    const media = await searchByVector(embedding.embedding)
    
    return {
        id: params.id,
        search: search,
        media: media,
    }
}

async function searchByVector(queryEmbedding: number[]) {
  return await Media.aggregate([
    {
      $vectorSearch: {
        index: 'media_search_index',
        path: 'embedding',
        queryVector: queryEmbedding, // Float array generated from your SigLIP python client/server
        numCandidates: 100,
        limit: 10
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
