// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces

declare global {
	interface ImportMetaEnv {
		readonly PUBLIC_API_BASE: string;
	}

	namespace App {
		// interface Error {}
		// interface Locals {}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
