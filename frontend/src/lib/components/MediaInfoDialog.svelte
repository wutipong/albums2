<script lang="ts">
	import { mdiClose, mdiLink } from '@mdi/js';
	import Icon from '$lib/components/Icon.svelte';
	import { ObjectId } from 'bson';
	//import { Temporal } from 'temporal-polyfill';

	interface Media {
		collectionName: any;
		albumId: any;
		albumName: any;
		collectionId: any;
		_id: ObjectId;
		name: string;
		createdAt: Date;
		modifiedAt: Date;
		type: string;
		videoDuration?: any;
		imageFrames?: number;
		original: string;
		preview: string;
		thumbnail: string;
	}

	let media: Media = $state({
		_id: new ObjectId(),
		name: '<placeholder>',
		createdAt: new Date(),
		modifiedAt: new Date(),
		type: 'image',
		imageFrames: 1,
		original: '',
		preview: '',
		thumbnail: '',

        albumId: new ObjectId(),
        albumName: '<placeholder>',
        collectionId: new ObjectId(),
        collectionName: '<placeholder>'
	});

	export function show(m: Media) {
		media = m;
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
						<td>{media._id.toString()}</td>
					</tr>
					<tr>
						<th>Name</th>
						<td>{media.name}</td>
					</tr>
					<tr>
						<th>Album</th>
						<td><a class="link flex gap-1" href="/albums/{media.albumId.toHexString()}">
                            {media.albumName} 
                            <Icon path={mdiLink}/>
                        </a></td>
					</tr>
					<tr>
						<th>Collection</th>
						<td><a class="link flex gap-1" href="/collections/{media.collectionId.toHexString()}">
                            {media.collectionName}  
                            <Icon path={mdiLink}/>
                        </a></td>
					</tr>
					<tr>
						<th>Created At</th>
						<td>{media.createdAt.toLocaleString()}</td>
					</tr>
					<tr>
						<th>Modified At</th>
						<td>{media.modifiedAt.toLocaleString()}</td>
					</tr>
					<tr>
						<th>Type</th>
						<td>{media.type}</td>
					</tr>
					{#if media.type === 'video'}
						<tr>
							<th>Video Duration</th>
							<td>{media.videoDuration.toLocaleString('en-US', { style: 'digital' })}</td>
						</tr>
					{/if}
					{#if media.imageFrames && media.imageFrames > 1}
						<tr>
							<th>Image Frames</th>
							<td>{media.imageFrames}</td>
						</tr>
					{/if}
					<tr>
						<th>Original </th>
						<td>{media.original}</td>
					</tr>
					<tr>
						<th>Preview </th>
						<td>{media.preview}</td>
					</tr>
					<tr>
						<th>Thumbnail </th>
						<td>{media.thumbnail}</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</dialog>
