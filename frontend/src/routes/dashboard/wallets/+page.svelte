<script lang="ts">
	import { onMount } from 'svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { CheckCircle, Warning, FloppyDisk, Lock, Eye, EyeSlash, Copy, Check, ArrowsClockwise } from 'phosphor-svelte';
	import * as m from '$lib/paraglide/messages.js';

	const coinColor: Record<string, string> = {
		BTC: '#F7931A',
		LTC: '#BEBEBE',
		SOL: '#9945FF'
	};

	const coinLabel: Record<string, string> = {
		BTC: 'Bitcoin',
		LTC: 'Litecoin',
		SOL: 'Solana'
	};

	const coinIcon: Record<string, string> = {
		BTC: 'B',
		LTC: 'L',
		SOL: 'S'
	};

	let supportedCoins = $state<string[]>([]);
	let webhookUrl = $state('');
	let saving = $state(false);
	let saved = $state(false);
	let error = $state('');

	// Webhook secret state
	let webhookSecret = $state('');
	let secretRevealed = $state(false);
	let secretCopied = $state(false);
	let confirmRotate = $state(false);
	let rotating = $state(false);
	let rotated = $state(false);
	let rotateError = $state('');

	const maskedSecret = $derived(
		webhookSecret ? `${webhookSecret.slice(0, 8)}${'•'.repeat(32)}${webhookSecret.slice(-8)}` : ''
	);
	const displaySecret = $derived(secretRevealed ? webhookSecret : maskedSecret);

	async function load(): Promise<void> {
		const res = await fetch('/api/user/me');
		if (res.ok) {
			const me = await res.json();
			supportedCoins = me.supported_coins ?? [];
			webhookUrl = me.webhook_url ?? '';
			webhookSecret = me.webhook_secret ?? '';
		}
	}

	async function save(): Promise<void> {
		saving = true;
		error = '';
		saved = false;
		try {
			const res = await fetch('/api/user/wallets', {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ webhook_url: webhookUrl || null })
			});
			if (res.ok) {
				saved = true;
				setTimeout(() => (saved = false), 3000);
			} else {
				const data = await res.json();
				error = data.message ?? m.wallets_error();
			}
		} catch {
			error = m.wallets_network_error();
		} finally {
			saving = false;
		}
	}

	function copySecret(): void {
		navigator.clipboard.writeText(webhookSecret);
		secretCopied = true;
		setTimeout(() => (secretCopied = false), 2000);
	}

	async function rotateSecret(): Promise<void> {
		if (!confirmRotate) {
			confirmRotate = true;
			return;
		}
		rotating = true;
		rotateError = '';
		rotated = false;
		confirmRotate = false;
		try {
			const res = await fetch('/api/user/webhook-secret', { method: 'POST' });
			if (res.ok) {
				const data = await res.json();
				webhookSecret = data.webhook_secret;
				secretRevealed = true;
				rotated = true;
				setTimeout(() => (rotated = false), 4000);
			} else {
				rotateError = 'Failed to rotate secret';
			}
		} catch {
			rotateError = m.wallets_network_error();
		} finally {
			rotating = false;
		}
	}

	onMount(load);
</script>

