<script lang="ts">
	import { authClient } from '$lib/auth';

	const { data } = $props();
</script>

<button
	class="btn btn-primary"
	onclick={() =>
		authClient.signIn.social({
			provider: data.oidcProviderId,
			callbackURL: '',
			errorCallbackURL: ''
		})}
>
	Sign-In
</button>
