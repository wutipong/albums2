<script lang="ts">
	import notAvailableSvg from '$lib/assets/not-available.svg?raw';

	interface Props {
		album: any;
	}

	let { album }: Props = $props();
</script>

<div class="block rounded-xl p-4 hover:bg-base-content/20 hover:shadow-xl" id={album._id}>
	<div class="w-[300px] overflow-clip rounded-xl border-1 border-base-content/20">
		<figure class="block h-[450px] w-[300px]">
			<a href="/albums/{album._id}/" class="block h-full w-full">
				<div class="relative h-[450px] w-[300px]">
					{#if album.cover_url == ''}
						{@html notAvailableSvg}
					{:else}
						<img class="h-[450px] w-[300px] object-cover" src={album.cover_url} alt={album.id} />
					{/if}
					<div
						class="absolute bottom-0 left-0 h-[150px] w-full
                            bg-gradient-to-b from-neutral/0 to-neutral/80
                            p-4 pt-12 break-all text-neutral-content"
					>
						{album.name}
					</div>
				</div>
			</a>
		</figure>
	</div>
</div>
