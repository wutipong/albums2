import { type PageServerLoad } from './$types'
import { encodeText } from '$lib/server/grpc/mlservice'
import { error } from '@sveltejs/kit'

export const load: PageServerLoad = async ({ params, url }) => {
    const search = url.searchParams.get('search')
    if (!search){
        error(400, "search parameter is required.")
    }
    const embedding = await encodeText(search)
    
    return {
        id: params.id,
        search: search,
        embedding: embedding
    }
}