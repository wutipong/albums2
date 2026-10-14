<script lang="ts">
	import notAvailableSvg from '$lib/assets/not-available.svg?raw';

	interface Props {
		album: any;
	}

	let { album }: Props = $props();
</script>

<div class="block rounded-xl p-4 hover:bg-base-content/20 hover:shadow-xl" id={album.id}>
	<div class="w-[300px] overflow-clip rounded-xl border-1 border-base-content/20">
		<figure class="block h-[300px] w-[200px]">
			<a href={`/albums/${album._id}/`} class="block h-full w-full">
				<div class="relative h-[200px] w-[300px]">
					{#if album.cover_url == ''}
						{@html notAvailableSvg}
					{:else}
						<img
							class="h-[300px] w-[200px] object-cover"
							src={album.cover_url}
							alt={album.id}
							loading="lazy"
						/>
					{/if}
				</div>
			</a>
		</figure>
		<div class="h-[100px] w-full overflow-clip p-4 break-all">
			<a href={`/album/${album.id}/`}>
				{album.name}
			</a>
		</div>
	</div>
</div>
