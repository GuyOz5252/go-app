<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { auth } from '$lib/auth/auth.svelte';
	import { api } from '$lib/api/client';
	import { toast } from '$lib/toast.svelte';
	import { ApiError } from '$lib/api/client';
	import { Users } from 'lucide-svelte';

	let name = $state('');
	let imageUrl = $state('');
	let memberIdsRaw = $state('');
	let loading = $state(false);

	function parseMemberIds(raw: string): string[] {
		return raw
			.split(/[\s,]+/)
			.map((s) => s.trim())
			.filter(Boolean);
	}

	async function onSubmit(e: Event) {
		e.preventDefault();
		const selfId = auth.user?.id;
		if (!selfId) {
			toast.show('Session invalid.', 'error');
			return;
		}
		const ids = parseMemberIds(memberIdsRaw);
		if (!ids.includes(selfId)) {
			ids.push(selfId);
		}
		if (ids.length <= 1) {
			toast.show('Add at least one other member (their user id). Your id is included automatically.', 'error');
			return;
		}
		const t = auth.token;
		if (!t) return;
		loading = true;
		try {
			const { id } = await api.createChat(t, {
				name: name.trim(),
				chat_member_ids: ids,
				image_url: imageUrl.trim() || undefined
			});
			await goto(resolve(`/chats/${id}`));
		} catch (err) {
			const msg =
				err instanceof ApiError ? err.detail || err.title : err instanceof Error ? err.message : 'Could not create chat';
			toast.show(msg, 'error');
		} finally {
			loading = false;
		}
	}
</script>

<div class="page">
	<h1>New chat</h1>
	<p class="hint">
		The API requires <strong>at least two members</strong> in the chat. Include yourself (your id:
		<strong>{auth.user?.id ?? '—'}</strong>) and one or more other users by their numeric ids.
	</p>

	<form class="form" onsubmit={onSubmit}>
		<label for="name">Chat name</label>
		<input id="name" required bind:value={name} placeholder="Team standup" />

		<label for="image">Image URL <span class="opt">(optional)</span></label>
		<input id="image" type="url" bind:value={imageUrl} placeholder="https://…" />

		<label for="members">Member user ids</label>
		<textarea
			id="members"
			rows="3"
			required
			bind:value={memberIdsRaw}
			placeholder="e.g. 2 3 or 2,3"
		></textarea>
		<p class="help">
			<Users size={16} class="inline" aria-hidden="true" />
			Separate ids with commas or spaces. You will be added if missing.
		</p>

		<div class="actions">
			<a href={resolve('/')} class="btn ghost">Cancel</a>
			<button type="submit" class="btn primary" disabled={loading}>{loading ? 'Creating…' : 'Create chat'}</button>
		</div>
	</form>
</div>

<style>
	.page {
		padding: 1.25rem;
		max-width: 560px;
		margin: 0 auto;
		width: 100%;
	}
	h1 {
		margin: 0 0 0.5rem;
		font-size: 1.35rem;
	}
	.hint {
		color: var(--muted);
		font-size: 0.95rem;
		margin: 0 0 1.25rem;
		line-height: 1.5;
	}
	.form {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}
	label {
		font-weight: 600;
		font-size: 0.9rem;
	}
	.opt {
		font-weight: 400;
		color: var(--muted);
	}
	input,
	textarea {
		padding: 0.6rem 0.75rem;
		border-radius: 8px;
		border: 1px solid var(--border);
		background: var(--bg);
		color: var(--text);
	}
	.help {
		display: flex;
		align-items: flex-start;
		gap: 0.35rem;
		font-size: 0.85rem;
		color: var(--muted);
		margin: 0;
	}
	:global(.inline) {
		flex-shrink: 0;
		margin-top: 2px;
	}
	.actions {
		display: flex;
		gap: 0.75rem;
		margin-top: 0.5rem;
	}
	.btn {
		padding: 0.55rem 1rem;
		border-radius: 8px;
		border: 1px solid var(--border);
		cursor: pointer;
		font-size: 0.95rem;
		text-decoration: none;
		display: inline-flex;
		align-items: center;
		justify-content: center;
	}
	.btn.primary {
		background: var(--accent);
		border-color: var(--accent);
		color: white;
	}
	.btn.ghost {
		background: transparent;
		color: var(--text);
	}
	.btn:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}
</style>
