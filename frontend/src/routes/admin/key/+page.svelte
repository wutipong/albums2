<script lang="ts">
	import { authClient } from '$lib/auth';
	import { mdiAlert, mdiInformation, mdiMinus, mdiPlus } from '@mdi/js';
	import Icon from '$lib/components/Icon.svelte';
	import { onMount } from 'svelte';
	import type { ApiKey } from '@better-auth/api-key/types';
	import Toast from '$lib/components/Toast.svelte';

	let apiKeys: Omit<ApiKey, 'key'>[] = $state([]);
	let apiKeyModal: HTMLDialogElement;
	let apiNewKey = $state('');

	let { data } = $props();

	let toast: Toast;

	onMount(async () => {
		const keys = await authClient.apiKey.list();
		if (keys.data) {
			apiKeys = keys.data.apiKeys;
		}
	});

	async function addNewApiKey() {
		const { data, error } = await authClient.apiKey.create({});
		if (error) {
			toast.add(`fail to create new key: ${error.statusText}`, 'error');
			return;
		}
		if (data) {
			apiNewKey = data.key;
			apiKeyModal.showModal();
			navigator.clipboard.writeText(apiNewKey);
		}

		const keys = await authClient.apiKey.list();
		if (keys.data) {
			apiKeys = keys.data.apiKeys;
		}
	}

	async function deleteApiKey(id: string) {
		const { error } = await authClient.apiKey.delete({ keyId: id });
		if (error) {
			toast.add(`fail to delete key: ${error.statusText}`, 'error');
			return;
		}

		const keys = await authClient.apiKey.list();
		if (keys.data) {
			apiKeys = keys.data.apiKeys;
		}
	}
</script>

<svelte:head>
	<title>Albums: User</title>
</svelte:head>

<div class="w-full p-4 pt-8">
	<h2 class="p-4 text-2xl font-semibold">API Keys</h2>
	<table class="table-pin-rows table table-zebra">
		<thead>
			<tr>
				<th>Key</th>
				<th>Created At</th>
				<th>Action</th>
			</tr>
		</thead>
		<tbody>
			{#each apiKeys as apiKey}
				<tr>
					<td>{apiKey.start}...</td>
					<td>{apiKey.createdAt.toLocaleString()}</td>
					<td>
						<button class="btn btn-sm" onclick={() => deleteApiKey(apiKey.id)}>
							<Icon path={mdiMinus} /> Delete
						</button>
					</td>
				</tr>
			{/each}
		</tbody>
		<tfoot>
			<tr>
				<th></th>
				<th></th>
				<th>
					<button class="btn btn-primary btn-sm" onclick={() => addNewApiKey()}>
						<Icon path={mdiPlus} /> Add
					</button>
				</th>
			</tr>
		</tfoot>
	</table>
</div>

<dialog class="modal" bind:this={apiKeyModal}>
	<div class="modal-box">
		<h3 class="text-lg font-bold">New key added</h3>
		<div class="my-1 alert alert-info">
			<Icon path={mdiInformation} />
			<span>The key has been copied to the clipboard.</span>
		</div>
		<div role="alert" class="my-1 alert alert-warning">
			<Icon path={mdiAlert} />
			<span> This API key will never be visible again! </span>
		</div>
		<p class="py-4 font-mono text-wrap break-all">{apiNewKey}</p>
		<div class="modal-action">
			<form method="dialog">
				<!-- if there is a button in form, it will close the modal -->
				<button class="btn">Close</button>
			</form>
		</div>
	</div>
</dialog>

<Toast bind:this={toast}></Toast>