<script lang="ts">
	import Icon from '#lib/components/Icon.svelte';
	import { mdiAlertCircle, mdiDelete, mdiPlus } from '@mdi/js';
	import { enhance } from '$app/forms';

	const { data, form } = $props();
	const excludedCoverNames = $derived(data.settings?.album?.excludedCoverNames ?? []);

	var inputText = $state('');
</script>

<svelte:head>
	<title>albums2: Albums administration</title>
</svelte:head>

<div class="m-8 rounded border border-1 border-base-300 bg-base-100 p-4 shadow-xl">
	<h2 class="p-4 text-2xl font-semibold">Cover Image Exclude List</h2>
	<div>
		The list controls the names of the cover images to exclude from the album covers populating
		process. It does not affect the manual selection of cover images.
	</div>

	{#if form?.message}
			<div role="alert" class="alert alert-error">
				<Icon path={mdiAlertCircle} />
				<span>{form.message}</span>
			</div>
		{/if}

	<table class="table table-zebra">
		<thead>
			<tr>
				<th>Name</th>
				<th>Action</th>
			</tr>
		</thead>
		<tbody>
			{#each excludedCoverNames as name}
				<tr>
					<td>{name}</td>
					<td>
						<form method="POST" action="?/delete">
							<input type="hidden" name="name" value={name} />
							<button type="submit" class="btn btn-error btn-sm">
								<Icon path={mdiDelete} />
								Remove
							</button>
						</form>
					</td>
				</tr>
			{/each}
		</tbody>
		<tfoot>
			<tr>
				<th>
					<input
						type="text"
						bind:value={inputText}
						placeholder="Name"
						required
						class="input-bordered input w-full"
					/>
				</th>
				<th>
					<form method="POST" action="?/add" use:enhance class="flex flex-1 flex-row gap-2">
						<input type="hidden" name="name" value={inputText} />
						<button type="submit" class="btn btn-primary btn-sm">
							<Icon path={mdiPlus} />
							Add
						</button>
					</form>
				</th>
			</tr>
		</tfoot>
	</table>
</div>
