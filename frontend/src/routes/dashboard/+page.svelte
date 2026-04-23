<script lang="ts">
	import { onMount } from 'svelte';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import { ArrowRight, Wallet, Key, CheckCircle } from 'phosphor-svelte';
	import * as m from '$lib/paraglide/messages.js';

	interface Me {
		id: string;
		email: string;
		api_key: string;
		webhook_url: string | null;
		account_index: number;
		supported_coins: string[];
	}

	let me = $state<Me | null>(null);
	let error = $state('');

	const coinColor: Record<string, string> = {
		BTC: '#F7931A',
		LTC: '#BEBEBE',
		SOL: '#9945FF'
	};

	async function load(): Promise<void> {
		const token = localStorage.getItem('token');
		if (!token) return;
		try {
			const meRes = await fetch('/api/user/me', { headers: { Authorization: `Bearer ${token}` } });
			if (meRes.ok) me = await meRes.json();
			else error = 'Failed to load account info';
		} catch {
			error = 'Network error';
		}
	}

	onMount(load);
</script>

<div class="max-w-2xl space-y-8">
	<div class="space-y-1">
		<p class="text-xs uppercase tracking-widest text-muted-foreground">{m.dash_nav_dashboard()}</p>
		<h1 class="text-xl font-medium text-foreground">{m.dash_nav_dashboard()}</h1>
	</div>

	{#if error}
		<p class="text-sm text-destructive">{error}</p>
	{:else if !me}
		<div class="space-y-3">
			{#each [1, 2, 3] as _}
				<div class="bg-muted/30 h-20 animate-pulse ring-1 ring-border"></div>
			{/each}
		</div>
	{:else}
		<!-- Account card -->
		<div class="ring-border space-y-4 p-5 ring-1">
			<div class="flex items-start justify-between">
				<div>
					<p class="text-xs text-muted-foreground">{m.dash_overview_account()}</p>
					<p class="mt-0.5 text-sm font-medium text-foreground">{me.email}</p>
				</div>
				<Badge variant="secondary">
					{me.supported_coins.length} coins
				</Badge>
			</div>
			<Separator />
			<div class="flex items-center justify-between">
				<p class="text-xs text-muted-foreground">{m.dash_overview_user_id()}</p>
				<span class="font-mono text-xs text-foreground">{me.id.slice(0, 8)}…</span>
			</div>
		</div>

		<!-- Coins status -->
		<div class="space-y-3">
			<p class="text-xs uppercase tracking-widest text-muted-foreground">
				{m.dash_overview_wallet_config()}
			</p>
			<div class="grid gap-3 sm:grid-cols-3">
				{#each me.supported_coins as sym}
					<div class="ring-border space-y-3 p-4 ring-1">
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-2">
								<span
									class="inline-flex size-5 items-center justify-center rounded-full text-[9px] font-bold"
									style="background:{(coinColor[sym] ?? '#888')}18;color:{coinColor[sym] ??
										'#888'};outline:1px solid {(coinColor[sym] ?? '#888')}40"
								>{sym[0]}</span>
								<span class="text-xs font-medium text-foreground">{sym}</span>
							</div>
							<CheckCircle size={12} color="#4ade80" weight="fill" />
						</div>
						<p class="text-xs text-muted-foreground">{m.dash_overview_ready()}</p>
					</div>
				{/each}
			</div>
		</div>

		<!-- Quick links -->
		<div class="space-y-3">
			<p class="text-xs uppercase tracking-widest text-muted-foreground">
				{m.dash_overview_quick_links()}
			</p>
			<div class="grid gap-3 sm:grid-cols-2">
				<a
					href="/dashboard/wallets"
					class="ring-border group flex items-center justify-between p-4 ring-1 transition-colors hover:bg-muted/30"
				>
					<div class="flex items-center gap-3">
						<Wallet size={16} class="text-muted-foreground group-hover:text-foreground" />
						<div>
							<p class="text-xs font-medium text-foreground">{m.dash_nav_wallets()}</p>
							<p class="text-xs text-muted-foreground">{m.dash_overview_wallets_desc()}</p>
						</div>
					</div>
					<ArrowRight size={12} class="text-muted-foreground" />
				</a>
				<a
					href="/dashboard/api-keys"
					class="ring-border group flex items-center justify-between p-4 ring-1 transition-colors hover:bg-muted/30"
				>
					<div class="flex items-center gap-3">
						<Key size={16} class="text-muted-foreground group-hover:text-foreground" />
						<div>
							<p class="text-xs font-medium text-foreground">{m.dash_nav_api_keys()}</p>
							<p class="text-xs text-muted-foreground">{m.dash_overview_api_keys_desc()}</p>
						</div>
					</div>
					<ArrowRight size={12} class="text-muted-foreground" />
				</a>
			</div>
		</div>
	{/if}
</div>
