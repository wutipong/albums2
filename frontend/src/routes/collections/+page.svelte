<script lang="ts">
	import { enhance } from '$app/forms';
	import Icon from '$lib/components/Icon.svelte';
	import { mdiAlertCircle, mdiCircleOffOutline, mdiLibraryShelves } from '@mdi/js';

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

<ul class="list mx-auto mt-16 max-w-96 rounded-box border-2 border-base-300 bg-base-300 shadow-md">
	<li class="p-4 pb-2 text-xs tracking-wide opacity-60">Available collections</li>
	{#if data.collections.length === 0}
		<li class="list-row opacity-60">
			<Icon path={mdiCircleOffOutline} size={1.5} />
			<div>
				<div>No Collection Available</div>
				<div class="text-xs font-semibold">An administrator can add a new collection.</div>
			</div>
		</li>
	{/if}
	{#each data.collections as c}
		<li class="list-row">
			<Icon path={mdiLibraryShelves} size={1.5} />
			<a href="/collections/{c._id}">
				<div>{c.name}</div>
				<div class="text-xs font-semibold opacity-60">{c.description}</div>
			</a>
		</li>
	{/each}
</ul>

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
