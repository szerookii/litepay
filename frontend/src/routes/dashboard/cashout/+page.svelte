<script lang="ts">
	import { onMount } from 'svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { ArrowUpRight, CheckCircle, Warning, ArrowRight, Spinner } from 'phosphor-svelte';
	import * as m from '$lib/paraglide/messages.js';

	interface CashoutTx {
		from_address: string;
		tx_hash: string;
		amount: number;
	}

	const coinColor: Record<string, string> = {
		BTC: '#F7931A',
		LTC: '#BEBEBE',
		SOL: '#9945FF'
	};

	function getExplorerUrl(symbol: string, hash: string) {
		switch (symbol) {
			case 'BTC': return `https://mempool.space/tx/${hash}`;
			case 'LTC': return `https://blockchair.com/litecoin/transaction/${hash}`;
			case 'SOL': return `https://solscan.io/tx/${hash}`;
			default: return '#';
		}
	}

	let supportedCoins = $state<string[]>([]);
	let balance = $state<Record<string, number>>({});
	let balanceLoading = $state(true);

	async function loadData(): Promise<void> {
		const token = localStorage.getItem('token');
		if (!token) return;
		try {
			const [meRes, balRes] = await Promise.all([
				fetch('/api/user/me', { headers: { Authorization: `Bearer ${token}` } }),
				fetch('/api/user/balance', { headers: { Authorization: `Bearer ${token}` } })
			]);
			if (meRes.ok) {
				const me = await meRes.json();
				supportedCoins = me.supported_coins ?? [];
				if (supportedCoins.length > 0) symbol = supportedCoins[0];
			}
			if (balRes.ok) balance = await balRes.json();
		} catch {
			// silent
		} finally {
			balanceLoading = false;
		}
	}

	onMount(loadData);

	let symbol = $state('');
	let destination = $state('');
	let loading = $state(false);
	let error = $state('');
	let results = $state<CashoutTx[]>([]);
	let done = $state(false);

	async function submit() {
		if (!destination.trim()) { error = 'Destination address required'; return; }
		loading = true;
		error = '';
		done = false;
		results = [];

		const token = localStorage.getItem('token');
		try {
			const res = await fetch('/api/user/cashout', {
				method: 'POST',
				headers: { Authorization: `Bearer ${token}`, 'Content-Type': 'application/json' },
				body: JSON.stringify({
					symbol,
					destination: destination.trim()
				})
			});
			const data = await res.json();
			if (!res.ok) {
				error = data.message ?? 'Cashout failed';
				if (data.details) error += ' — ' + data.details.map((d: {address: string; reason: string}) => `${d.address.slice(0,8)}…: ${d.reason}`).join('; ');
			} else {
				results = data.transactions ?? [];
				done = true;
				loadData();
			}
		} catch {
			error = 'Network error';
		} finally {
			loading = false;
		}
	}

	function reset() {
		done = false;
		results = [];
		destination = '';
		amount = '';
		loadData();
	}
</script>

