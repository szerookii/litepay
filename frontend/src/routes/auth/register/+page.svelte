<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import * as m from '$lib/paraglide/messages.js';

	interface ConfigResponse {
		allow_register: boolean;
	}

	let email = $state('');
	let password = $state('');
	let error = $state('');
	let loading = $state(false);
	let allowRegister = $state<boolean | null>(null);

	async function loadConfig(): Promise<void> {
		try {
			const res = await fetch('/api/config');
			if (res.ok) {
				const data: ConfigResponse = await res.json();
				allowRegister = data.allow_register;
			} else {
				allowRegister = false;
			}
		} catch {
			allowRegister = false;
		}
	}

	$effect(() => {
		loadConfig();
	});

	async function submit(e: SubmitEvent): Promise<void> {
		e.preventDefault();
		error = '';
		loading = true;

		try {
			const res = await fetch('/api/auth/register', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ email, password })
			});

			const data = await res.json();

			if (!res.ok) {
				error = data.message === 'email already in use' ? m.auth_register_error() : data.message;
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

<div class="flex min-h-screen items-center justify-center px-6">
	<div class="w-full max-w-sm space-y-8">
		<div class="space-y-2">
			<div class="flex items-center gap-2">
				<span class="text-sm font-medium tracking-tight text-foreground"
					>lite<span class="text-muted-foreground">pay</span></span
				>
				<Badge variant="outline">v0.1</Badge>
			</div>
			<h1 class="text-2xl font-medium text-foreground">{m.auth_register_title()}</h1>
			<p class="text-sm text-muted-foreground">{m.auth_register_subtitle()}</p>
		</div>

		{#if allowRegister === null}
			<p class="text-sm text-muted-foreground">Loading…</p>
		{:else if !allowRegister}
			<div class="ring-border space-y-2 p-6 ring-1 rounded-none">
				<p class="text-sm font-medium text-foreground">Registration disabled</p>
				<p class="text-xs text-muted-foreground">
					New account creation is currently disabled on this instance. Contact the administrator.
				</p>
			</div>
			<p class="text-center text-xs text-muted-foreground">
				{m.auth_register_have_account()}
				<a href="/auth/login" class="text-foreground hover:underline"
					>{m.auth_register_login_link()}</a
				>
			</p>
		{:else}
			<form onsubmit={submit} class="space-y-4">
				<div class="space-y-2">
					<Label for="email">{m.auth_register_email() || 'Email'}</Label>
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
						placeholder={m.auth_register_password_hint()}
						bind:value={password}
						required
						minlength={8}
						autocomplete="new-password"
					/>
				</div>

				{#if error}
					<p class="text-sm text-destructive">{error}</p>
				{/if}

				<Button type="submit" class="w-full" disabled={loading}>
					{loading ? m.auth_register_submitting() : m.auth_register_submit()}
				</Button>
			</form>

			<p class="text-center text-xs text-muted-foreground">
				{m.auth_register_have_account()}
				<a href="/auth/login" class="text-foreground hover:underline"
					>{m.auth_register_login_link()}</a
				>
			</p>
		{/if}

		<p class="text-center text-xs text-muted-foreground">
			<a href="/" class="transition-colors hover:text-foreground">← Back to home</a>
		</p>
	</div>
</div>
