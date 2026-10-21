<script lang="ts">
	import { mdiClose, mdiLink } from '@mdi/js';
	import Icon from '$lib/components/Icon.svelte';
	import { ObjectId } from 'bson';
	//import { Temporal } from 'temporal-polyfill';

	interface Album {
		_id: ObjectId;
		name: string;
		createdAt: Date;
		modifiedAt: Date;
	}

	interface Collection {
		_id: ObjectId;
		name: string;
	}

	let album: Album = $state({
		_id: new ObjectId(),
		name: '<placeholder>',
		createdAt: new Date(),
		modifiedAt: new Date()
	});

	let collection: Collection = $state({
		_id: new ObjectId(),
		name: '<placeholder>'
	});

	export function show(a: Album, c: Collection) {
		album = a;
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
						<td>{album._id.toString()}</td>
					</tr>
					<tr>
						<th>Name</th>
						<td>{album.name}</td>
					</tr>
					<tr>
						<th>Collection</th>
						<td
							><a class="flex link gap-1" href="/collections/{collection._id.toHexString()}">
								{collection.name}
								<Icon path={mdiLink} />
							</a></td
						>
					</tr>
					<tr>
						<th>Created At</th>
						<td>{album.createdAt.toLocaleString()}</td>
					</tr>
					<tr>
						<th>Modified At</th>
						<td>{album.modifiedAt.toLocaleString()}</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</dialog>
