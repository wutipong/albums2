<script lang="ts">
	import Icon from '#lib/components/Icon.svelte';
	import {
		mdiPageNext,
		mdiPageFirst,
		mdiChevronLeft,
		mdiChevronRight,
		mdiPageLast,
		mdiPlus,
		mdiMinus
	} from '@mdi/js';
	import { goto } from '$app/navigation';

	let { currentPage = 1, totalPage = 0, createLink = (_n: number) => {} } = $props();
	let customPage = $derived(currentPage);

	let dialog: HTMLDialogElement;

	export function showModal() {
		dialog.showModal();
	}
</script>

<dialog class="modal modal-bottom" bind:this={dialog}>
	<div class="modal-box mx-auto w-full max-w-160">
		<h3 class="text-lg font-bold">Move to page</h3>
		<div class="flex flex-col py-4">
			<div class="join mt-3 flex">
				<button class="btn join-item flex-none" onclick={() => goto(createLink(0))}>
					<Icon path={mdiPageFirst} />
					<div class="text hidden sm:block">First</div></button
				>
				<button
					class="btn join-item flex-none"
					class:btn-disabled={currentPage <= 1}
					onclick={() => goto(createLink(currentPage))}
				>
					<Icon path={mdiChevronLeft} />
					<div class="text hidden sm:block">Previous</div>
				</button>

				<input type="number" readonly class="input join flex-1" value={currentPage} />
				<button
					class="btn join-item flex-none"
					class:btn-disabled={currentPage >= totalPage}
					onclick={() => goto(createLink(currentPage))}
				>
					<Icon path={mdiChevronRight} />
					<div class="text hidden sm:block">Next</div>
				</button>
				<button class="btn join-item" onclick={() => goto(createLink(totalPage - 1))}>
					<Icon path={mdiPageLast} />
					<div class="text hidden sm:block">Last</div>
				</button>
			</div>

			<div class="divider">Enter page number manually.</div>

			<div class="join mt-3 flex">
				<input
					type="number"
					inputmode="numeric"
					class="input join-item flex-1"
					bind:value={customPage}
					placeholder="page #"
					max={totalPage}
					min={1}
				/>
				<button
					class="btn join-item flex-none"
					class:btn-disabled={customPage == 1}
					onclick={() => (customPage = Math.max(customPage - 1, 1))}
				>
					<Icon path={mdiMinus} />
				</button>
				<button
					class="btn join-item flex-none"
					class:btn-disabled={customPage == totalPage}
					onclick={() => (customPage = Math.min(customPage + 1, totalPage))}
				>
					<Icon path={mdiPlus} />
				</button>

				<button class="btn join-item flex-none" onclick={() => goto(createLink(customPage))}>
					<Icon path={mdiPageNext}></Icon>
					<div class="text hidden sm:block">Go</div>
				</button>
			</div>
		</div>
	</div>
	<form method="dialog" class="modal-backdrop">
		<button>close</button>
	</form>
</dialog>
