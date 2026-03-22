<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import { page } from '$app/stores';
	import { resolve } from '$app/paths';
	import { auth } from '$lib/auth/auth.svelte';
	import { api, ApiError } from '$lib/api/client';
	import { toast } from '$lib/toast.svelte';
	import { CHAT_WS } from '$lib/context';
	import type { ChatWebSocket } from '$lib/ws/connection';
	import { MessageType, isChatMessagePayload, upsertMessages, type WSMessage } from '$lib/ws/protocol';
	import type { ChatMessage } from '$lib/api/types';
	import * as messagesDb from '$lib/db/messagesDb';
	import { linkifyToHtml, isProbablyImageUrl } from '$lib/utils/text';
	import { ArrowLeft, Reply, Send } from 'lucide-svelte';

	type LocalMessage = ChatMessage & { pending?: boolean };

	const socket = getContext<ChatWebSocket>(CHAT_WS);

	let messages = $state<LocalMessage[]>([]);
	let draft = $state('');
	let mediaUrl = $state('');
	let replyTo = $state<ChatMessage | null>(null);
	let loading = $state(true);
	let typingUsers = $state<Set<string>>(new Set());
	let typingTimer: ReturnType<typeof setTimeout> | null = null;
	let localTypingActive = $state(false);
	const nameCache = new Map<string, string>();
	let displayNames = $state<Record<string, string>>({});
	const readSent = new Set<string>();
	let unsub: (() => void) | null = null;
	let listEl = $state<HTMLDivElement | null>(null);

	let chatId = $derived($page.params.chatId ?? '');

	$effect(() => {
		const id = chatId;
		if (!id) return;
		void bootstrap(id);
		return () => {
			if (unsub) {
				unsub();
				unsub = null;
			}
		};
	});

	$effect(() => {
		const ids = [...new Set(messages.map((m) => m.user_id))];
		void (async () => {
			for (const uid of ids) {
				if (displayNames[uid]) continue;
				const n = await resolveName(uid);
				displayNames = { ...displayNames, [uid]: n };
			}
		})();
	});

	onMount(() => {
		return () => {
			if (typingTimer) clearTimeout(typingTimer);
		};
	});

	async function bootstrap(id: string) {
		if (!id) return;
		loading = true;
		messages = [];
		readSent.clear();
		try {
			const stored = await messagesDb.loadMessagesForChat(id);
			messages = stored.map((m) => ({ ...m }));
		} catch {
			toast.show('Could not load chat history from this device.', 'error');
		} finally {
			loading = false;
		}
		if (unsub) unsub();
		unsub = socket.subscribeMessages((msg) => handleWs(msg, id));
	}

	async function resolveName(userId: string): Promise<string> {
		if (nameCache.has(userId)) return nameCache.get(userId)!;
		if (userId === auth.user?.id) {
			nameCache.set(userId, auth.user.name);
			return auth.user.name;
		}
		const t = auth.token;
		if (!t) return userId;
		try {
			const u = await api.getUser(t, userId);
			nameCache.set(userId, u.name);
			return u.name;
		} catch {
			return `User ${userId}`;
		}
	}

	function handleWs(msg: WSMessage, currentChatId: string) {
		if (msg.message_type === MessageType.NewMessage && isChatMessagePayload(msg.payload)) {
			const m = msg.payload;
			if (m.chat_id !== currentChatId) return;
			messages = upsertMessages(messages, [m]) as LocalMessage[];
			void messagesDb.saveMessage(m);
			return;
		}
		if (msg.message_type === MessageType.MessageServerAck) {
			const p = msg.payload as { message_id?: string } | undefined;
			const mid = p?.message_id;
			if (!mid) return;
			const list = [...messages];
			let idx = -1;
			for (let i = list.length - 1; i >= 0; i--) {
				const row = list[i];
				if (row.pending && row.user_id === auth.user?.id && row.chat_id === currentChatId) {
					idx = i;
					break;
				}
			}
			if (idx === -1) return;
			const updated = { ...list[idx], id: mid, pending: false };
			list[idx] = updated;
			messages = list.sort(
				(a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
			);
			void messagesDb.saveMessage(updated);
			return;
		}
		if (msg.message_type === MessageType.UserTypingStart) {
			if (msg.destination_chat_id !== currentChatId) return;
			const uid = msg.initiator_user_id;
			if (!uid || uid === auth.user?.id) return;
			const s = new Set(typingUsers);
			s.add(uid);
			typingUsers = s;
			return;
		}
		if (msg.message_type === MessageType.UserTypingEnd) {
			if (msg.destination_chat_id !== currentChatId) return;
			const uid = msg.initiator_user_id;
			if (!uid) return;
			const s = new Set(typingUsers);
			s.delete(uid);
			typingUsers = s;
		}
	}

	function scheduleTypingEnd() {
		if (typingTimer) clearTimeout(typingTimer);
		typingTimer = setTimeout(() => {
			const id = chatId;
			const u = auth.user?.id;
			if (u && id && localTypingActive) {
				socket.sendTypingEnd(u, id);
				localTypingActive = false;
			}
		}, 320);
	}

	function onComposerInput() {
		const id = chatId;
		const u = auth.user?.id;
		if (!u || !id) return;
		if (!localTypingActive) {
			socket.sendTypingStart(u, id);
			localTypingActive = true;
		}
		scheduleTypingEnd();
	}

	async function send() {
		const content = draft.trim();
		const media = mediaUrl.trim();
		const replyToId = replyTo?.id;
		if (!content && !media) {
			toast.show('Enter a message or a media URL.', 'error');
			return;
		}
		const id = chatId;
		const u = auth.user?.id;
		const t = auth.token;
		if (!t || !u || !id) return;

		const tempId = `temp-${crypto.randomUUID()}`;
		const optimistic: LocalMessage = {
			id: tempId,
			user_id: u,
			chat_id: id,
			content: content || '(attachment)',
			media_url: media || undefined,
			reply_to_id: replyToId,
			created_at: new Date().toISOString(),
			pending: true
		};
		messages = [...messages, optimistic];

		draft = '';
		mediaUrl = '';
		replyTo = null;
		if (localTypingActive) {
			socket.sendTypingEnd(u, id);
			localTypingActive = false;
		}
		if (typingTimer) clearTimeout(typingTimer);

		try {
			if (socket.isOpen) {
				socket.sendNewMessage(u, id, {
					content: content || ' ',
					media_url: media || undefined,
					reply_to_id: replyToId
				});
			} else {
				const saved = await api.sendMessage(t, id, {
					content: content || ' ',
					media_url: media || undefined,
					reply_to_id: replyToId
				});
				messages = upsertMessages(
					messages.filter((m) => m.id !== tempId),
					[saved]
				) as LocalMessage[];
				await messagesDb.saveMessage(saved);
			}
		} catch (e) {
			messages = messages.filter((m) => m.id !== tempId);
			const msg =
				e instanceof ApiError ? e.detail || e.title : e instanceof Error ? e.message : 'Send failed';
			toast.show(msg, 'error');
		}
	}

	function setReply(m: ChatMessage) {
		replyTo = m;
	}

	function observeRead(node: HTMLElement, message: ChatMessage) {
		const selfId = auth.user?.id;
		if (!selfId || message.user_id === selfId) return { destroy() {} };

		const io = new IntersectionObserver(
			(entries) => {
				for (const e of entries) {
					if (!e.isIntersecting) continue;
					if (readSent.has(message.id)) continue;
					readSent.add(message.id);
					if (socket.isOpen) {
						socket.sendMessageReadAck(selfId, message.user_id, message.id);
					}
					io.disconnect();
				}
			},
			{ root: listEl, threshold: 0.35 }
		);
		io.observe(node);
		return {
			destroy() {
				io.disconnect();
			}
		};
	}
</script>

<div class="room">
	<header class="bar">
		<a href={resolve('/')} class="back"><ArrowLeft size={20} /> Chats</a>
		<div class="title">
			<h1>Chat</h1>
			<span class="cid">#{chatId}</span>
		</div>
	</header>

	{#if loading}
		<p class="muted pad">Loading…</p>
	{:else}
		<div class="list" bind:this={listEl}>
			{#each messages as m (m.id)}
				<article
					class="bubble"
					class:own={m.user_id === auth.user?.id}
					class:pending={m.pending}
					use:observeRead={m}
				>
					<div class="row">
						<span class="who">{displayNames[m.user_id] ?? m.user_id}</span>
						<time datetime={m.created_at}>{new Date(m.created_at).toLocaleString()}</time>
					</div>
					{#if m.reply_to_id}
						<p class="replyref">Replying to message {m.reply_to_id}</p>
					{/if}
					<div class="body text">
						<!-- eslint-disable-next-line svelte/no-at-html-tags -->
						{@html linkifyToHtml(m.content)}
					</div>
					{#if m.media_url}
						{#if isProbablyImageUrl(m.media_url)}
							<div class="media">
								<img src={m.media_url} alt="" referrerpolicy="no-referrer" />
							</div>
						{:else}
							<p class="media-link">
								<!-- eslint-disable-next-line svelte/no-navigation-without-resolve -- external URL -->
								<a href={m.media_url} target="_blank" rel="noopener noreferrer">{m.media_url}</a>
							</p>
						{/if}
					{/if}
					<div class="actions">
						<button type="button" class="linkish" onclick={() => setReply(m)}>
							<Reply size={14} /> Reply
						</button>
					</div>
				</article>
			{/each}
		</div>

		{#if typingUsers.size > 0}
			<p class="typing" aria-live="polite">
				{#each [...typingUsers] as uid (uid)}
					<span>{displayNames[uid] ?? uid} is typing…</span>
				{/each}
			</p>
		{/if}

		<div class="composer">
			{#if replyTo}
				<div class="reply-banner">
					Reply to {replyTo.id.slice(0, 8)}…
					<button type="button" onclick={() => (replyTo = null)}>Cancel</button>
				</div>
			{/if}
			<textarea
				rows="3"
				placeholder="Write a message…"
				bind:value={draft}
				oninput={onComposerInput}
				aria-label="Message text"
			></textarea>
			<input type="url" placeholder="Media URL (optional)" bind:value={mediaUrl} aria-label="Media URL" />
			<button type="button" class="send" onclick={() => send()}>
				<Send size={18} /> Send
			</button>
		</div>
	{/if}
</div>

<style>
	.room {
		display: flex;
		flex-direction: column;
		height: 100%;
		min-height: 0;
	}
	.bar {
		display: flex;
		align-items: center;
		gap: 1rem;
		padding: 0.75rem 1rem;
		border-bottom: 1px solid var(--border);
		background: var(--surface);
	}
	.back {
		display: inline-flex;
		align-items: center;
		gap: 0.35rem;
		color: var(--muted);
		font-size: 0.9rem;
	}
	.title h1 {
		margin: 0;
		font-size: 1.1rem;
	}
	.cid {
		font-size: 0.8rem;
		color: var(--muted);
	}
	.pad {
		padding: 1rem;
	}
	.muted {
		color: var(--muted);
	}
	.list {
		flex: 1;
		overflow: auto;
		padding: 1rem;
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}
	.bubble {
		align-self: flex-start;
		max-width: min(560px, 100%);
		padding: 0.65rem 0.85rem;
		border-radius: 12px;
		background: var(--surface);
		border: 1px solid var(--border);
	}
	.bubble.own {
		align-self: flex-end;
		background: var(--surface2);
	}
	.bubble.pending {
		opacity: 0.75;
	}
	.row {
		display: flex;
		justify-content: space-between;
		gap: 0.75rem;
		font-size: 0.8rem;
		color: var(--muted);
		margin-bottom: 0.35rem;
	}
	.who {
		font-weight: 600;
		color: var(--text);
	}
	.replyref {
		font-size: 0.8rem;
		color: var(--muted);
		margin: 0 0 0.25rem;
	}
	.body {
		word-break: break-word;
	}
	.media {
		margin-top: 0.5rem;
		border-radius: 8px;
		overflow: hidden;
		max-width: 280px;
	}
	.media img {
		display: block;
		width: 100%;
		height: auto;
	}
	.media-link {
		margin: 0.35rem 0 0;
		font-size: 0.9rem;
	}
	.actions {
		margin-top: 0.35rem;
	}
	.linkish {
		background: none;
		border: none;
		color: var(--accent);
		cursor: pointer;
		display: inline-flex;
		align-items: center;
		gap: 0.2rem;
		font-size: 0.8rem;
		padding: 0;
	}
	.typing {
		font-size: 0.85rem;
		color: var(--muted);
		padding: 0 1rem;
		margin: 0;
	}
	.composer {
		border-top: 1px solid var(--border);
		padding: 0.75rem 1rem;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		background: var(--surface);
	}
	.composer textarea,
	.composer input {
		width: 100%;
		padding: 0.75rem;
		border-radius: 8px;
		border: 1px solid var(--border);
		background: var(--bg);
		color: var(--text);
	}
	.send {
		align-self: flex-end;
		display: inline-flex;
		align-items: center;
		gap: 0.4rem;
		padding: 0.65rem 1rem;
		border-radius: 8px;
		border: 1px solid var(--accent);
		background: var(--accent);
		color: white;
		cursor: pointer;
		font-weight: 600;
	}
	.reply-banner {
		display: flex;
		justify-content: space-between;
		align-items: center;
		font-size: 0.85rem;
		color: var(--muted);
	}
	.reply-banner button {
		background: none;
		border: none;
		color: var(--accent);
		cursor: pointer;
	}
</style>
