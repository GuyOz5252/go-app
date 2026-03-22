import type { ChatWebSocket } from '$lib/ws/connection';

export const CHAT_WS = Symbol('chat-ws');

export type ChatWsContext = ChatWebSocket;
