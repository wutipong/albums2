<script lang="ts">
	import Icon from '$lib/components/Icon.svelte';
	import MediaThumbnail from '$lib/components/MediaThumbnail.svelte';
	import { mdiInformationOutline, mdiImageOff, mdiLibraryShelves, mdiBookshelf } from '@mdi/js';

	const { data } = $props();

	let media = $derived(data.media);
	let album = $derived(data.album);

	$inspect(data);
</script>

<div class="navbar bg-base-100 shadow-sm">
	<div class="flex-1">
		<a class="btn btn-ghost text-xl">albums2</a>
	</div>
	<div class="flex-none">
		<ul class="menu menu-horizontal px-1">
			<li>
				<a href="/collections/{album.collectionId}">
					<Icon path={mdiBookshelf} />
					<span class="hidden sm:block">Albums</span>
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
<div class="h-12 w-full bg-base-300 p-3 flex flex-row gap-1">
	<Icon path={mdiInformationOutline} />
	<div>{album.name}</div>
</div>

{#if media.docs.length === 0}
	<div class="mx-auto my-12 max-h-96 max-w-96 opacity-60">
		<Icon path={mdiImageOff} />
		This album contains no media.
	</div>
{/if}
{#each media.docs as doc}
	<MediaThumbnail media={doc} />
{/each}
