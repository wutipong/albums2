import { PAGE_SIZE, SEARCH_LIMIT } from "$lib/const";
import { ObjectId } from "bson";
import { Media } from "./database";
import { populateUrls } from "./media";

export async function searchMedia(collectionId: string, queryEmbedding: number[], page: number): Promise<any> {
    const aggregate = Media.aggregate([
        {
            $vectorSearch: {
                index: 'media_search_index',
                path: 'embedding',
                queryVector: queryEmbedding,
                numCandidates: SEARCH_LIMIT,
                limit: SEARCH_LIMIT,
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

    const results = await Media.aggregatePaginate(aggregate, {
        page: page,
        limit: PAGE_SIZE,
        useFacet: false,
    });

    for (let i = 0; i < results.docs.length; i++) {
        let media = await populateUrls(results.docs[i])

        results.docs[i] = media
    }

    return results
}