<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { auth } from '$lib/auth/auth.svelte';
	import { toast } from '$lib/toast.svelte';
	import { ApiError } from '$lib/api/client';
	import { LogIn } from 'lucide-svelte';

	let email = $state('');
	let password = $state('');
	let loading = $state(false);

	async function onSubmit(e: Event) {
		e.preventDefault();
		loading = true;
		try {
			await auth.login(email.trim(), password);
			await goto(resolve('/'));
		} catch (err) {
			const msg =
				err instanceof ApiError ? err.detail || err.title : err instanceof Error ? err.message : 'Login failed';
			toast.show(msg, 'error');
		} finally {
			loading = false;
		}
	}
</script>

<div class="auth-wrap">
	<div class="card">
		<div class="brand">
			<LogIn size={28} aria-hidden="true" />
			<h1>Sign in</h1>
			<p class="muted">Use your account credentials to connect to the chat API.</p>
		</div>

		<form onsubmit={onSubmit} class="form">
			<label for="email">Email</label>
			<input
				id="email"
				name="email"
				type="email"
				autocomplete="email"
				required
				bind:value={email}
			/>

			<label for="password">Password</label>
			<input
				id="password"
				name="password"
				type="password"
				autocomplete="current-password"
				required
				bind:value={password}
			/>

			<button type="submit" class="btn primary" disabled={loading}>
				{loading ? 'Signing in…' : 'Sign in'}
			</button>
		</form>

		<p class="footer">
			No account? <a href={resolve('/register')}>Create one</a>
		</p>
	</div>
</div>

<style>
	.auth-wrap {
		min-height: 100vh;
		display: grid;
		place-items: center;
		padding: 1.5rem;
	}
	.card {
		width: 100%;
		max-width: 420px;
		background: var(--surface);
		border: 1px solid var(--border);
		border-radius: var(--radius);
		padding: 2rem;
		box-shadow: var(--shadow);
	}
	.brand h1 {
		margin: 0.5rem 0 0.25rem;
		font-size: 1.5rem;
	}
	.muted {
		color: var(--muted);
		font-size: 0.9rem;
		margin: 0;
	}
	.form {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
		margin-top: 1.5rem;
	}
	label {
		font-size: 0.875rem;
		font-weight: 600;
	}
	input {
		padding: 0.6rem 0.75rem;
		border-radius: 8px;
		border: 1px solid var(--border);
		background: var(--bg);
		color: var(--text);
	}
	.btn {
		margin-top: 0.5rem;
		padding: 0.65rem 1rem;
		border-radius: 8px;
		border: 1px solid var(--border);
		background: var(--surface2);
		color: var(--text);
		cursor: pointer;
	}
	.btn.primary {
		background: var(--accent);
		border-color: var(--accent);
	}
	.btn.primary:hover:not(:disabled) {
		background: var(--accent-hover);
	}
	.btn:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}
	.footer {
		margin-top: 1.5rem;
		font-size: 0.9rem;
		color: var(--muted);
		text-align: center;
	}
</style>
