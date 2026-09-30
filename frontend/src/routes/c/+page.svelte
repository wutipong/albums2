<script lang="ts">
	import { enhance } from '$app/forms';

	let { data, form } = $props();
	let modal: HTMLDialogElement;
</script>

{#if form?.error}
	<div role="alert" class="alert alert-error">
		<svg
			xmlns="http://www.w3.org/2000/svg"
			fill="none"
			viewBox="0 0 24 24"
			class="h-6 w-6 shrink-0 stroke-current"
		>
			<path
				stroke-linecap="round"
				stroke-linejoin="round"
				stroke-width="2"
				d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
			></path>
		</svg>
		<span>{form.error}</span>
	</div>
{/if}

<ul class="list rounded-box bg-base-100 shadow-md">
	{#each data.collections as c}
		<li class="list-row">
			<div>
				<div>{c.name}</div>
			</div>
		</li>
	{/each}

	<li class="list-row">
		<btn class="w-full" onclick={() => modal.showModal()}>Add</btn>
	</li>
</ul>

<dialog bind:this={modal} class="modal">
	<div class="modal-box">
		<h3 class="text-lg font-bold">New Collection</h3>
		<div class="modal-action">
			<form method="POST" action="?/create" use:enhance>
				<input type="text" class="input" required name="name" placeholder="Name" />
				{#if form?.error}
					<p class="mt-1 text-xs text-error">{form.error}</p>
				{/if}
				<button class="btn">Add</button>
			</form>
		</div>
	</div>
</dialog>
