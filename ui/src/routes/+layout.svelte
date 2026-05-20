<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<nav>
	<a href="/">Home</a>
	<a href="/map">Map</a>
	<a href="/about">About</a>
</nav>

<style>
	nav {
		display: flex;
		justify-content: flex-end;
		gap: 1rem;
		padding: 0.75rem 1.5rem;
		background-color: var(--surface);
		border-bottom: 1px solid var(--border);
	}

	nav a {
		color: var(--text-muted);
		font-size: 0.9rem;
		font-weight: 500;
		transition: color 0.15s;
	}

	nav a:hover {
		color: var(--text);
	}
</style>


<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

{@render children()}
