<script lang="ts">
	const { data } = $props();
</script>

{#each data.docs as doc}
	<a class='btn btn-link' href='/media/{doc._id}'>{doc.name}</a>
{/each}

