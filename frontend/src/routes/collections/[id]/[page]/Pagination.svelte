<script lang="ts">
	import { resolve } from '$app/paths';
	import Icon from '#lib/components/Icon.svelte';
	import { mdiPageFirst, mdiPageLast, mdiDotsVertical } from '@mdi/js';
	import { goto } from '$app/navigation';
	import PaginationDialog from './PaginationDialog.svelte';
	import type { RouteId } from '$app/types';
	import { page } from '$app/state';

	interface Props {
		currentPage?: number;
		totalPage?: number;
		pageToShow?: number;
	}

	let { currentPage = 1, totalPage = 1, pageToShow = 5 }: Props = $props();

	let first = $derived.by(() => Math.max(1, currentPage - Math.floor(pageToShow / 2)));
	let last = $derived.by(() => Math.min(totalPage, first + pageToShow - 1));

	let pageNumbers: number[] = $derived.by(() =>
		Array.from({ length: last - first + 1 }, (_, i) => first + i)
	);

	let customInput: PaginationDialog;

	function createLink(p: number): string {
		const params = {
			id: page.params.id ?? '',
			page: String(p)
		};

		return resolve('/collections/[id]/[page]', params);
	}
</script>

<div class="fixed inset-x-1/2 bottom-10 hidden md:flex">
	<div class="join -translate-x-1/2 shadow-xl">
		<a class="btn join-item" href={createLink(1)}>
			<Icon path={mdiPageFirst} />
		</a>

		{#each pageNumbers as i (i)}
			<a
				class="btn join-item"
				class:btn-active={currentPage == i}
				class:btn-primary={currentPage == i}
				href={createLink(i)}
			>
				{i}
			</a>
		{/each}

		<button class="btn join-item" onclick={() => customInput.showModal()}>
			<Icon path={mdiDotsVertical} />
		</button>

		<a class="btn join-item" href={createLink(totalPage)}>
			<Icon path={mdiPageLast} />
		</a>
	</div>
</div>

<PaginationDialog bind:this={customInput} {currentPage} {totalPage} {createLink} />