<div class="max-w-lg space-y-8">
	<div class="space-y-1">
		<p class="text-xs uppercase tracking-widest text-muted-foreground">{m.dash_nav_dashboard()}</p>
		<h1 class="text-xl font-medium text-foreground">Cashout</h1>
		<p class="text-xs text-muted-foreground">Sweep received funds to your personal wallet.</p>
	</div>

	<!-- Balances -->
	<div class="space-y-2">
		<p class="text-xs uppercase tracking-widest text-muted-foreground">Available balance</p>
		<div class="grid gap-2" style="grid-template-columns: repeat({Math.min(supportedCoins.length || 3, 4)}, minmax(0, 1fr))">
			{#each (balanceLoading ? ['BTC', 'LTC', 'SOL'] : supportedCoins) as sym}
				<div class="ring-1 ring-border p-3 space-y-1 {balanceLoading ? 'animate-pulse' : ''}">
					<div class="flex items-center gap-1.5">
						<span
							class="inline-flex size-3.5 items-center justify-center rounded-full text-[7px] font-bold"
							style="background:{(coinColor[sym] ?? '#888')}18;color:{coinColor[sym] ?? '#888'};outline:1px solid {(coinColor[sym] ?? '#888')}40"
						>{sym[0]}</span>
						<span class="text-[10px] text-muted-foreground">{sym}</span>
					</div>
					<p class="font-mono text-sm font-medium text-foreground">
						{balanceLoading ? '…' : ((balance[sym] ?? 0).toFixed(8).replace(/\.?0+$/, '') || '0')}
					</p>
				</div>
			{/each}
		</div>
		<p class="text-[10px] text-muted-foreground">Confirmed PAID payments only</p>
	</div>

	{#if done}
		<!-- Success state -->
		<div class="space-y-4">
			<div class="flex items-center gap-2 text-sm font-medium text-foreground">
				<CheckCircle size={16} color="#4ade80" weight="fill" />
				{results.length} transaction{results.length !== 1 ? 's' : ''} broadcast
			</div>
			<div class="space-y-2">
				{#each results as tx}
					<div class="ring-1 ring-border p-4 space-y-1 bg-background">
						<div class="flex items-center justify-between">
							<span class="text-xs text-muted-foreground font-mono">{tx.from_address.slice(0, 16)}…</span>
							<span class="text-xs font-medium">{tx.amount.toFixed(6)} {symbol}</span>
						</div>
						<a
							href={getExplorerUrl(symbol, tx.tx_hash)}
							target="_blank"
							rel="noopener noreferrer"
							class="flex items-center gap-1 text-xs font-mono text-muted-foreground hover:text-foreground transition-colors"
						>
							{tx.tx_hash.slice(0, 20)}…
							<ArrowUpRight size={10} />
						</a>
					</div>
				{/each}
			</div>
			<Button variant="outline" size="sm" class="text-xs" onclick={reset}>New cashout</Button>
		</div>
	{:else}
		<!-- Form -->
		<div class="space-y-6">
			<!-- Asset selector -->
			<div class="space-y-2">
				<Label class="text-xs">Asset</Label>
				<div class="grid gap-2" style="grid-template-columns: repeat({Math.min(supportedCoins.length || 3, 4)}, minmax(0, 1fr))">
					{#each supportedCoins as sym}
						<button
							onclick={() => (symbol = sym)}
							class="ring-1 ring-border p-3 text-center transition-colors
								{symbol === sym ? 'bg-muted ring-foreground' : 'hover:bg-muted/50'}"
						>
							<span
								class="mb-1 inline-flex size-5 items-center justify-center rounded-full text-[9px] font-bold"
								style="background:{(coinColor[sym] ?? '#888')}18;color:{coinColor[sym] ?? '#888'};outline:1px solid {(coinColor[sym] ?? '#888')}40"
							>{sym[0]}</span>
							<p class="text-[10px] font-medium text-foreground">{sym}</p>
						</button>
					{/each}
				</div>
			</div>

			<!-- Destination -->
			<div class="space-y-2">
				<Label class="text-xs">Destination address</Label>
				<Input
					placeholder="Your {symbol} wallet address"
					bind:value={destination}
					class="text-xs font-mono"
				/>
			</div>

			<!-- Warning -->
			<div class="flex gap-2 ring-1 ring-border p-3 text-xs text-muted-foreground">
				<Warning size={14} class="mt-0.5 shrink-0 text-yellow-400" />
				<span>This sends funds on-chain. Double-check the destination address — transactions cannot be reversed.</span>
			</div>

			{#if error}
				<p class="text-xs text-destructive">{error}</p>
			{/if}

			<Button
				onclick={submit}
				disabled={loading || !destination.trim() || !symbol}
				class="w-full gap-1.5"
			>
				{#if loading}
					<Spinner size={14} class="animate-spin" />
					Broadcasting…
				{:else}
					Cashout {symbol}
					<ArrowRight size={14} />
				{/if}
			</Button>
		</div>
	{/if}
</div>
