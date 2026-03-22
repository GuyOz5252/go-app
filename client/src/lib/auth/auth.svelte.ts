import { api } from '$lib/api/client';
import type { User } from '$lib/api/types';
import { clearToken, loadToken, parseJwtPayload, saveToken } from './session';

class AuthStore {
	token = $state<string | null>(loadToken());
	user = $state<User | null>(null);

	async hydrate() {
		const t = this.token;
		if (!t) {
			this.user = null;
			return;
		}
		const claims = parseJwtPayload(t);
		const uid = claims?.userId;
		if (!uid) {
			this.user = null;
			return;
		}
		try {
			this.user = await api.getUser(t, uid);
		} catch {
			this.user = null;
			this.token = null;
			clearToken();
		}
	}

	async login(email: string, password: string) {
		const { token: t } = await api.login({ email, password });
		this.token = t;
		saveToken(t);
		await this.hydrate();
	}

	async register(username: string, email: string, password: string) {
		await api.register({ username, email, password });
	}

	logout() {
		this.token = null;
		this.user = null;
		clearToken();
	}
}

export const auth = new AuthStore();
