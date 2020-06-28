<script>
	export let segment;
</script>

<script context="module">
	import { waitLocale } from 'svelte-i18n'
	
	export async function preload() {
		// Wait for locales to load
		return waitLocale()
	}
</script>

<style>
</style>

<svelte:head>
	<title>{segment || ""} Ewen Le Bihan</title>
</svelte:head>

<main>
	<slot></slot>
</main>
