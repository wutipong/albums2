import { type PageServerLoad } from './$types'

export const load: PageServerLoad = async ({ params, url }) => {
    const search = url.searchParams.get('search')

    return {
        id: params.id,
        search: search
    }
}