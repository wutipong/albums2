<script lang="ts">
	import AlbumThumbnail from '#lib/components/AlbumThumbnail.svelte';
	import CollectionInfoDialog from '#lib/components/CollectionInfoDialog.svelte';
	import Icon from '#lib/components/Icon.svelte';
	import NavBar from '#lib/components/NavBar.svelte';
	import Toast from '#lib/components/Toast.svelte';
	import ConfirmationDialog from '#lib/components/ConfirmationDialog.svelte';
	import log from '#lib/log.js';
	import Pagination from './Pagination.svelte';
	import {
		mdiBlender,
		mdiFeatureSearchOutline,
		mdiInformationOutline,
		mdiLibraryShelves,
		mdiBookshelf,
		mdiSquareOff,
		mdiSync
	} from '@mdi/js';

	const { data } = $props();
	const collection = $derived(data.collection);
	const albums = $derived(data.albums);

	$inspect(data);

	async function createSyncJobs() {
		if (!(await confirmationDialog.show('Synchronize this collection?'))) {
			return;
		}

		log.debug('request sync jobs.');
		const resp = await fetch(`/api/collections/${data.collection._id}/sync-jobs`, {
			method: 'POST'
		});

		const obj = await resp.json();

		toast.add('Collection synchronization started.', 'info');
		return true;
	}

	async function createEmbeddingJobs() {
		if (
			!(await confirmationDialog.show('Generate search metadata for all media in this collection?'))
		) {
			return;
		}

		log.debug('request renditions.');
		const resp = await fetch(`/api/collections/${data.collection._id}/embedding-jobs`, {
			method: 'POST'
		});

		const obj = await resp.json();

		toast.add('Media search metadata generation started.', 'info');
		return true;
	}

	let collectionInfoDlg: CollectionInfoDialog;
	let toast: Toast;
	let confirmationDialog: ConfirmationDialog;
</script>

<svelte:head>
	<title>albums2: {data.collection.name}</title>
</svelte:head>

{#snippet menu()}
	<li>
		<button onclick={() => collectionInfoDlg.show(data.collection)}>
			<Icon path={mdiInformationOutline} />Collection information
		</button>
	</li>
	{#if data.user?.role == 'admin'}
		<li></li>
		<li class="menu-item">Administration actions</li>
		<li>
			<button onclick={createSyncJobs}>
				<Icon path={mdiSync} /> Synchronize
			</button>
		</li>
		<li>
			<button onclick={createEmbeddingJobs}>
				<Icon path={mdiFeatureSearchOutline} />
				Generate media search metadata
			</button>
		</li>
	{/if}
{/snippet}

<NavBar
	user={data.user}
	up={{ description: 'Collections', icon: mdiLibraryShelves, url: '/' }}
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
<Toast bind:this={toast} />
<ConfirmationDialog bind:this={confirmationDialog} />

<Pagination currentPage={albums.page} totalPage={albums.totalPages} />
