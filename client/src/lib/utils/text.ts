const URL_RE = /(https?:\/\/[^\s<]+[^<.,:;"')\]\s])/g;

export function escapeHtml(s: string): string {
	return s
		.replace(/&/g, '&amp;')
		.replace(/</g, '&lt;')
		.replace(/>/g, '&gt;')
		.replace(/"/g, '&quot;')
		.replace(/'/g, '&#39;');
}

/**
 * Returns HTML-safe segments: plain text chunks and anchor tags for http(s) URLs.
 */
export function linkifyToHtml(text: string): string {
	const escaped = escapeHtml(text);
	const parts: string[] = [];
	let last = 0;
	let m: RegExpExecArray | null;
	const re = new RegExp(URL_RE.source, 'g');
	while ((m = re.exec(escaped)) !== null) {
		parts.push(escaped.slice(last, m.index));
		const url = m[1];
		parts.push(
			`<a href="${url}" target="_blank" rel="noopener noreferrer" class="text-link">${url}</a>`
		);
		last = m.index + m[0].length;
	}
	parts.push(escaped.slice(last));
	return parts.join('');
}

export function isProbablyImageUrl(url: string): boolean {
	if (!/^https?:\/\//i.test(url)) return false;
	return /\.(png|jpe?g|gif|webp|svg|avif)(\?|$)/i.test(url);
}