<div class="max-w-2xl space-y-8">
	<div class="space-y-1">
		<p class="text-xs uppercase tracking-widest text-muted-foreground">{m.wallets_title()}</p>
		<h1 class="text-xl font-medium text-foreground">{m.wallets_title()}</h1>
		<p class="text-xs text-muted-foreground">{m.wallets_subtitle()}</p>
	</div>

	<!-- Supported coins (read-only, managed by server) -->
	<div class="space-y-3">
		{#each supportedCoins as sym}
			<div class="ring-border flex items-center justify-between p-4 ring-1">
				<div class="flex items-center gap-2">
					<span
						class="inline-flex size-5 items-center justify-center rounded-full text-[9px] font-bold"
						style="background:{coinColor[sym] ?? '#888'}18;color:{coinColor[sym] ?? '#888'};outline:1px solid {coinColor[sym] ?? '#888'}40"
					>{coinIcon[sym] ?? sym[0]}</span>
					<span class="text-sm font-medium text-foreground">{coinLabel[sym] ?? sym}</span>
					<Badge variant="outline" class="text-[10px]">{sym}</Badge>
				</div>
				<div class="flex items-center gap-1.5 text-xs text-green-400">
					<CheckCircle size={12} weight="fill" />
					<span>{m.wallets_ready_status()}</span>
					<Lock size={11} class="text-muted-foreground" />
				</div>
			</div>
		{/each}
		<p class="text-xs text-muted-foreground">
			{m.wallets_auto_derive()}
		</p>
	</div>

	<Separator />

	<!-- Webhook URL -->
	<div class="space-y-2">
		<Label for="webhook" class="text-xs">{m.wallets_webhook_url()}</Label>
		<Input
			id="webhook"
			type="url"
			placeholder="https://yourstore.com/webhooks/litepay"
			bind:value={webhookUrl}
			class="text-xs"
		/>
		<p class="text-xs text-muted-foreground">{m.wallets_webhook_desc()}</p>
	</div>

	{#if error}
		<div class="flex items-center gap-2 text-xs text-destructive">
			<Warning size={12} />
			{error}
		</div>
	{/if}

	<div class="flex items-center gap-3">
		<Button onclick={save} disabled={saving} class="gap-1.5">
			<FloppyDisk size={13} />
			{saving ? m.wallets_saving() : m.wallets_save()}
		</Button>
		{#if saved}
			<div class="flex items-center gap-1.5 text-xs text-green-400">
				<CheckCircle size={12} weight="fill" />
				{m.wallets_success()}
			</div>
		{/if}
	</div>

	<Separator />

	<!-- Webhook Secret -->
	<div class="space-y-4">
		<div class="space-y-1">
			<p class="text-xs font-medium text-foreground">{m.wallets_webhook_secret_label()}</p>
			<p class="text-xs text-muted-foreground">{m.wallets_webhook_secret_desc()}</p>
		</div>

		{#if webhookSecret}
			<div class="ring-border flex items-center gap-2 bg-muted/20 px-3 py-2.5 ring-1">
				<span class="flex-1 truncate font-mono text-xs text-foreground">{displaySecret}</span>
				<div class="flex shrink-0 items-center gap-1">
					<button
						onclick={() => (secretRevealed = !secretRevealed)}
						class="p-1 text-muted-foreground transition-colors hover:text-foreground"
						aria-label={secretRevealed ? m.wallets_webhook_secret_hide() : m.wallets_webhook_secret_reveal()}
					>
						{#if secretRevealed}
							<EyeSlash size={13} />
						{:else}
							<Eye size={13} />
						{/if}
					</button>
					<button
						onclick={copySecret}
						class="p-1 text-muted-foreground transition-colors hover:text-foreground"
						aria-label={m.wallets_webhook_secret_copy()}
					>
						{#if secretCopied}
							<Check size={13} color="#4ade80" />
						{:else}
							<Copy size={13} />
						{/if}
					</button>
				</div>
			</div>

			{#if secretCopied}
				<p class="text-xs text-green-400">{m.wallets_webhook_secret_copied()}</p>
			{/if}
		{/if}

		{#if rotated}
			<div class="flex items-center gap-1.5 text-xs text-green-400">
				<CheckCircle size={12} weight="fill" />
				{m.wallets_webhook_secret_rotated()}
			</div>
		{/if}

		{#if rotateError}
			<p class="flex items-center gap-1.5 text-xs text-destructive">
				<Warning size={12} />
				{rotateError}
			</p>
		{/if}

		{#if confirmRotate}
			<div class="ring-border bg-destructive/5 p-3 ring-1 ring-destructive/20">
				<p class="mb-3 text-xs text-foreground">{m.wallets_webhook_secret_warning()}</p>
				<div class="flex gap-2">
					<Button
						variant="destructive"
						size="sm"
						onclick={rotateSecret}
						disabled={rotating}
						class="gap-1.5"
					>
						<ArrowsClockwise size={12} />
						{rotating ? m.wallets_webhook_secret_rotating() : m.wallets_webhook_secret_rotate()}
					</Button>
					<Button variant="ghost" size="sm" onclick={() => (confirmRotate = false)}>{m.transactions_cancel()}</Button>
				</div>
			</div>
		{:else}
			<Button variant="outline" size="sm" onclick={rotateSecret} class="gap-1.5">
				<ArrowsClockwise size={12} />
				{m.wallets_webhook_secret_rotate()}
			</Button>
		{/if}
	</div>
</div>
