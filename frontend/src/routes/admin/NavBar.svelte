<script lang="ts">
	import { goto } from '$app/navigation';
	import {
		mdiImageSearch,
		mdiImageAlbum,
		mdiLogout,
		mdiAccount,
		mdiClose,
		mdiShieldCrownOutline,
		mdiMenu
	} from '@mdi/js';
	import Icon from '$lib/components/Icon.svelte';
	import { onMount } from 'svelte';
	import { createHash } from '@better-auth/utils/hash';
	//import logo from '$lib/assets/logo.svg?raw';

	let search = $state('');
	let searchDialog: HTMLDialogElement;

	let { album = null, onMenuBtn = null, user } = $props();
	let avatarSrc = $state('');

	let albumsUrl = $derived.by(() => {
		if (album) {
			return `/album#${album.id}`;
		} else {
			return '/album';
		}
	});

	function doSearch() {
		if (!search) return;
		const url = new URL('/search', location.origin);
		url.searchParams.append('search', search);

		goto(url);
	}

	onMount(async () => {
		const hashVal = await createHash('SHA-256', 'hex').digest(user.email);

		avatarSrc = `https://gravatar.com/avatar/${hashVal} `;
	});
</script>

<div class="navbar shadow-sm">
	<div class="flex flex-1 flex-row">
		<button
			class="btn btn-square btn-ghost lg:hidden"
			class:hidden={onMenuBtn == null}
			onclick={() => onMenuBtn()}
		>
			<Icon path={mdiMenu} />
		</button>

		<div class="ms-2 w-32 text-base-content">
            albums2
			<!-- {@html logo} -->
		</div>
	</div>
	
	
	<div class="dropdown dropdown-end">
		<div tabindex="0" role="button" class="btn avatar btn-circle btn-ghost">
			<div class="w-10 rounded-full">
				<img alt="User" src={avatarSrc} />
			</div>
		</div>
		<ul
			tabindex="-1"
			class="dropdown-content menu z-1 mt-3 w-52 menu-sm rounded-box bg-base-100 p-2 shadow-xl"
		>
			<li>
				<a href="/user"><Icon path={mdiAccount} /> User</a>
			</li>
			<li></li>
			<li>
				<a href="/logout" data-sveltekit-preload-data="off"><Icon path={mdiLogout} /> Logout</a>
			</li>
		</ul>
	</div>
</div>

<dialog class="modal" bind:this={searchDialog}>
	<div class="modal-box">
		<form method="dialog">
			<button class="btn absolute top-2 right-2 btn-circle btn-ghost btn-sm">
				<Icon path={mdiClose} />
				<span class="sr-only">Close</span>
			</button>
		</form>
		<h3 class="text-lg font-bold">Search</h3>
		<div class="py-4">
			<input
				class="input-bordered input w-full"
				type="text"
				placeholder="Search..."
				bind:value={search}
				onkeydown={(e) => {
					if (e.key === 'Enter') {
						doSearch();
						searchDialog.close();
					}
				}}
			/>
		</div>
	</div>
</dialog>