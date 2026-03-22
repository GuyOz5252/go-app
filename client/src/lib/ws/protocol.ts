import type { ChatMessage } from '$lib/api/types';

/** Matches internal/websocket/types.go WSMessageType iota order. */
export enum MessageType {
	NewMessage = 0,
	MessageServerAck = 1,
	MessageUserAck = 2,
	MessageUserReadAck = 3,
	UserTypingStart = 4,
	UserTypingEnd = 5,
	UserOnline = 6,
	UserOffline = 7,
	UserAway = 8,
	ServerError = 9
}

export type WSMessage = {
	message_type: MessageType;
	initiator_user_id?: string;
	destination_chat_id?: string;
	destination_user_id?: string;
	payload?: unknown;
};

export type NewMessagePayload = {
	content: string;
	media_url?: string;
	reply_to_id?: string;
};

export type MessageIdPayload = {
	message_id: string;
};

export type ServerErrorPayload = {
	error: string;
};

export function isChatMessagePayload(p: unknown): p is ChatMessage {
	if (!p || typeof p !== 'object') return false;
	const o = p as Record<string, unknown>;
	return (
		typeof o.id === 'string' &&
		typeof o.user_id === 'string' &&
		typeof o.chat_id === 'string' &&
		typeof o.content === 'string' &&
		typeof o.created_at === 'string'
	);
}

export function mergeMessagesById(existing: ChatMessage[], incoming: ChatMessage): ChatMessage[] {
	const map = new Map<string, ChatMessage>();
	for (const m of existing) {
		map.set(m.id, m);
	}
	map.set(incoming.id, incoming);
	return [...map.values()].sort(
		(a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
	);
}

export function upsertMessages(existing: ChatMessage[], incoming: ChatMessage[]): ChatMessage[] {
	const map = new Map<string, ChatMessage>();
	for (const m of existing) map.set(m.id, m);
	for (const m of incoming) map.set(m.id, m);
	return [...map.values()].sort(
		(a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
	);
}
