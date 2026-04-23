import { browser } from '$app/environment';
import { redirect } from '@sveltejs/kit';

export const ssr = false;
export const prerender = false;

export async function load({ fetch }: { fetch: typeof globalThis.fetch }) {
	if (!browser) return;
	const res = await fetch('/api/user/me');
	if (res.status === 401) {
		redirect(302, '/auth/login');
	}
}
