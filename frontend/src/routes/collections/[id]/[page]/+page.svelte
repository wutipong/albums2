<script lang="ts">
	import AlbumThumbnail from '$lib/components/AlbumThumbnail.svelte';
	import Icon from '$lib/components/Icon.svelte';
	import { mdiInformationOutline, mdiLibraryShelves } from '@mdi/js';

	const { data } = $props();

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
<div class="h-12 w-full bg-base-300 p-3">
	<Icon path={mdiInformationOutline} />
</div>

{#each data.docs as album}
	<AlbumThumbnail {album} />
{/each}
