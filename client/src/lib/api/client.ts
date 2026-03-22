import type {
	ChatDto,
	ChatMessage,
	CreateChatBody,
	CreateChatResponse,
	LoginBody,
	LoginResponse,
	ProblemDetails,
	RegisterBody,
	RegisterResponse,
	SendMessageBody,
	User
} from './types';

export class ApiError extends Error {
	readonly status: number;
	readonly title: string;
	readonly detail?: string;
	readonly instance?: string;

	constructor(status: number, title: string, detail?: string, instance?: string) {
		super(title);
		this.name = 'ApiError';
		this.status = status;
		this.title = title;
		this.detail = detail;
		this.instance = instance;
	}
}

function baseUrl(): string {
	const b = import.meta.env.PUBLIC_API_BASE;
	return typeof b === 'string' ? b.replace(/\/$/, '') : '';
}

export function apiUrl(path: string): string {
	const p = path.startsWith('/') ? path : `/${path}`;
	return `${baseUrl()}${p}`;
}

async function parseJson<T>(res: Response): Promise<T> {
	const ct = res.headers.get('content-type') || '';
	if (!res.ok) {
		if (ct.includes('application/problem+json')) {
			const p = (await res.json()) as ProblemDetails;
			throw new ApiError(p.status, p.title, p.detail, p.instance);
		}
		let detail: string | undefined;
		try {
			detail = await res.text();
		} catch {
			detail = undefined;
		}
		throw new ApiError(res.status, res.statusText, detail);
	}
	if (res.status === 204) return undefined as T;
	if (ct.includes('application/json')) {
		return (await res.json()) as T;
	}
	return (await res.text()) as T;
}

type RequestOpts = Omit<RequestInit, 'body'> & { body?: unknown };

async function request<T>(path: string, token: string | null, opts: RequestOpts = {}): Promise<T> {
	const headers = new Headers(opts.headers);
	if (opts.body !== undefined && !(opts.body instanceof FormData)) {
		headers.set('Content-Type', 'application/json');
	}
	if (token) {
		headers.set('Authorization', `Bearer ${token}`);
	}
	const { body, ...rest } = opts;
	const init: RequestInit = {
		...rest,
		headers
	};
	if (body !== undefined && !(body instanceof FormData)) {
		init.body = JSON.stringify(body);
	} else if (body instanceof FormData) {
		init.body = body;
	}
	const res = await fetch(apiUrl(path), init);
	return parseJson<T>(res);
}

export const api = {
	register: (body: RegisterBody) =>
		request<RegisterResponse>('/api/users/', null, { method: 'POST', body }),

	login: (body: LoginBody) =>
		request<LoginResponse>('/api/users/login', null, { method: 'POST', body }),

	getUser: (token: string, id: string) =>
		request<User>(`/api/users/${encodeURIComponent(id)}`, token, { method: 'GET' }),

	listChats: (token: string) => request<ChatDto[]>('/api/chats/', token, { method: 'GET' }),

	createChat: (token: string, body: CreateChatBody) =>
		request<CreateChatResponse>('/api/chats/', token, { method: 'POST', body }),

	sendMessage: (token: string, chatId: string, body: SendMessageBody) =>
		request<ChatMessage>(`/api/chats/${encodeURIComponent(chatId)}/messages`, token, {
			method: 'POST',
			body
		})
};
