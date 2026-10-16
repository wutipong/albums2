<script lang="ts">
	import AlbumThumbnail from '$lib/components/AlbumThumbnail.svelte';
	import Icon from '$lib/components/Icon.svelte';
	import Pagination from './Pagination.svelte';
	import { mdiInformationOutline, mdiLibraryShelves, mdiSquareOff } from '@mdi/js';

	const { data } = $props();
	const collection = $derived(data.collection);
	const albums = $derived(data.albums);

	$inspect(data);
</script>

<div class="navbar bg-base-100 shadow-sm">
	<div class="flex-1">
		<a class="btn btn-ghost text-xl">albums2</a>
	</div>
	<div class="flex-none">
		<ul class="menu menu-horizontal px-1">
			<li>
				<a href="/collections">
					<Icon path={mdiLibraryShelves} />
					<span class="hidden sm:block">Collections</span>
				</a>
			</li>
			<li>
				<details>
					<summary>Parent</summary>
					<ul class="rounded-t-none bg-base-100 p-2">
						<li><a>Link 1</a></li>
						<li><a>Link 2</a></li>
					</ul>
				</details>
			</li>
		</ul>
	</div>
</div>
<div class="flex h-12 w-full flex-row gap-2 bg-base-300 p-3">
	<Icon path={mdiInformationOutline} />
	<span>{collection.name}</span>
</div>

<div class="mx-4 overflow-auto pt-4 pb-20">
	{#if albums.docs.length === 0}
		<div class="mx-auto my-12 max-h-96 max-w-96 opacity-60">
			<Icon path={mdiSquareOff} />
			This collection contains no albums.
		</div>
	{:else}
		<div class="flex flex-wrap justify-evenly gap-2">
			{#each albums.docs as album}
				<AlbumThumbnail {album} />
			{/each}
		</div>
	{/if}
</div>

<Pagination currentPage={albums.page} totalPage={albums.totalPages}/>