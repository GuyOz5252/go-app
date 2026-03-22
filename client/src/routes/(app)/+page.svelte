<script lang="ts">
	import { onMount } from 'svelte';
	import { resolve } from '$app/paths';
	import { auth } from '$lib/auth/auth.svelte';
	import { api, ApiError } from '$lib/api/client';
	import { toast } from '$lib/toast.svelte';
	import type { ChatDto } from '$lib/api/types';
	import { ChevronRight } from 'lucide-svelte';

	let chats = $state<ChatDto[]>([]);
	let loading = $state(true);

	onMount(() => {
		void load();
	});

	async function load() {
		const t = auth.token;
		if (!t) return;
		loading = true;
		try {
			chats = await api.listChats(t);
		} catch (e) {
			const msg =
				e instanceof ApiError ? e.detail || e.title : e instanceof Error ? e.message : 'Could not load chats';
			toast.show(msg, 'error');
		} finally {
			loading = false;
		}
	}
</script>

<div class="page">
	<div class="head">
		<h1>Your chats</h1>
		<p class="sub">
			Signed in as <strong>{auth.user?.name}</strong>
			{#if auth.user}
				<span class="id">· user id {auth.user.id}</span>
			{/if}
		</p>
	</div>

	{#if loading}
		<p class="muted">Loading conversations…</p>
	{:else if chats.length === 0}
		<div class="empty">
			<p>No chats yet. Create one with at least one other member (you need their user id).</p>
			<a href={resolve('/chats/new')} class="btn primary">Start a chat</a>
		</div>
	{:else}
		<ul class="list" role="list">
			{#each chats as c (c.id)}
				<li>
					<a href={resolve(`/chats/${c.id}`)} class="row">
						<div class="avatar" aria-hidden="true">
							{#if c.image_url}
								<img src={c.image_url} alt="" />
							{:else}
								<span>{c.name.slice(0, 1).toUpperCase()}</span>
							{/if}
						</div>
						<div class="meta">
							<span class="name">{c.name}</span>
							<span class="id">#{c.id}</span>
						</div>
						<ChevronRight size={18} class="chev" aria-hidden="true" />
					</a>
				</li>
			{/each}
		</ul>
	{/if}
</div>

<style>
	.page {
		padding: 1.25rem;
		max-width: 720px;
		margin: 0 auto;
		width: 100%;
	}
	.head h1 {
		margin: 0 0 0.25rem;
		font-size: 1.35rem;
	}
	.sub {
		margin: 0;
		color: var(--muted);
		font-size: 0.95rem;
	}
	.id {
		font-size: 0.85rem;
		opacity: 0.85;
	}
	.muted {
		color: var(--muted);
	}
	.empty {
		margin-top: 2rem;
		padding: 1.5rem;
		border: 1px dashed var(--border);
		border-radius: var(--radius);
		text-align: center;
	}
	.btn {
		display: inline-block;
		margin-top: 1rem;
		padding: 0.5rem 1rem;
		border-radius: 8px;
		text-decoration: none;
	}
	.btn.primary {
		background: var(--accent);
		color: white;
		border: 1px solid var(--accent);
	}
	.list {
		list-style: none;
		padding: 0;
		margin: 1rem 0 0;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}
	.row {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.75rem 1rem;
		border-radius: var(--radius);
		border: 1px solid var(--border);
		background: var(--surface);
		color: inherit;
		text-decoration: none;
		transition: background 0.15s;
	}
	.row:hover {
		background: var(--surface2);
	}
	.avatar {
		width: 40px;
		height: 40px;
		border-radius: 10px;
		background: var(--surface2);
		display: grid;
		place-items: center;
		overflow: hidden;
		font-weight: 700;
	}
	.avatar img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}
	.meta {
		flex: 1;
		min-width: 0;
	}
	.name {
		display: block;
		font-weight: 600;
	}
	.meta .id {
		font-size: 0.8rem;
		color: var(--muted);
	}
	:global(.chev) {
		color: var(--muted);
	}
</style>
