import { browser } from '$app/environment';
import { redirect } from '@sveltejs/kit';

export const ssr = false;
export const prerender = false;

export function load() {
	if (browser && !localStorage.getItem('token')) {
		redirect(302, '/auth/login');
	}
}
