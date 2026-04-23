<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import * as m from '$lib/paraglide/messages.js';
	import SEO from '$lib/components/seo.svelte';
	import { getLocale } from '$lib/paraglide/runtime.js';

	let email = $state('');
	let password = $state('');
	let error = $state('');
	let loading = $state(false);
	let version = $state('v0.1');

	$effect(() => {
		fetch('/api/config')
			.then((r) => r.json())
			.then((d) => {
				version = d.version ? `v${d.version}` : 'v0.1';
			})
			.catch(() => {});
	});

	async function submit(e: SubmitEvent): Promise<void> {
		e.preventDefault();
		error = '';
		loading = true;

		try {
			const res = await fetch('/api/auth/login', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ email, password })
			});

			const data = await res.json();

			if (!res.ok) {
				error = data.message === 'Invalid credentials' ? m.auth_login_error() : data.message;
				return;
			}

			localStorage.setItem('token', data.token);
			goto('/dashboard');
		} catch {
			error = 'Network error';
		} finally {
			loading = false;
	}
}
</script>

<SEO
	config={{
		title: m.seo_login_title(),
		description: m.seo_login_desc(),
		ogType: 'website'
	}}
	url={page.url.pathname}
	locale={getLocale()}
/>

<div class="flex min-h-screen items-center justify-center px-6">
	<div class="w-full max-w-sm space-y-8">
		<div class="space-y-2">
			<div class="flex items-center gap-2">
				<span class="text-sm font-medium tracking-tight text-foreground"
					>lite<span class="text-muted-foreground">pay</span></span
				>
				<Badge variant="outline">{version}</Badge>
			</div>
			<h1 class="text-2xl font-medium text-foreground">{m.auth_login_title()}</h1>
			<p class="text-sm text-muted-foreground">{m.auth_login_subtitle()}</p>
		</div>

		<form onsubmit={submit} class="space-y-4">
			<div class="space-y-2">
				<Label for="email">{m.auth_login_email()}</Label>
				<Input
					id="email"
					type="email"
					placeholder="you@example.com"
					bind:value={email}
					required
					autocomplete="email"
				/>
			</div>

			<div class="space-y-2">
				<Label for="password">{m.auth_login_password()}</Label>
				<Input
					id="password"
					type="password"
					placeholder="••••••••"
					bind:value={password}
					required
					autocomplete="current-password"
				/>
			</div>

			{#if error}
				<p class="text-sm text-destructive">{error}</p>
			{/if}

			<Button type="submit" class="w-full" disabled={loading}>
				{loading ? m.auth_login_submitting() : m.auth_login_submit()}
			</Button>
		</form>

		<div class="flex flex-col items-center gap-4 pt-2">
			<p class="text-xs text-muted-foreground">
				{m.auth_login_no_account()}
				<a href="/auth/register" class="text-foreground hover:underline"
					>{m.auth_login_register_link()}</a
				>
			</p>

			<p class="text-center text-xs text-muted-foreground">
				<a href="/" class="transition-colors hover:text-foreground">← Back to home</a>
			</p>
		</div>
	</div>
</div>
