<script lang="ts">
	import { authClient } from '$lib/auth';
	import Icon from '$lib/components/Icon.svelte';
	import { mdiLogin } from '@mdi/js';

	const { data } = $props();
</script>

<svelte:head>
	<title>Albums2: Login</title>
</svelte:head>

<div class="relative h-screen w-screen items-center justify-center">
	<button
		class="btn absolute top-1/2 left-1/2 -translate-1/2 gap-4 btn-outline btn-primary btn-xl"
		onclick={() => authClient.signIn.social({ provider: data.oidcProviderId })}
	>
		<Icon path={mdiLogin} />
		Login
	</button>

	<div class="absolute right-0 bottom-0 left-0 m-4">albums2 &COPY; 2026 Wutipong Wongsakuldej</div>
</div>
