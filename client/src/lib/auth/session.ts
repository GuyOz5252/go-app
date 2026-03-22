const TOKEN_KEY = 'go-chat-token';

export function saveToken(t: string) {
	sessionStorage.setItem(TOKEN_KEY, t);
}

export function loadToken(): string | null {
	if (typeof sessionStorage === 'undefined') return null;
	return sessionStorage.getItem(TOKEN_KEY);
}

export function clearToken() {
	sessionStorage.removeItem(TOKEN_KEY);
}

/** Read JWT claims (client convenience only; API calls still send the token to the server). */
export function parseJwtPayload(token: string): { userId?: string; exp?: number } | null {
	try {
		const parts = token.split('.');
		if (parts.length !== 3) return null;
		const json = atob(parts[1].replace(/-/g, '+').replace(/_/g, '/'));
		return JSON.parse(json) as { userId?: string; exp?: number };
	} catch {
		return null;
	}
}
