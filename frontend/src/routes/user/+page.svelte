<script lang="ts">
	import NavBar from '#lib/components/NavBar.svelte';
	import { mdiLogout, mdiLibraryShelves } from '@mdi/js';
	import Icon from '#lib/components/Icon.svelte';
	import { createHash } from '@better-auth/utils/hash';

	let { data } = $props();

	const hashVal = $derived.by(() => createHash('SHA-256', 'hex').digest(data.user?.email ?? ''));
	const avatarSrc = $derived(`https://gravatar.com/avatar/${hashVal}`);
</script>

<svelte:head>
	<title>Albums: User</title>
</svelte:head>

<div class="flex h-screen w-screen flex-col">
	<NavBar
		user={data.user}
		up={{ description: 'Collections', icon: mdiLibraryShelves, url: '/collections' }}
	/>

	<div class="overflow-auto p-4 pt-8">
		<article class="mx-auto prose h-full w-full md:w-200">
			<div class="flex flex-row gap-8">
				<div class="avatar">
					<div class="h-24 w-24 rounded-full">
						<img src={avatarSrc} alt="avatar" class="my-0!" />
					</div>
				</div>
				<div>
					<h2 class="mt-0">{data.user?.name}</h2>
					{#if data.user?.role && data.user?.role == 'admin'}
						<h3>Administrator</h3>
					{/if}
					<p><a href={`mailto:${data.user?.email}`}>{data.user?.email}</a></p>
					<a class="btn btn-soft" href="/logout" data-sveltekit-preload-data="false">
						<Icon path={mdiLogout} />Logout
					</a>
				</div>
			</div>
		</article>
	</div>
</div>
