<script lang="ts">
	import Icon from '$lib/components/Icon.svelte';
	import MediaThumbnail from '$lib/components/MediaThumbnail.svelte';
	import {
		mdiInformationOutline,
		mdiImageOff,
		mdiBookshelf,
		mdiFolderMultipleImage,
		mdiMagnify,

		mdiLoading

	} from '@mdi/js';
	import { untrack } from 'svelte';
	import MediaViewer from '$lib/components/MediaViewer.svelte';
	import { enhance } from '$app/forms';
	import log from '$lib/log'

	const { data } = $props();
	let collectionId = $derived(data.collectionId)
	$inspect(data);

	let mediaItems = $state(untrack(() => data.media.docs));
	let currentPage = $state(untrack(() => data.media));
	let index = $state(0);
	let current = $derived(mediaItems[index]);
	let showViewer = $state(false);
	let isLoading = $state(false);

	function next() {
		if (index == mediaItems.length - 1) {
			return;
		}

		index = index + 1;
	}

	function previous() {
		if (index === 0) {
			return;
		}
		index = index - 1;
	}

	async function createCollectionCover() {
		const resp = await fetch(`/api/media/${current._id}/collectioncover`, {
			method: 'POST'
		});

		const obj = await resp.json();

		return true;
	}

	async function loadNextPage() {
		if (currentPage.page == currentPage.totalPages) {
			return;
		}

		isLoading = true;
		const resp = await fetch(`/api/collections/${data.collectionId}/search?search=${data.search}&page=${currentPage.page + 1}`);
		currentPage = await resp.json();
		log.debug(currentPage);

		mediaItems = [...mediaItems, ...currentPage.docs];

		isLoading = false;
	}

	function observeVisibility(node: Element, options = { threshold: 0.1 }) {
		const observer = new IntersectionObserver(([entry]) => {
			if (entry.isIntersecting) {
				loadNextPage();
			}
		}, options);

		observer.observe(node);

		return {
			destroy() {
				observer.disconnect();
			}
		};
	}
</script>

<div class="relative flex h-screen w-screen flex-col">
	<div class="navbar bg-base-100 shadow-sm">
		<div class="flex-1">
			<a class="btn btn-ghost text-xl">albums2</a>
		</div>
		<div class="flex gap-2">
			<a class="btn btn-ghost" href="/collections/{data.collectionId}">
				<Icon path={mdiBookshelf}></Icon>
				Albums
			</a>
			<form class="join" method="POST" action="?/search" use:enhance>
				<input type="hidden" name="collectionId" value={data.collectionId} />
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
					<!-- <li>
						<button onclick={createRenditionJobs}>
							<Icon path={mdiBlender} />
							Process media
						</button>
					</li>
					<li>
						<button onclick={createEmbeddingJobs}>
							<Icon path={mdiFeatureSearchOutline} />
							Generate media search metadata
						</button>
					</li> -->
				</ul>
			</div>
		</div>
	</div>
	<div class="flex h-12 w-full flex-row gap-1 bg-base-300 p-3">
		<Icon path={mdiMagnify} />
		<div>{data.search}</div>
	</div>

	<div class="mx-4 overflow-auto pt-8">
		{#if mediaItems.length === 0}
			<div class="mx-auto my-12 max-h-96 max-w-96 opacity-60">
				<Icon path={mdiImageOff} />
				This search contains no results.
			</div>
		{:else}
			<div class="flex flex-wrap justify-evenly gap-1">
				{#each mediaItems as doc, i (doc._id)}
					<MediaThumbnail
						media={doc}
						onclick={(asset: any) => {
							index = i;
							showViewer = true;
						}}
					/>
				{/each}
			</div>
			<button
				class="btn w-full btn-neutral"
				class:hidden={currentPage.page == currentPage.totalPages}
				class:disabled={isLoading}
				onclick={loadNextPage}
				use:observeVisibility
			>
				<Icon path={mdiLoading} spin={isLoading}></Icon>
				Load More
			</button>
		{/if}
	</div>

	{#if current}
		<MediaViewer
			assetType={current.type}
			viewURL={current.viewUrl}
			filename={current.name}
			bind:show={showViewer}
			{next}
			{previous}
			hasNext={index != mediaItems.length - 1}
			hasPrevious={index != 0}
			menu={viewMenu}
		/>
	{/if}
</div>

{#snippet viewMenu()}
	<li>
		<button
			onclick={() => {
				// assetInfoDialog.show(asset);
			}}
		>
			<Icon path={mdiInformationOutline} /> Media information
		</button>
	</li>
	<!-- li>
		<button
			disabled={asset.copy_url == undefined || asset.copy_url == ''}
			onclick={() => {
				copyImageToClipboard(asset.copy_url);
				toast.add('Image copied to clipboard', 'success');
			}}
		>
			<Icon path={mdiClipboardOutline} /> Copy to clipboard
		</button>
	</li>
	<li>
		<a href={asset.original_url} target="_blank">
			<Icon path={mdiDownload} /> Download
		</a>
	</li-->
	{#if /*data.user.role == 'admin' ||*/ true}
		<li></li>
		<li class="menu-item">Administration actions</li>
		<li>
			<button
				onclick={() => {
					createCollectionCover();
				}}
			>
				<Icon path={mdiFolderMultipleImage} /> Set as collection cover
			</button>
		</li>
	{/if}
{/snippet}
