<script lang="ts">
	import AlbumThumbnail from "$lib/components/AlbumThumbnail.svelte";

	const { data } = $props();

	$inspect(data)
</script>

{#each data.docs as album}
	<AlbumThumbnail {album} />
{/each}

