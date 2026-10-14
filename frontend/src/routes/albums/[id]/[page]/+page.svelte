<script lang="ts">
	import MediaThumbnail from '$lib/components/MediaThumbnail.svelte';

	const { data } = $props();

	$inspect(data)
</script>

{#each data.docs as doc}
	<MediaThumbnail media={doc}/>
{/each}

