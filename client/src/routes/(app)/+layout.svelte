<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { setContext } from 'svelte';
	import { auth } from '$lib/auth/auth.svelte';
	import { ChatWebSocket } from '$lib/ws/connection';
	import { CHAT_WS } from '$lib/context';
	import { toast } from '$lib/toast.svelte';
	import { LogOut, MessageSquarePlus, MessagesSquare } from 'lucide-svelte';

	let { children } = $props();

	let ready = $state(false);
	let wsStatus = $state<'idle' | 'connecting' | 'open' | 'reconnecting' | 'closed'>('idle');

	const socket = new ChatWebSocket({
		onStatus: (s) => {
			wsStatus = s;
		},
		onServerError: (e) => toast.show(e, 'error')
	});
	setContext(CHAT_WS, socket);

	onMount(() => {
		void (async () => {
			await auth.hydrate();
			if (!auth.token) {
				await goto(resolve('/login'));
				return;
			}
			ready = true;
			socket.connect(auth.token);
		})();

		return () => socket.disconnect();
	});

	async function logout() {
		socket.disconnect();
		auth.logout();
		await goto(resolve('/login'));
	}

	function statusLabel(): string {
		switch (wsStatus) {
			case 'open':
				return 'Live connection';
			case 'connecting':
				return 'Connecting…';
			case 'reconnecting':
				return 'Reconnecting…';
			case 'closed':
				return 'Disconnected';
			default:
				return '';
		}
	}
</script>

{#if !ready}
	<div class="loading">Loading your session…</div>
{:else}
	<div class="shell">
		<header class="top">
			<div class="brand">
				<MessagesSquare size={22} aria-hidden="true" />
				<span class="title">Go Chat</span>
			</div>
			<div class="status" role="status">
				<span class="dot" class:ok={wsStatus === 'open'} class:warn={wsStatus !== 'open' && wsStatus !== 'idle'}></span>
				{statusLabel()}
			</div>
			<nav class="nav">
				<a href={resolve('/chats/new')} class="nav-btn"><MessageSquarePlus size={18} /> New chat</a>
				<button type="button" class="nav-btn ghost" onclick={logout}>
					<LogOut size={18} /> Sign out
				</button>
			</nav>
		</header>
		<main class="main">
			{@render children()}
		</main>
	</div>
{/if}

<style>
	.loading {
		min-height: 100vh;
		display: grid;
		place-items: center;
		color: var(--muted);
	}
	.shell {
		min-height: 100vh;
		display: flex;
		flex-direction: column;
	}
	.top {
		display: flex;
		align-items: center;
		gap: 1rem;
		padding: 0.75rem 1.25rem;
		border-bottom: 1px solid var(--border);
		background: var(--surface);
		flex-wrap: wrap;
	}
	.brand {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-weight: 700;
	}
	.title {
		font-size: 1.1rem;
	}
	.status {
		flex: 1;
		display: flex;
		align-items: center;
		gap: 0.4rem;
		font-size: 0.85rem;
		color: var(--muted);
	}
	.dot {
		width: 8px;
		height: 8px;
		border-radius: 50%;
		background: var(--muted);
	}
	.dot.ok {
		background: var(--success);
	}
	.dot.warn {
		background: #fbbf24;
	}
	.nav {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}
	.nav-btn {
		display: inline-flex;
		align-items: center;
		gap: 0.35rem;
		padding: 0.45rem 0.75rem;
		border-radius: 8px;
		border: 1px solid var(--border);
		background: var(--surface2);
		color: var(--text);
		font-size: 0.9rem;
		cursor: pointer;
		text-decoration: none;
	}
	.nav-btn:hover {
		background: var(--border);
	}
	.nav-btn.ghost {
		background: transparent;
	}
	.main {
		flex: 1;
		display: flex;
		flex-direction: column;
		min-height: 0;
	}
</style>
