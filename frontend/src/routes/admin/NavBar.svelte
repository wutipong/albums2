<script lang="ts">
	import { mdiLogout, mdiAccount, mdiMenu, mdiLibraryShelves } from '@mdi/js';
	import Icon from '$lib/components/Icon.svelte';
	import { onMount } from 'svelte';
	import { createHash } from '@better-auth/utils/hash';
	//import logo from '$lib/assets/logo.svg?raw';

	let { onMenuBtn = null, user } = $props();
	let avatarSrc = $state('');

	onMount(async () => {
		const hashVal = await createHash('SHA-256', 'hex').digest(user.email);

		avatarSrc = `https://gravatar.com/avatar/${hashVal} `;
	});
</script>

<div class="navbar shadow-sm">
	<div class="flex flex-1 flex-row">
		<button
			class="btn btn-circle btn-square btn-ghost lg:hidden"
			class:hidden={onMenuBtn == null}
			onclick={() => onMenuBtn()}
		>
			<Icon path={mdiMenu} />
		</button>

		<a class="btn btn-ghost text-xl">albums2</a>
	</div>
	<a class="btn rounded-full btn-ghost max-lg:btn-circle" href="/collections">
		<Icon path={mdiLibraryShelves}></Icon>
		<span class="hidden lg:block">Collections</span>
	</a>

	<div class="dropdown dropdown-end">
		<div tabindex="0" role="button" class="btn avatar btn-circle btn-ghost">
			<div class="w-10 rounded-full">
				<img alt="User" src={avatarSrc} />
			</div>
		</div>
		<ul
			tabindex="-1"
			class="menu dropdown-content z-1 mt-3 w-52 menu-sm rounded-box bg-base-100 p-2 shadow-xl"
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
