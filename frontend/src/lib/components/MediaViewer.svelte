<script lang="ts">
	import { mdiChevronLeft, mdiChevronRight, mdiClose, mdiDotsVertical } from '@mdi/js';
	import Hammer from 'hammerjs';
	import Icon from '$lib/components/Icon.svelte';
	import '@videojs/html/video/player';
	import '@videojs/html/video/skin';

	let {
		viewURL = 'http://example.com',
		assetType = 'image',
		filename = '',
		show = $bindable(false),
		next,
		previous,
		hasNext = false,
		hasPrevious = false,
		menu
	} = $props();

	let loading = $state(true);

	$effect(() => {
		loading = true;
	});

	function doNext() {
		next();
	}

	function doPrevious() {
		previous();
	}

	function doClose() {
        assetType = 'image'
		show = false;
	}

	function handleKeyDown(event: KeyboardEvent) {
		if (!show) return;
		if (event.key === 'Escape') {
			doClose();
		} else if (event.key === 'ArrowRight' && hasNext) {
			if (hasNext) {
				doNext();
			}
		} else if (event.key === 'ArrowLeft' && hasPrevious) {
			if (hasPrevious) {
				doPrevious();
			}
		}
	}

	function hammerJsAttachment(element: HTMLElement) {
		let manager = new Hammer.Manager(element);
		let swipe = new Hammer.Swipe();
		manager.add(swipe);
		manager.on('swipeleft', () => {
			if (hasNext) {
				doNext();
			}
		});

		manager.on('swiperight', () => {
			if (hasPrevious) {
				doPrevious();
			}
		});
	}
</script>

<svelte:window onkeydown={handleKeyDown} />

<div
	role="presentation"
	class="absolute top-0 right-0 bottom-0 left-0 backdrop-blur-lg backdrop-brightness-50"
	class:hidden={!show}
	{@attach hammerJsAttachment}
>
	{#if assetType === 'image'}
		<div class="h-full w-full">
			<img
				class:hidden={loading}
				onload={() => (loading = false)}
				onerror={() => (loading = false)}
				src={viewURL}
				alt={filename}
				class="m-auto h-full w-full object-contain"
			/>
		</div>
	{/if}
	{#if assetType === 'video'}
		{#key viewURL}
			<video-player title={filename} src={viewURL}>
				<video-skin class="block h-full w-full">
					<video class="h-full w-full">
						<source src={viewURL} type="video/mp4" />
					</video>
				</video-skin>
			</video-player>
		{/key}
	{/if}
	<div class="absolute top-1/2 left-4 -translate-y-1/2 bg-transparent">
		<button
			class="btn btn-circle btn-soft btn-lg"
			class:btn-disabled={!hasPrevious}
			onclick={() => {
				doPrevious();
			}}
		>
			<Icon path={mdiChevronLeft} />
		</button>
	</div>
	<div class="absolute top-1/2 right-4 -translate-y-1/2 bg-transparent">
		<button
			class="btn btn-circle btn-soft btn-lg"
			class:btn-disabled={!hasNext}
			onclick={() => {
				doNext();
			}}
		>
			<Icon path={mdiChevronRight} />
		</button>
	</div>
	<div class="absolute top-4 right-4 flex flex-row-reverse gap-4 rounded-full bg-transparent">
		<button
			class="btn btn-circle btn-soft btn-lg"
			onclick={() => {
				doClose();
			}}
		>
			<Icon path={mdiClose} />
		</button>
		{#if menu}
			<button
				class="btn btn-circle btn-soft btn-lg"
				popovertarget="popover-1"
				style="anchor-name:--anchor-1"
			>
				<Icon path={mdiDotsVertical} />
			</button>
			<ul
				class="menu dropdown w-52 rounded-box bg-base-100 shadow-sm"
				popover
				id="popover-1"
				style="position-anchor:--anchor-1"
			>
				{@render menu()}
			</ul>
		{/if}
	</div>
</div>
