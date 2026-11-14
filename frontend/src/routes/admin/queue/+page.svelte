<script lang="ts">
	const { data } = $props();

	$inspect(data);
</script>

<svelte:head>
	<title>albums2: Processing queue</title>
</svelte:head>

<div class="w-full p-4 pt-8">
	<h2 class="p-4 text-2xl font-semibold">Processing</h2>
	<div class="h-64 overflow-auto">
		<table class="table table-pin-rows table-zebra">
			<thead>
				<tr>
					<th>Type</th>
					<th>Payload</th>
				</tr>
			</thead>
			<tbody>
				{#each data.processings as t}
					<tr class="text-bold bg-base-200">
						<td>{t.type}</td>
						<td><span class="font-mono">{JSON.stringify(t.payload, null, 2)}</span></td>
					</tr>
				{/each}
				{#each data.tasks as t}
					<tr>
						<td>{t.type}</td>
						<td><span class="font-mono">{JSON.stringify(t.payload, null, 2)}</span></td>
					</tr>
				{/each}
			</tbody>
			<tfoot>
				<tr>
					<th>Total</th>
					<td>{data.taskCount}</td>
				</tr>
			</tfoot>
		</table>
	</div>

	<h2 class="p-4 text-2xl font-semibold">Failed Jobs</h2>
	<div class="h-64 overflow-auto">
		<table class="table table-pin-rows table-zebra">
			<thead>
				<tr>
					<th>Type</th>
					<th>Payload</th>
				</tr>
			</thead>
			<tbody>
				{#each data.failures as t}
					<tr>
						<td>{t.type}</td>
						<td><span class="font-mono">{JSON.stringify(t.payload, null, 2)}</span></td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>
