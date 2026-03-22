import {
	MessageType,
	type MessageIdPayload,
	type NewMessagePayload,
	type ServerErrorPayload,
	type WSMessage
} from './protocol';

export type WsConnectionStatus = 'idle' | 'connecting' | 'open' | 'reconnecting' | 'closed';

function buildWebSocketUrl(accessToken: string): string {
	const prefix = (import.meta.env.PUBLIC_API_BASE || '').replace(/\/$/, '');
	const qs = new URLSearchParams({ access_token: accessToken });
	const path = `${prefix}/api/chats/ws?${qs.toString()}`;
	const wsProto = typeof window !== 'undefined' && window.location.protocol === 'https:' ? 'wss:' : 'ws:';
	const host = typeof window !== 'undefined' ? window.location.host : '';
	return `${wsProto}//${host}${path}`;
}

export type WsHandlers = {
	onStatus: (status: WsConnectionStatus) => void;
	onServerError: (message: string) => void;
};

/**
 * Manages a single WebSocket to /api/chats/ws with reconnect backoff.
 * Browser cannot send Authorization; URL carries access_token for the dev proxy to rewrite.
 */
export class ChatWebSocket {
	private ws: WebSocket | null = null;
	private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
	private attempt = 0;
	private readonly maxDelayMs = 30_000;
	private token: string | null = null;
	private shouldReconnect = false;
	private readonly handlers: WsHandlers;
	private readonly messageSubscribers = new Set<(msg: WSMessage) => void>();

	constructor(handlers: WsHandlers) {
		this.handlers = handlers;
	}

	/** Receive every inbound frame after JSON parse (including ServerError, before global handler). */
	subscribeMessages(handler: (msg: WSMessage) => void) {
		this.messageSubscribers.add(handler);
		return () => this.messageSubscribers.delete(handler);
	}

	connect(token: string) {
		this.token = token;
		this.shouldReconnect = true;
		this.attempt = 0;
		this.clearReconnectTimer();
		this.closeSocket();
		this.open();
	}

	disconnect() {
		this.shouldReconnect = false;
		this.clearReconnectTimer();
		this.closeSocket();
		this.handlers.onStatus('closed');
	}

	private clearReconnectTimer() {
		if (this.reconnectTimer !== null) {
			clearTimeout(this.reconnectTimer);
			this.reconnectTimer = null;
		}
	}

	private closeSocket() {
		if (this.ws) {
			this.ws.onopen = null;
			this.ws.onclose = null;
			this.ws.onerror = null;
			this.ws.onmessage = null;
			try {
				this.ws.close();
			} catch {
				/* ignore */
			}
			this.ws = null;
		}
	}

	private open() {
		if (!this.token) return;
		this.handlers.onStatus(this.attempt > 0 ? 'reconnecting' : 'connecting');
		let url: string;
		try {
			url = buildWebSocketUrl(this.token);
		} catch (e) {
			this.handlers.onServerError(e instanceof Error ? e.message : 'Invalid WebSocket URL');
			return;
		}
		const ws = new WebSocket(url);
		this.ws = ws;

		ws.onopen = () => {
			this.attempt = 0;
			this.handlers.onStatus('open');
		};

		ws.onmessage = (ev) => {
			try {
				const raw = JSON.parse(String(ev.data)) as WSMessage;
				for (const h of this.messageSubscribers) {
					try {
						h(raw);
					} catch {
						/* isolate subscriber errors */
					}
				}
				if (raw.message_type === MessageType.ServerError) {
					const p = raw.payload as ServerErrorPayload | undefined;
					this.handlers.onServerError(p?.error ?? 'Unknown server error');
					return;
				}
			} catch (e) {
				this.handlers.onServerError(e instanceof Error ? e.message : 'Invalid WebSocket frame');
			}
		};

		ws.onerror = () => {
			/* onclose will handle reconnect */
		};

		ws.onclose = () => {
			this.ws = null;
			if (!this.shouldReconnect || !this.token) {
				this.handlers.onStatus('closed');
				return;
			}
			const delay = Math.min(this.maxDelayMs, 1000 * 2 ** this.attempt);
			this.attempt += 1;
			this.handlers.onStatus('reconnecting');
			this.reconnectTimer = setTimeout(() => this.open(), delay);
		};
	}

	sendRaw(msg: WSMessage) {
		if (this.ws?.readyState === WebSocket.OPEN) {
			this.ws.send(JSON.stringify(msg));
		}
	}

	sendNewMessage(initiatorUserId: string, chatId: string, payload: NewMessagePayload) {
		this.sendRaw({
			message_type: MessageType.NewMessage,
			initiator_user_id: initiatorUserId,
			destination_chat_id: chatId,
			payload
		});
	}

	sendTypingStart(initiatorUserId: string, chatId: string) {
		this.sendRaw({
			message_type: MessageType.UserTypingStart,
			initiator_user_id: initiatorUserId,
			destination_chat_id: chatId
		});
	}

	sendTypingEnd(initiatorUserId: string, chatId: string) {
		this.sendRaw({
			message_type: MessageType.UserTypingEnd,
			initiator_user_id: initiatorUserId,
			destination_chat_id: chatId
		});
	}

	sendMessageReadAck(initiatorUserId: string, destinationUserId: string, messageId: string) {
		this.sendRaw({
			message_type: MessageType.MessageUserReadAck,
			initiator_user_id: initiatorUserId,
			destination_user_id: destinationUserId,
			payload: { message_id: messageId } satisfies MessageIdPayload
		});
	}

	get isOpen(): boolean {
		return this.ws?.readyState === WebSocket.OPEN;
	}
}
