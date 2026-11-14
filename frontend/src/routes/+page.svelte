<script lang="ts">
	import { enhance, type SubmitFunction } from '$app/forms';
	import Icon from '#lib/components/Icon.svelte';
	import CollectionThumbnail from '#lib/components/CollectionThumbnail.svelte';
	import { mdiAlertCircle, mdiSquareOff } from '@mdi/js';
	import NavBar from '#lib/components/NavBar.svelte';
	import log from '#lib/log.js';
	import ConfirmationDialog from '#lib/components/ConfirmationDialog.svelte';

	let { data, form } = $props();
	let confirmationDialog: ConfirmationDialog;

	$inspect(data);

	const confirmCollectionCreation: SubmitFunction = async ({ cancel }) => {
		if (!(await confirmationDialog.show('Create this collection?'))) {
			cancel();
		}
	};
</script>

<svelte:head>
	<title>albums2: Collections</title>
</svelte:head>

<NavBar user={data.user} />

<div class="mx-4 overflow-auto pt-4 pb-20">
	{#if data.collections.length === 0}
		<div class="mx-auto my-12 max-h-96 max-w-96 opacity-60">
			<Icon path={mdiSquareOff} />
			There is no collection available.
		</div>
	{:else}
		<div class="flex flex-wrap justify-evenly gap-2">
			{#each data.collections as collection}
				<CollectionThumbnail {collection} />
			{/each}
		</div>
	{/if}
</div>

{log.debug(data.user?.role)}

{#if data.user?.role === 'admin'}
	<div class="divider">Or! Create a new collection.</div>

	<div
		class="m-4 mx-auto mt-8 max-w-96 rounded-box border-2 border-base-300 bg-base-200 p-4 shadow-md"
	>
		<div class="text-xs tracking-wide opacity-60">New collection</div>

		{#if form?.message}
			<div role="alert" class="alert alert-error">
				<Icon path={mdiAlertCircle} />
				<span>{form.message}</span>
			</div>
		{/if}

		<form method="POST" action="?/create" use:enhance={confirmCollectionCreation}>
			<fieldset class="fieldset">
				<legend class="fieldset-legend">collection</legend>
				<input
					type="text"
					class="input w-full"
					required
					name="name"
					placeholder="collection name"
				/>
				<textarea class="textarea h-24 w-full" name="description" placeholder="description"
				></textarea>

				<input
					type="text"
					class="input w-full"
					name="prefix"
					required
					placeholder="S3 Object Prefix"
				/>

				<button class="btn btn-neutral">Add</button>
			</fieldset>
		</form>
	</div>
{/if}

<ConfirmationDialog bind:this={confirmationDialog} />
