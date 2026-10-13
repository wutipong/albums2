<script lang="ts">
	const { data } = $props();
</script>

{#each data.docs as album}
	<a class='btn btn-link' href='/albums/{album._id}'>{album.name}</a>
{/each}

