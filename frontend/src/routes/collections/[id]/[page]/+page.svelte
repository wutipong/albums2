<script lang="ts">
	import { enhance } from '$app/forms';
	import AlbumInfoDialog from '$lib/components/AlbumInfoDialog.svelte';
	import AlbumThumbnail from '$lib/components/AlbumThumbnail.svelte';
	import CollectionInfoDialog from '$lib/components/CollectionInfoDialog.svelte';
	import Icon from '$lib/components/Icon.svelte';
	import NavBar from '$lib/components/์NavBar.svelte';
	import log from '$lib/log';
	import Pagination from './Pagination.svelte';
	import {
		mdiBlender,
		mdiFeatureSearchOutline,
		mdiInformationOutline,
		mdiLibraryShelves,
		mdiBookshelf,
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

{#snippet menu()}
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
{/snippet}

<NavBar
	up={{ description: 'Collections', icon: mdiLibraryShelves, url: '/collections' }}
	collectionId={data.collection._id.toString()}
	{menu}
/>

<div class="flex h-12 w-full flex-row gap-2 bg-base-300 p-3">
	<Icon path={mdiBookshelf} />
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
