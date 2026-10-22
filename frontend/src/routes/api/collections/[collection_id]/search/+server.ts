import type { RequestHandler } from './$types';
import { searchMedia } from '$lib/server/collection';
import { json } from '@sveltejs/kit';
import { encodeText } from '$lib/server/grpc/mlservice';

export const GET: RequestHandler = async ({ params, url }) => {
    const search = url.searchParams.get("search")
    const page = parseInt(url.searchParams.get("page")?? "0")

    const embedding = await encodeText(search?? "");
    const resp = await searchMedia(params.collection_id, embedding.embedding ,page)

    return json(resp)
}