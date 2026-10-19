<script lang='ts'>
    const { data} = $props()

    $inspect(data)
</script>

<div class="text text-xl">{data.id}</div>
<div class="text text-xl">{data.search}</div>