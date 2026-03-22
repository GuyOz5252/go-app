import { openDB, type DBSchema, type IDBPDatabase } from 'idb';
import type { ChatMessage } from '$lib/api/types';

const DB_NAME = 'go-chat-messages';
const DB_VERSION = 1;

interface GoChatSchema extends DBSchema {
	messages: {
		key: string;
		value: ChatMessage;
		indexes: { 'by-chat': string };
	};
}

let dbPromise: Promise<IDBPDatabase<GoChatSchema>> | null = null;

function getDb(): Promise<IDBPDatabase<GoChatSchema>> {
	if (!dbPromise) {
		dbPromise = openDB<GoChatSchema>(DB_NAME, DB_VERSION, {
			upgrade(database) {
				if (!database.objectStoreNames.contains('messages')) {
					const store = database.createObjectStore('messages', { keyPath: 'id' });
					store.createIndex('by-chat', 'chat_id');
				}
			}
		});
	}
	return dbPromise;
}

export async function loadMessagesForChat(chatId: string): Promise<ChatMessage[]> {
	const db = await getDb();
	const tx = db.transaction('messages', 'readonly');
	const idx = tx.store.index('by-chat');
	const rows = await idx.getAll(chatId);
	await tx.done;
	return rows.sort(
		(a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
	);
}

export async function saveMessage(message: ChatMessage): Promise<void> {
	const db = await getDb();
	await db.put('messages', message);
}

export async function saveMessages(messages: ChatMessage[]): Promise<void> {
	const db = await getDb();
	const tx = db.transaction('messages', 'readwrite');
	for (const m of messages) {
		await tx.store.put(m);
	}
	await tx.done;
}
