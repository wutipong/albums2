import { env } from '$env/dynamic/private';
import { type PageServerLoad } from './$types';

export const load: PageServerLoad = ({ params }) => {
	return {
		oidcProviderId: env.OIDC_PROVIDER_ID
	};
};
