<script lang="ts">
	import { mdiClose, mdiLink } from '@mdi/js';
	import Icon from '$lib/components/Icon.svelte';
	import { ObjectId } from 'bson';
	//import { Temporal } from 'temporal-polyfill';

	interface Collection {
		_id: ObjectId;
		name: string;
		createdAt: Date;
		modifiedAt: Date;
	}

	let collection: Collection = $state({
		_id: new ObjectId(),
		name: '<placeholder>',
		createdAt: new Date(),
		modifiedAt: new Date()
	});

	export function show(c: Collection) {
		collection = c;
		dialog.showModal();
	}

	let dialog: HTMLDialogElement;
</script>

<dialog class="modal" bind:this={dialog}>
	<div class="modal-box">
		<form method="dialog">
			<button class="btn absolute top-2 right-2 btn-circle btn-ghost btn-sm">
				<Icon path={mdiClose} />
				<span class="sr-only">Close</span>
			</button>
		</form>
		<h3 class="text-lg font-bold">Infomation</h3>
		<div class="py-4">
			<table class="table w-full">
				<tbody>
					<tr>
						<th>ID</th>
						<td>{collection._id.toString()}</td>
					</tr>
					<tr>
						<th>Name</th>
						<td>{collection.name}</td>
					</tr>
					<tr>
						<th>Created At</th>
						<td>{collection.createdAt.toLocaleString()}</td>
					</tr>
					<tr>
						<th>Modified At</th>
						<td>{collection.modifiedAt.toLocaleString()}</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</dialog>
