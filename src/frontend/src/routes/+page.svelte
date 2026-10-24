<script lang="ts">
	let checking = $state(false);
	let status = $state('Not checked');

	async function checkBackend() {
		checking = true;
		status = 'Checking…';
		try {
			const response = await fetch('/api/ready', { signal: AbortSignal.timeout(5000) });
			if (response.ok) {
				status = 'Go API and PostgreSQL are ready.';
			} else if (response.status === 503) {
				status = 'Go API is running. PostgreSQL is unavailable.';
			} else {
				status = 'Cannot reach the Go API. Check that the backend is running.';
			}
		} catch {
			status = 'Cannot reach the Go API. Check that the backend is running.';
		} finally {
			checking = false;
		}
	}
</script>

<svelte:head>
	<title>Portal67</title>
	<meta name="description" content="Portal67 project starter" />
</svelte:head>

<main>
	<h1>Portal67</h1>
	<p>SvelteKit + Vite+, Go, and PostgreSQL.</p>
	<button onclick={checkBackend} disabled={checking}>Check backend connection</button>
	<p role="status">{status}</p>
</main>

<style>
	:global(body) {
		margin: 0;
		font-family: system-ui, sans-serif;
		color: #172033;
		background: #f6f7fb;
	}

	main {
		max-width: 40rem;
		margin: 12vh auto;
		padding: 2rem;
	}

	button {
		padding: 0.75rem 1rem;
		font: inherit;
		cursor: pointer;
	}
</style>
