import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, loadEnv } from 'vite';
import type { ProxyOptions } from 'vite';
import type HttpProxy from 'http-proxy';
import type { ClientRequest, IncomingMessage } from 'node:http';

function configureApiProxy(proxy: HttpProxy) {
	proxy.on('proxyReq', (proxyReq: ClientRequest, req: IncomingMessage) => {
		const raw = req.url;
		if (!raw?.includes('access_token=')) return;
		try {
			const qIndex = raw.indexOf('?');
			const path = qIndex >= 0 ? raw.slice(0, qIndex) : raw;
			const qs = qIndex >= 0 ? raw.slice(qIndex + 1) : '';
			const params = new URLSearchParams(qs);
			const token = params.get('access_token');
			if (token) {
				proxyReq.setHeader('Authorization', `Bearer ${decodeURIComponent(token)}`);
			}
			params.delete('access_token');
			const nextQs = params.toString();
			proxyReq.path = nextQs ? `${path}?${nextQs}` : path;
		} catch {
			/* ignore */
		}
	});
	proxy.on('proxyReqWs', (proxyReq: ClientRequest, req: IncomingMessage) => {
		const raw = req.url;
		if (!raw?.includes('access_token=')) return;
		try {
			const qIndex = raw.indexOf('?');
			const path = qIndex >= 0 ? raw.slice(0, qIndex) : raw;
			const qs = qIndex >= 0 ? raw.slice(qIndex + 1) : '';
			const params = new URLSearchParams(qs);
			const token = params.get('access_token');
			if (token) {
				proxyReq.setHeader('Authorization', `Bearer ${decodeURIComponent(token)}`);
			}
			params.delete('access_token');
			const nextQs = params.toString();
			proxyReq.path = nextQs ? `${path}?${nextQs}` : path;
		} catch {
			/* ignore */
		}
	});
}

export default defineConfig(({ mode }) => {
	const env = loadEnv(mode, process.cwd(), '');
	const target = env.PUBLIC_DEV_PROXY_TARGET || 'http://localhost:8080';

	const apiProxy: ProxyOptions = {
		target,
		changeOrigin: true,
		ws: true,
		configure: configureApiProxy as NonNullable<ProxyOptions['configure']>
	};

	return {
		plugins: [sveltekit()],
		server: {
			proxy: {
				'/api': apiProxy
			}
		},
		test: {
			include: ['src/**/*.{test,spec}.{js,ts}']
		}
	};
});
