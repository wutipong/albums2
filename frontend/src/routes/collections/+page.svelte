<script lang="ts">
	import { enhance } from '$app/forms';
	import Icon from '$lib/components/Icon.svelte';
	import CollectionThumbnail from '$lib/components/CollectionThumbnail.svelte'
	import { mdiAlertCircle, mdiCircleOffOutline, mdiLibraryShelves, mdiSquareOff } from '@mdi/js';

	let { data, form } = $props();

	$inspect(data);
</script>

<div class="navbar bg-base-100 shadow-sm">
	<div class="flex-1">
		<a class="btn btn-ghost text-xl">daisyUI</a>
	</div>
	<div class="flex-none">
		<ul class="menu menu-horizontal px-1">
			<li><a>Link</a></li>
			<li>
				<details>
					<summary>Parent</summary>
					<ul class="rounded-t-none bg-base-100 p-2">
						<li><a>Link 1</a></li>
						<li><a>Link 2</a></li>
					</ul>
				</details>
			</li>
		</ul>
	</div>
</div>

<div class="mx-4 overflow-auto pt-4 pb-20">
	{#if data.collections.length === 0}
		<div class="mx-auto my-12 max-h-96 max-w-96 opacity-60">
			<Icon path={mdiSquareOff} />
			This collection contains no albums.
		</div>
	{:else}
		<div class="flex flex-wrap justify-evenly gap-2">
			{#each data.collections as collection}
				<CollectionThumbnail {collection} />
			{/each}
		</div>
	{/if}
</div>

<div class="m-4 mx-auto max-w-96 rounded-box border-2 border-base-300 bg-base-200 p-4 shadow-md">
	<div class="text-xs tracking-wide opacity-60">New collection</div>

	{#if form?.message}
		<div role="alert" class="alert alert-error">
			<Icon path={mdiAlertCircle} />
			<span>{form.message}</span>
		</div>
	{/if}

	<form method="POST" action="?/create" use:enhance>
		<fieldset class="fieldset">
			<legend class="fieldset-legend">collection</legend>
			<input type="text" class="input w-full" required name="name" placeholder="collection name" />
			<textarea class="textarea h-24 w-full" name="description" placeholder="description"
			></textarea>

			<button class="btn btn-neutral">Add</button>
		</fieldset>
	</form>
</div>
