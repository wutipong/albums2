<script lang="ts">
	import Icon from '$lib/components/Icon.svelte';
	import MediaThumbnail from '$lib/components/MediaThumbnail.svelte';
	import { mdiCircleOffOutline, mdiImageOff } from '@mdi/js';

	const { data } = $props();

	$inspect(data)
</script>

{#if data.docs.length === 0}
<div class='mx-auto my-auto max-w-96 max-h-96'>
	<Icon path={mdiImageOff} />
	<span>This album contains no media.</span>
</div>
{/if}
{#each data.docs as doc}
	<MediaThumbnail media={doc}/>
{/each}

