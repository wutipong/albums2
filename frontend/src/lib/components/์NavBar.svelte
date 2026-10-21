<script lang="ts">
	import { mdiLogout, mdiMagnify, mdiMenu } from '@mdi/js';
	import Icon from './Icon.svelte';
	import { enhance } from '$app/forms';
	import { onMount } from 'svelte';
	import { createHash } from '@better-auth/utils/hash';
	import { authClient } from '$lib/auth';

	interface UpDetails {
		description: string;
		icon: string;
		url: string;
	}

	interface Props {
		up?: UpDetails;
		collectionId?: string;
		menu?: any;
	}

	const { up, collectionId, menu }: Props = $props();
	let avatarSrc = $state('');
    let user = $state({user: '<placeholder>', email: "placeholder@example.com"})

	onMount(async () => {
        const session = await authClient.getSession()
        user = session.data.user

		const hashVal = await createHash('SHA-256', 'hex').digest(user.email);

		avatarSrc = `https://gravatar.com/avatar/${hashVal} `;
	});
</script>

<div class="navbar bg-base-100 shadow-sm">
	<div class="flex-1">
		<a class="btn btn-ghost text-xl">albums2</a>
	</div>
	<div class="flex gap-2">
		{#if up}
			<a class="btn btn-ghost" href={up.url}>
				<Icon path={up.icon}></Icon>
				{up.description}
			</a>
		{/if}
		{#if collectionId}
			<form class="join" method="POST" action="?/search" use:enhance>
				<input type="hidden" name="collectionId" value={collectionId} />
				<input
					type="text"
					placeholder="Search"
					required
					name="text"
					class="input join-item w-24 md:w-auto"
				/>
				<button class="btn join-item"><Icon path={mdiMagnify} /></button>
			</form>
		{/if}
		{#if menu}
			<div class="dropdown dropdown-end">
				<div tabindex="0" role="button" class="btn avatar btn-circle btn-ghost">
					<div class="w-10 rounded-full">
						<Icon path={mdiMenu} />
					</div>
				</div>
				<ul
					tabindex="-1"
					class="menu dropdown-content z-1 mt-3 w-52 menu-sm rounded-box bg-base-100 p-2 shadow"
				>
					<ul class="rounded-t-none bg-base-100 p-2">
						{@render menu()}
					</ul>
				</ul>
			</div>
		{/if}
		<div class="dropdown dropdown-end">
			<div tabindex="0" role="button" class="btn avatar btn-circle btn-ghost">
				<div class="w-10 rounded-full">
					<img alt={user.email} src={avatarSrc} />
				</div>
			</div>
			<ul
				tabindex="-1"
				class="menu dropdown-content z-1 mt-3 w-52 menu-sm rounded-box bg-base-100 p-2 shadow"
			>
				<li>
					<button onclick={() => authClient.logout()}>
						<Icon path={mdiLogout} />
						Log out
					</button>
				</li>
			</ul>
		</div>
	</div>
</div>
