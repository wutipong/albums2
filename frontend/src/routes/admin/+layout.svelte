<script lang="ts">
	import NavBar from './NavBar.svelte';
	import { mdiImage, mdiImageAlbum, mdiKeyChain, mdiTrayFull } from '@mdi/js';
	import Icon from '$lib/components/Icon.svelte';
	import { page } from '$app/state';

	let { data, children } = $props();
	let drawerToggle: HTMLInputElement;
</script>

<NavBar
    user={data.user}
	onMenuBtn={() => {
		drawerToggle.checked = !drawerToggle.checked;
	}}
/>
<div class="drawer lg:drawer-open">
	<input id="my-drawer-3" type="checkbox" class="drawer-toggle" bind:this={drawerToggle} />
	<div class="drawer-content">
		{@render children()}
	</div>
	<div class="drawer-side">
		<label for="my-drawer-3" aria-label="close sidebar" class="drawer-overlay"></label>
		<ul class="menu min-h-full w-80 bg-base-200 p-4">
			<li><h2 class="menu-title">Administration Menu</h2></li>
			<li class:menu-active={page.url.pathname === '/admin/albums'}>
				<a href="/admin/albums">
					<Icon path={mdiImageAlbum} />Albums
				</a>
			</li>

			<!-- <li class:menu-active={page.url.pathname === '/admin/asset'}>
				<a href="/admin/asset">
					<Icon path={mdiImage} /> Assets
				</a>
			</li>

			<li class:menu-active={page.url.pathname === '/admin/queue'}>
				<a href="/admin/queue">
					<Icon path={mdiTrayFull} /> Processing Queue
				</a>
			</li> -->
			<li class:menu-active={page.url.pathname === '/admin/key'}>
				<a href="/admin/key">
					<Icon path={mdiKeyChain} /> API Keys
				</a>
			</li>
		</ul>
	</div>
</div>