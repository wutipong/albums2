<script lang="ts">
	let dialog: HTMLDialogElement;
	let message = $state('');
	let resolveConfirmation: ((confirmed: boolean) => void) | undefined;

	export function show(prompt: string): Promise<boolean> {
		message = prompt;
		dialog.showModal();

		return new Promise((resolve) => {
			resolveConfirmation = resolve;
		});
	}

	function finish(confirmed: boolean) {
		const resolve = resolveConfirmation;
		resolveConfirmation = undefined;
		dialog.close();
		resolve?.(confirmed);
	}

	function handleClose() {
		const resolve = resolveConfirmation;
		resolveConfirmation = undefined;
		resolve?.(false);
	}
</script>

<dialog class="modal" bind:this={dialog} onclose={handleClose}>
	<div class="modal-box">
		<h3 class="text-lg font-bold">Confirm administrator action</h3>
		<p class="py-4">{message}</p>
		<div class="modal-action">
			<button type="button" class="btn" onclick={() => finish(false)}>NO</button>
			<button type="button" class="btn btn-primary" onclick={() => finish(true)}>YES</button>
		</div>
	</div>
</dialog>
