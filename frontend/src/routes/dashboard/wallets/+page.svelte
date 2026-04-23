<script lang="ts">
	import { onMount } from 'svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { CheckCircle, Warning, FloppyDisk, Lock } from 'phosphor-svelte';
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

	async function load(): Promise<void> {
		const token = localStorage.getItem('token');
		if (!token) return;
		const res = await fetch('/api/user/me', {
			headers: { Authorization: `Bearer ${token}` }
		});
		if (res.ok) {
			const me = await res.json();
			supportedCoins = me.supported_coins ?? [];
			webhookUrl = me.webhook_url ?? '';
		}
	}

	async function save(): Promise<void> {
		saving = true;
		error = '';
		saved = false;
		const token = localStorage.getItem('token');
		try {
			const res = await fetch('/api/user/wallets', {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${token}`
				},
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
			error = 'Network error';
		} finally {
			saving = false;
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
					<span>Ready</span>
					<Lock size={11} class="text-muted-foreground" />
				</div>
			</div>
		{/each}
		<p class="text-xs text-muted-foreground">
			Addresses are derived automatically from the server master seed. No configuration required.
		</p>
	</div>

	<Separator />

	<!-- Webhook -->
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
</div>
