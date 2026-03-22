export type ToastKind = 'error' | 'info';

class ToastStore {
	message = $state<string | null>(null);
	kind = $state<ToastKind>('info');

	show(text: string, kind: ToastKind = 'info') {
		this.message = text;
		this.kind = kind;
		window.setTimeout(() => {
			this.message = null;
		}, 6000);
	}
}

export const toast = new ToastStore();
