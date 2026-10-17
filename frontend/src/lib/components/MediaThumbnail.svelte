<script lang="ts">
	import notAvailableSvg from '$lib/assets/not-available.svg?raw';
	import { mdiAlertCircle, mdiClock, mdiFilmstrip, mdiMovie } from '@mdi/js';
	import Icon from '$lib/components/Icon.svelte';
	import log from '$lib/log'

	let { media, onclick = (id: string) => {} } = $props();

	const TARGET_HEIGHT = 300;

	let preview = $state(false);
	let thumbnailWidth = $derived((TARGET_HEIGHT * media.thumbnail_width) / media.thumbnail_height);
	let thumbnailHeight = TARGET_HEIGHT;
	let assetType = $derived(media.type);
	let imageFrames = $derived(media.image_frames);

	let thumbnailLoading = $state(true);
	let previewLoading = $state(true);

</script>

<button
	tabindex="0"
	class={`block cursor-pointer overflow-hidden rounded-xl p-4 
		hover:bg-base-content/20 hover:shadow-xl`}
	onmouseenter={() => (preview = true)}
	onmouseleave={() => (preview = false)}
	onclick={() => {
		log.debug(media, 'thumbnail clicked.')
		onclick(media);
	}}
	disabled={media.processStatus != 'processed'}
>
	<div class="relative h-full w-full">
		<div
			class:hidden={preview}
			class="box-border h-full w-full overflow-hidden rounded-xl border-1 border-base-content/20"
			style={`width: ${thumbnailWidth}px; height: ${thumbnailHeight}px;`}
		>
			{#if media.thumbnail_url === undefined || media.thumbnail_url === ''}
				{@html notAvailableSvg}
			{:else}
				<img
					width={media.thumbnailWidth}
					height={media.thumbnailHeight}
					src={media.thumbnail_url}
					alt={media.id}
					class="h-full w-full"
					onload={() => (thumbnailLoading = false)}
				/>
				<div class="h-full w-full skeleton bg-base-100" class:hidden={!thumbnailLoading}></div>
			{/if}
		</div>

		<div
			class:hidden={!preview}
			class="box-border h-full w-full overflow-hidden rounded-xl border-1 border-base-content/20"
			style={`width: ${thumbnailWidth}px; height: ${thumbnailHeight}px;`}
		>
			{#if media.preview_url === undefined || media.preview_url == ''}
				{@html notAvailableSvg}
			{:else}
				<img
					width={media.thumbnailWidth}
					height={media.thumbnailHeight}
					src={media.preview_url}
					alt={media.id}
					class="h-full w-full"
					onload={() => (previewLoading = false)}
				/>
				<div
					class="h-full w-full skeleton animate-pulse bg-base-200"
					class:hidden={!previewLoading}
				></div>
			{/if}
		</div>

		<div class="absolute top-1 right-2 grid grid-cols-1 place-items-end gap-2">
			{#if assetType === 'video'}
				<div class="badge badge-soft badge-info"><Icon path={mdiMovie} /> Video</div>
			{/if}

			{#if imageFrames > 1}
				<div class="badge badge-soft badge-info"><Icon path={mdiFilmstrip} /> Animation</div>
			{/if}

			{#if media.processStatus == 'error'}
				<div class="badge badge-soft badge-error"><Icon path={mdiAlertCircle} /> Failed</div>
			{/if}

			{#if media.processStatus == 'processing'}
				<div class="badge badge-soft badge-warning"><Icon path={mdiClock} /> Processing</div>
			{/if}
		</div>
	</div>
</button>
