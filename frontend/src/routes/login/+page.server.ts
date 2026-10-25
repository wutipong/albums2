import { OIDC_PROVIDER_ID } from '$app/env/private';
import { type PageServerLoad } from './$types';

export const load: PageServerLoad = ({ params }) => {
	return { oidcProviderId: OIDC_PROVIDER_ID };
};
