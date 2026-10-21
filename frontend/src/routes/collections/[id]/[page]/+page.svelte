<script lang="ts">
	import { enhance } from '$app/forms';
	import AlbumInfoDialog from '$lib/components/AlbumInfoDialog.svelte';
	import AlbumThumbnail from '$lib/components/AlbumThumbnail.svelte';
	import CollectionInfoDialog from '$lib/components/CollectionInfoDialog.svelte';
	import Icon from '$lib/components/Icon.svelte';
	import log from '$lib/log';
	import Pagination from './Pagination.svelte';
	import {
		mdiBlender,
		mdiFeatureSearchOutline,
		mdiInformationOutline,
		mdiLibraryShelves,
		mdiMagnify,
		mdiSquareOff
	} from '@mdi/js';

	const { data } = $props();
	const collection = $derived(data.collection);
	const albums = $derived(data.albums);

	$inspect(data);

	async function createRenditionJobs() {
		log.debug('request renditions.');
		const resp = await fetch(`/api/collections/${data.collection._id}/rendition-jobs`, {
			method: 'POST'
		});

		const obj = await resp.json();

		return true;
	}

	async function createEmbeddingJobs() {
		log.debug('request renditions.');
		const resp = await fetch(`/api/collections/${data.collection._id}/embedding-jobs`, {
			method: 'POST'
		});

		const obj = await resp.json();

		return true;
	}

	let collectionInfoDlg: CollectionInfoDialog;
</script>

<div class="navbar bg-base-100 shadow-sm">
	<div class="flex-1">
		<a class="btn btn-ghost text-xl">albums2</a>
	</div>
	<div class="flex gap-2">
		<a class="btn btn-ghost" href="/collections/">
			<Icon path={mdiLibraryShelves}></Icon>
			Collections
		</a>
		<form class="join" method="POST" action="?/search" use:enhance>
			<input type="hidden" name="collectionId" value={data.collection._id.toString()} />
			<input
				type="text"
				placeholder="Search"
				required
				name="text"
				class="input join-item w-24 md:w-auto"
			/>
			<button class="btn join-item"><Icon path={mdiMagnify} /></button>
		</form>
		<div class="dropdown dropdown-end">
			<div tabindex="0" role="button" class="btn avatar btn-circle btn-ghost">
				<div class="w-10 rounded-full">
					<img
						alt="Tailwind CSS Navbar component"
						src="https://img.daisyui.com/images/stock/photo-1534528741775-53994a69daeb.webp"
					/>
				</div>
			</div>
			<ul
				tabindex="-1"
				class="menu dropdown-content z-1 mt-3 w-52 menu-sm rounded-box bg-base-100 p-2 shadow"
			>
				<li>
					<button onclick={() => collectionInfoDlg.show(data.collection)}>
						<Icon path={mdiInformationOutline} />Collection information
					</button>
				</li>
				<li></li>
				<li>
					<button onclick={createRenditionJobs}>
						<Icon path={mdiBlender} /> Process media
					</button>
				</li>
				<li>
					<button onclick={createEmbeddingJobs}>
						<Icon path={mdiFeatureSearchOutline} />
						Generate media search metadata
					</button>
				</li>
			</ul>
		</div>
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

<CollectionInfoDialog bind:this={collectionInfoDlg} />

<Pagination currentPage={albums.page} totalPage={albums.totalPages} />
