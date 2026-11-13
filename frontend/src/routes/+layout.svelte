<script lang="ts">
	import './layout.css';
	import { pwaAssetsHead } from 'virtual:pwa-assets/head';

	let { children } = $props();
</script>

<svelte:head>
	{#each pwaAssetsHead.links as link}
		<link {...link} />
	{/each}
	{#if pwaAssetsHead.themeColor}
		<meta name={pwaAssetsHead.themeColor.name} content={pwaAssetsHead.themeColor.content} />
	{/if}
</svelte:head>
{@render children()}
