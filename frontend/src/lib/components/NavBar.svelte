<script lang="ts">
	import {
		mdiCircleOffOutline,
		mdiClose,
		mdiDotsVertical,
		mdiLogout,
		mdiMagnify,
		mdiShieldCrown
	} from '@mdi/js';
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
	let user = $state({ user: '<placeholder>', email: 'placeholder@example.com', role: 'user' });

	onMount(async () => {
		const session = await authClient.getSession();
		user = session.data.user;

		const hashVal = await createHash('SHA-256', 'hex').digest(user.email);

		avatarSrc = `https://gravatar.com/avatar/${hashVal} `;
	});

	let searchDlg: HTMLDialogElement;
</script>

<div class="navbar bg-base-100 shadow-sm">
	<div class="flex-1">
		<a class="btn btn-ghost text-xl">albums2</a>
	</div>
	<div class="flex gap-2">
		{#if up}
			<a class="btn rounded-full btn-ghost max-lg:btn-circle" href={up.url}>
				<Icon path={up.icon}></Icon>
				<span class="hidden lg:block">{up.description}</span>
			</a>
		{/if}
		{#if collectionId}
			<button class="bth-circle btn btn-ghost lg:hidden" onclick={() => searchDlg.showModal()}>
				<Icon path={mdiMagnify} />
			</button>
			<form
				class="join max-lg:hidden"
				method="POST"
				action="/collections/{collectionId}/search?/search"
				use:enhance
			>
				<input type="hidden" name="collectionId" value={collectionId} />
				<input
					type="text"
					placeholder="Search"
					required
					name="text"
					class="input join-item w-24 rounded-l-full md:w-auto"
				/>

				<button class="btn join-item rounded-r-full"><Icon path={mdiMagnify} /></button>
			</form>
		{/if}

		<div class="dropdown dropdown-end">
			<div tabindex="0" role="button" class="btn avatar btn-circle btn-ghost">
				<div class="rounded-full">
					<Icon path={mdiDotsVertical} />
				</div>
			</div>
			<ul
				tabindex="-1"
				class="menu dropdown-content z-1 mt-3 w-52 menu-sm rounded-box bg-base-100 p-2 shadow"
			>
				<ul class="rounded-t-none bg-base-100 p-2">
					{#if menu}
						{@render menu()}
					{:else}
						<li class="menu-disabled">
							<button class="disabled">
								<Icon path={mdiCircleOffOutline} />
								There is no actionable item.
							</button>
						</li>
					{/if}
				</ul>
			</ul>
		</div>

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
				{#if user.role}
					<li>
						<a href="/admin"><Icon path={mdiShieldCrown} />Administrator</a>
					</li>
					<li></li>
				{/if}
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

<dialog class="modal modal-bottom sm:modal-middle" bind:this={searchDlg}>
	<div class="modal-box">
		<h3 class="text-lg font-bold">Search</h3>
		<form method="POST" action="/collections/{collectionId}/search?/search" use:enhance>
			<input type="text" placeholder="Search" required name="text" class="input w-full" />
			<input type="hidden" name="collectionId" value={collectionId} />
			<div class="modal-action justify-end gap-2">
				<button type="submit" class="btn btn-primary">
					<Icon path={mdiMagnify} />Search
				</button>
				<button type="button" class="btn" onclick={() => searchDlg.close()}>
					<Icon path={mdiClose} />Close
				</button>
			</div>
		</form>
	</div>
</dialog>
