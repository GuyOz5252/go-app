<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { toast } from '$lib/toast.svelte';

	let { children } = $props();
</script>

<svelte:head>
	<title>Go Chat</title>
	<link rel="icon" href={favicon} />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
</svelte:head>

<div aria-live="polite" class="sr-only">{toast.message ?? ''}</div>

{#if toast.message}
	<div
		class="toast"
		role="status"
		class:error={toast.kind === 'error'}
		class:info={toast.kind === 'info'}
	>
		{toast.message}
	</div>
{/if}

{@render children()}

<style>
	.toast {
		position: fixed;
		bottom: 1.25rem;
		right: 1.25rem;
		max-width: min(420px, calc(100vw - 2rem));
		padding: 0.75rem 1rem;
		border-radius: var(--radius);
		box-shadow: var(--shadow);
		z-index: 1000;
		background: var(--surface2);
		border: 1px solid var(--border);
	}
	.toast.error {
		border-color: var(--danger);
	}
	.toast.info {
		border-color: var(--accent);
	}
</style>
