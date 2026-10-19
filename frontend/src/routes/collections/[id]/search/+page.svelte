<script lang='ts'>
    const { data} = $props()
</script>

<div class="text text-xl">{data.id}</div>
<div class="text text-xl">{data.search}</div>
<ol class="list-decimal list-outside pl-32">
    {#each data.embedding.embedding as f}
        <li>{f}</li>
    {/each}
</ol>