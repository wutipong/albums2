<script lang="ts">
	import Icon from '$lib/components/Icon.svelte';
	import MediaThumbnail from '$lib/components/MediaThumbnail.svelte';
	import {
		mdiInformationOutline,
		mdiImageOff,
	} from '@mdi/js';
	import log from '$lib/log'
	import { untrack } from 'svelte';

	const { data } = $props();

	let media = $derived(data.media);
	let album = $derived(data.album);

	$inspect(data);

	let mediaItems = $state(untrack(() =>data.media.docs));
	let currentPage = $state(untrack(()=>media));
	let isLoading = $state(false)

	async function loadNextPage(){
		if (currentPage.page == currentPage.totalPages) {
			return
		}
		
		isLoading = true;
		const resp = await fetch(`/api/albums/${album._id}/${currentPage.page + 1}`)
		currentPage = await resp.json()
		mediaItems = [...mediaItems, ...currentPage.docs]

		isLoading = false;
	}

	async function createRenditionJobs() {
		log.debug('request renditions.')
		const resp = await fetch(`/api/albums/${album._id}/rendition-jobs`, {
			method: 'POST'
		});

		const obj = await resp.json();

		return true
	}
</script>

<div class="navbar bg-base-100 shadow-sm">
	<div class="flex-1">
		<a class="btn btn-ghost text-xl">daisyUI</a>
	</div>
	<div class="flex gap-2">
		<input type="text" placeholder="Search" class="input w-24 md:w-auto" />
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
					<button onclick={createRenditionJobs}> Process media </button>
				</li>
			</ul>
		</div>
	</div>
</div>
<div class="flex h-12 w-full flex-row gap-1 bg-base-300 p-3">
	<Icon path={mdiInformationOutline} />
	<div>{album.name}</div>
</div>

<div class="mx-4 overflow-auto pt-8">
	{#if mediaItems.length === 0}
		<div class="mx-auto my-12 max-h-96 max-w-96 opacity-60">
			<Icon path={mdiImageOff} />
			This album contains no media.
		</div>
	{:else}
		<div class="flex flex-wrap justify-evenly gap-1">
			{#each mediaItems as doc}
				<MediaThumbnail media={doc} />
			{/each}
		</div>
	{/if}
	<button 
		class='btn btn-neutral w-full' 
		class:hidden={currentPage.page == currentPage.totalPages}
		class:disabled={isLoading} 
		onclick={loadNextPage}
		>Load More</button>
</div>
