<script>
	export let segment;
</script>

<style>
</style>

<svelte:head>
	<title>{segment || ""} TODO:NAME</title>
</svelte:head>

<main>
	<slot></slot>
</main>
