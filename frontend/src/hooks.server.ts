
import type { Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
	const cookie = event.request.headers.get('cookie') ?? '';
	const pathname = event.url.pathname;

	if (pathname.startsWith('/admin')) {
		try {
			const res = await fetch('http://localhost:8080/api/v1/auth/me', {
				headers: { cookie },
				credentials: 'include'
			});

			const contentType = res.headers.get('content-type') ?? '';
			const isJSON = contentType.includes('application/json');
			const raw = await res.text();


			if (res.ok && isJSON) {
				try {
					const json = JSON.parse(raw);
					console.log('[hooks] user fetched:', json.data);
					event.locals.user = json.data ?? null;
				} catch (err) {
					console.error('[hooks] JSON parse error:', err);
					event.locals.user = null;
				}
			} else {
				console.warn('[hooks] non-OK or non-JSON:', res.status);
				event.locals.user = null;
			}
		} catch (err) {
			console.error('[hooks] request failed:', err);
			event.locals.user = null;
		}
	} else {
		event.locals.user = null; // safe fallback for public pages
	}

	return resolve(event);
};

