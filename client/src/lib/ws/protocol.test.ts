import { describe, expect, it } from 'vitest';
import type { ChatMessage } from '$lib/api/types';
import { mergeMessagesById, upsertMessages } from './protocol';

const m = (id: string, chatId: string, t: string): ChatMessage => ({
	id,
	user_id: '1',
	chat_id: chatId,
	content: 'x',
	created_at: t
});

describe('mergeMessagesById', () => {
	it('merges and sorts by created_at', () => {
		const a = mergeMessagesById(
			[m('1', 'c', '2025-01-01T00:00:00Z')],
			m('2', 'c', '2025-01-02T00:00:00Z')
		);
		expect(a.map((x) => x.id)).toEqual(['1', '2']);
	});

	it('replaces same id', () => {
		const a = mergeMessagesById([m('1', 'c', '2025-01-01T00:00:00Z')], {
			...m('1', 'c', '2025-01-01T00:00:00Z'),
			content: 'updated'
		});
		expect(a).toHaveLength(1);
		expect(a[0].content).toBe('updated');
	});
});

describe('upsertMessages', () => {
	it('merges arrays', () => {
		const a = upsertMessages([m('1', 'c', '2025-01-01T00:00:00Z')], [m('2', 'c', '2025-01-02T00:00:00Z')]);
		expect(a.map((x) => x.id)).toEqual(['1', '2']);
	});
});
