<script lang="ts">
	import { onMount } from 'svelte';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import {
		Receipt,
		CheckCircle,
		Clock,
		XCircle,
		ArrowSquareOut,
		ArrowsClockwise,
		Plus,
		ArrowRight,
		MagnifyingGlass,
		CaretLeft,
		CaretRight,
		ArrowUUpLeft,
		ArrowUpRight,
		Spinner
	} from 'phosphor-svelte';
	import * as m from '$lib/paraglide/messages.js';

	interface Payment {
		id: string;
		wallet_address: string;
		amount_crypto: number;
		currency_crypto: string;
		amount_fiat: number;
		currency_fiat: string;
		status: 'PENDING' | 'CONFIRMING' | 'PAID' | 'EXPIRED' | 'REFUNDED' | 'CASHED_OUT';
		transaction_hash?: string;
		create_time: string;
		expires_at: string;
	}

	function resolveStatus(p: Payment): Payment['status'] {
		if (p.status === 'PAID') return 'PAID';
		if (p.status === 'REFUNDED') return 'REFUNDED';
		if (p.status === 'CASHED_OUT') return 'CASHED_OUT';
		if (p.status === 'EXPIRED') return 'EXPIRED';
		if (new Date(p.expires_at) < new Date()) return 'EXPIRED';
		return p.status;
	}

	let payments = $state<Payment[]>([]);
	let supportedCoins = $state<string[]>([]);
	let loading = $state(true);
	let refreshing = $state(false);
	let error = $state('');

	// Refund state
	let refundPayment = $state<Payment | null>(null);
	let refunding = $state(false);
	let refundError = $state('');
	let refundResult = $state<{ tx_hash: string; to: string } | null>(null);

	async function doRefund() {
		if (!refundPayment) return;
		refunding = true;
		refundError = '';
		refundResult = null;
		try {
			const res = await fetch(`/api/payment/${refundPayment.id}/refund`, { method: 'POST' });
			const data = await res.json();
			if (!res.ok) {
				refundError = data.message ?? m.transactions_refund_failed();
			} else {
				refundResult = { tx_hash: data.tx_hash, to: data.to };
				load(true);
			}
		} catch {
			refundError = 'Network error';
		} finally {
			refunding = false;
		}
	}

	function closeRefundDialog() {
		refundPayment = null;
		refundError = '';
		refundResult = null;
	}

	// Filters
	let search = $state('');
	let filterStatus = $state('ALL');
	let filterCurrency = $state('ALL');

	// Pagination
	const PAGE_SIZE = 10;
	let page = $state(1);

	// Create payment state
	let showCreate = $state(false);
	let createAmount = $state('');
	let createCurrency = $state('USD');
	let createSymbol = $state('');
	let creating = $state(false);
	let createError = $state('');

	const filtered = $derived(
		payments.filter((p) => {
			const status = resolveStatus(p);
			if (filterStatus !== 'ALL' && status !== filterStatus) return false;
			if (filterCurrency !== 'ALL' && p.currency_crypto !== filterCurrency) return false;
			if (search) {
				const q = search.toLowerCase();
				if (
					!p.transaction_hash?.toLowerCase().includes(q) &&
					!p.wallet_address.toLowerCase().includes(q) &&
					!p.id.toLowerCase().includes(q)
				)
					return false;
			}
			return true;
		})
	);

	const totalPages = $derived(Math.max(1, Math.ceil(filtered.length / PAGE_SIZE)));
	const paginated = $derived(filtered.slice((page - 1) * PAGE_SIZE, page * PAGE_SIZE));

	$effect(() => {
		// Reset to page 1 when filters change
		filtered;
		page = 1;
	});

	async function load(isRefresh = false): Promise<void> {
		if (isRefresh) refreshing = true;
		else loading = true;

		try {
			const [paymentsRes, meRes] = await Promise.all([
				fetch('/api/user/payments'),
				supportedCoins.length === 0 ? fetch('/api/user/me') : Promise.resolve(null)
			]);
			if (paymentsRes.ok) {
				payments = await paymentsRes.json();
			} else {
				error = m.transactions_failed_load();
			}
			if (meRes?.ok) {
				const me = await meRes.json();
				supportedCoins = me.supported_coins ?? [];
				if (!createSymbol && supportedCoins.length > 0) createSymbol = supportedCoins[0];
			}
		} catch {
			error = m.transactions_network_error();
		} finally {
			loading = false;
			refreshing = false;
		}
	}

	async function createPayment() {
		creating = true;
		createError = '';
		try {
			const meRes = await fetch('/api/user/me');
			const me = await meRes.json();

			const res = await fetch('/api/payment', {
				method: 'POST',
				headers: {
					Authorization: `Bearer ${me.api_key}`,
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					amount: parseFloat(createAmount),
					currency: createCurrency,
					symbol: createSymbol
				})
			});

			if (res.ok) {
				const data = await res.json();
				showCreate = false;
				createAmount = '';
				load(true);
				window.open(`/pay?id=${data.id}`, '_blank');
			} else {
				const data = await res.json();
				createError = data.message ?? m.transactions_create_error();
			}
		} catch {
			createError = 'Network error';
		} finally {
			creating = false;
		}
	}

	onMount(load);

	function formatDate(dateStr: string) {
		return new Intl.DateTimeFormat(undefined, {
			dateStyle: 'medium',
			timeStyle: 'short'
		}).format(new Date(dateStr));
	}

	const coinColor: Record<string, string> = {
		BTC: '#F7931A',
		LTC: '#BEBEBE',
		SOL: '#9945FF'
	};

	function getExplorerUrl(symbol: string, hash: string) {
		switch (symbol) {
			case 'BTC':
				return `https://mempool.space/tx/${hash}`;
			case 'LTC':
				return `https://blockchair.com/litecoin/transaction/${hash}`;
			case 'SOL':
				return `https://solscan.io/tx/${hash}`;
			default:
				return '#';
		}
	}

	const STATUS_OPTIONS = [
		{ value: 'ALL', label: m.transactions_status_all() },
		{ value: 'PENDING', label: m.transactions_status_pending() },
		{ value: 'CONFIRMING', label: m.transactions_status_confirming() },
		{ value: 'PAID', label: m.transactions_status_confirmed() },
		{ value: 'REFUNDED', label: m.transactions_status_refunded() },
		{ value: 'CASHED_OUT', label: m.transactions_status_cashed_out() },
		{ value: 'EXPIRED', label: m.transactions_status_expired() }
	];

	const CURRENCY_OPTIONS = $derived([
		{ value: 'ALL', label: m.transactions_all_assets() },
		...supportedCoins.map((c) => ({ value: c, label: c }))
	]);
</script>

<div class="max-w-5xl space-y-6">
	<div class="flex items-start justify-between">
		<div class="space-y-1">
			<p class="text-xs uppercase tracking-widest text-muted-foreground">{m.dash_nav_dashboard()}</p>
			<h1 class="text-xl font-medium text-foreground">{m.transactions_title()}</h1>
			<p class="text-xs text-muted-foreground">{m.transactions_subtitle()}</p>
		</div>
		<div class="flex gap-2">
			<Button
				variant="outline"
				size="sm"
				onclick={() => load(true)}
				disabled={refreshing || loading}
				class="gap-1.5 text-xs"
			>
				<ArrowsClockwise size={14} class={refreshing ? 'animate-spin' : ''} />
				{m.transactions_refresh()}
			</Button>
			<Button size="sm" onclick={() => (showCreate = true)} class="gap-1.5 text-xs">
				<Plus size={14} />
				{m.transactions_create()}
			</Button>
		</div>
	</div>

	<!-- Filters -->
	<div class="flex flex-wrap gap-2">
		<div class="relative flex-1 min-w-48">
			<MagnifyingGlass size={13} class="absolute left-2.5 top-1/2 -translate-y-1/2 text-muted-foreground" />
			<Input
				placeholder={m.transactions_search_placeholder()}
				bind:value={search}
				class="pl-7 text-xs h-8"
			/>
		</div>

		<select
			bind:value={filterStatus}
			class="h-8 rounded-md border border-input bg-background px-2.5 py-1 text-xs shadow-sm focus:outline-none focus:ring-1 focus:ring-ring"
		>
			{#each STATUS_OPTIONS as opt}
				<option value={opt.value}>{opt.label}</option>
			{/each}
		</select>

		<select
			bind:value={filterCurrency}
			class="h-8 rounded-md border border-input bg-background px-2.5 py-1 text-xs shadow-sm focus:outline-none focus:ring-1 focus:ring-ring"
		>
			{#each CURRENCY_OPTIONS as opt}
				<option value={opt.value}>{opt.label}</option>
			{/each}
		</select>

		{#if search || filterStatus !== 'ALL' || filterCurrency !== 'ALL'}
			<Button
				variant="ghost"
				size="sm"
				class="h-8 text-xs text-muted-foreground"
				onclick={() => { search = ''; filterStatus = 'ALL'; filterCurrency = 'ALL'; }}
			>
				{m.transactions_clear_filters()}
			</Button>
		{/if}
	</div>

	{#if error}
		<p class="text-sm text-destructive">{error}</p>
	{:else if loading}
		<div class="space-y-3">
			{#each [1, 2, 3, 4, 5] as _}
				<div class="bg-muted/30 h-14 animate-pulse ring-1 ring-border"></div>
			{/each}
		</div>
	{:else if filtered.length === 0}
		<div class="flex flex-col items-center justify-center space-y-3 py-12 text-center ring-1 ring-border">
			<Receipt size={32} class="text-muted-foreground" />
			<p class="text-xs text-muted-foreground">
				{payments.length === 0 ? m.transactions_empty() : m.transactions_no_match()}
			</p>
		</div>
	{:else}
		<div class="ring-1 ring-border overflow-hidden bg-background">
			<div class="overflow-x-auto">
				<table class="w-full text-left text-xs">
					<thead>
						<tr class="border-b border-border bg-muted/30 font-medium uppercase tracking-tighter text-muted-foreground">
							<th class="px-4 py-3">{m.transactions_table_date()}</th>
							<th class="px-4 py-3">{m.transactions_table_status()}</th>
							<th class="px-4 py-3">{m.transactions_table_asset()}</th>
							<th class="px-4 py-3 text-right">{m.transactions_table_amount()}</th>
							<th class="px-4 py-3">{m.transactions_table_txhash()}</th>
						<th class="px-4 py-3"></th>
						</tr>
					</thead>
					<tbody class="divide-y divide-border">
						{#each paginated as p}
							<tr class="hover:bg-muted/30 transition-colors">
								<td class="whitespace-nowrap px-4 py-4 text-muted-foreground">
									{formatDate(p.create_time)}
								</td>
								<td class="px-4 py-4">
									{#if resolveStatus(p) === 'PAID'}
										<div class="flex items-center gap-1.5 text-foreground">
											<CheckCircle size={12} color="#4ade80" weight="fill" />
											<span>{m.transactions_status_confirmed()}</span>
										</div>
									{:else if resolveStatus(p) === 'CONFIRMING'}
										<div class="flex items-center gap-1.5 text-yellow-400">
											<Clock size={12} />
											<span>{m.transactions_status_confirming()}</span>
										</div>
									{:else if resolveStatus(p) === 'REFUNDED'}
										<div class="flex items-center gap-1.5 text-blue-400">
											<ArrowUUpLeft size={12} />
											<span>{m.transactions_status_refunded()}</span>
										</div>
									{:else if resolveStatus(p) === 'CASHED_OUT'}
										<div class="flex items-center gap-1.5 text-emerald-400">
											<ArrowUpRight size={12} />
											<span>{m.transactions_status_cashed_out()}</span>
										</div>
									{:else if resolveStatus(p) === 'EXPIRED'}
										<div class="flex items-center gap-1.5 text-muted-foreground/50">
											<XCircle size={12} />
											<span>{m.transactions_status_expired()}</span>
										</div>
									{:else}
										<div class="flex items-center gap-1.5 text-muted-foreground">
											<Clock size={12} />
											<span>{m.transactions_status_pending()}</span>
										</div>
									{/if}
								</td>
								<td class="px-4 py-4">
									<div class="flex items-center gap-2">
										<span
											class="inline-flex size-4 items-center justify-center rounded-full text-[8px] font-bold"
											style="background:{(coinColor[p.currency_crypto] ?? '#888')}18;color:{coinColor[p.currency_crypto] ?? '#888'};outline:1px solid {(coinColor[p.currency_crypto] ?? '#888')}40"
										>{p.currency_crypto[0]}</span>
										<span class="font-medium text-foreground">{p.currency_crypto}</span>
									</div>
								</td>
								<td class="whitespace-nowrap px-4 py-4 text-right">
									<div class="font-medium text-foreground">
										{p.amount_crypto} {p.currency_crypto}
									</div>
									<div class="text-[10px] text-muted-foreground">
										{p.amount_fiat.toFixed(2)} {p.currency_fiat}
									</div>
								</td>
								<td class="px-4 py-4">
									{#if p.transaction_hash}
										<a
											href={getExplorerUrl(p.currency_crypto, p.transaction_hash)}
											target="_blank"
											rel="noopener noreferrer"
											class="inline-flex items-center gap-1 font-mono text-muted-foreground transition-colors hover:text-foreground"
										>
											{p.transaction_hash.slice(0, 8)}…
											<ArrowSquareOut size={10} />
										</a>
									{:else}
										<span class="text-muted-foreground/30">—</span>
									{/if}
								</td>
								<td class="px-4 py-4">
									{#if resolveStatus(p) === 'PAID' && p.transaction_hash && resolveStatus(p) !== 'REFUNDED'}
										<Button
											variant="ghost"
											size="sm"
											class="h-6 gap-1 px-2 text-[10px] text-muted-foreground hover:text-foreground"
											onclick={() => { refundPayment = p; refundError = ''; refundResult = null; }}
										>
											<ArrowUUpLeft size={10} />
											{m.transactions_refund_button()}
										</Button>
									{/if}
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>

		<!-- Pagination -->
		<div class="flex items-center justify-between text-xs text-muted-foreground">
			<span>
				{filtered.length} {filtered.length === 1 ? m.transactions_result_singular() : m.transactions_result_plural()}
				{#if totalPages > 1}— {m.transactions_page_label({ current: page, total: totalPages })}{/if}
			</span>
			{#if totalPages > 1}
				<div class="flex items-center gap-1">
					<Button
						variant="outline"
						size="sm"
						class="h-7 w-7 p-0"
						disabled={page === 1}
						onclick={() => (page = Math.max(1, page - 1))}
					>
						<CaretLeft size={12} />
					</Button>
					{#each Array.from({ length: totalPages }, (_, i) => i + 1) as pg}
						{#if totalPages <= 7 || pg === 1 || pg === totalPages || Math.abs(pg - page) <= 1}
							<Button
								variant={pg === page ? 'default' : 'outline'}
								size="sm"
								class="h-7 w-7 p-0 text-xs"
								onclick={() => (page = pg)}
							>
								{pg}
							</Button>
						{:else if Math.abs(pg - page) === 2}
							<span class="px-1">…</span>
						{/if}
					{/each}
					<Button
						variant="outline"
						size="sm"
						class="h-7 w-7 p-0"
						disabled={page === totalPages}
						onclick={() => (page = Math.min(totalPages, page + 1))}
					>
						<CaretRight size={12} />
					</Button>
				</div>
			{/if}
		</div>
	{/if}
</div>

<!-- Create Payment Dialog -->
<Dialog.Root bind:open={showCreate}>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header>
			<Dialog.Title>{m.transactions_create_title()}</Dialog.Title>
		</Dialog.Header>
		<div class="space-y-4 py-4">
			<div class="grid grid-cols-2 gap-4">
				<div class="space-y-2">
					<Label class="text-xs">{m.transactions_create_amount()}</Label>
					<Input type="number" placeholder="0.00" bind:value={createAmount} class="text-xs" />
				</div>
				<div class="space-y-2">
					<Label class="text-xs">{m.transactions_create_currency()}</Label>
					<select
						bind:value={createCurrency}
						class="flex h-9 w-full rounded-md border border-input bg-background px-3 py-1 text-xs shadow-sm"
					>
						<option value="USD">USD</option>
						<option value="EUR">EUR</option>
					</select>
				</div>
			</div>
			<div class="space-y-2">
				<Label class="text-xs">{m.transactions_create_asset()}</Label>
				<div class="grid grid-cols-3 gap-2">
					{#each supportedCoins as sym}
						<button
							onclick={() => (createSymbol = sym)}
							class="ring-1 ring-border p-3 text-center transition-colors
								{createSymbol === sym ? 'bg-muted ring-foreground' : 'hover:bg-muted/50'}"
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

			{#if createError}
				<p class="text-xs text-destructive">{createError}</p>
			{/if}
		</div>
		<Dialog.Footer class="sm:justify-start">
			<Button onclick={createPayment} disabled={creating || !createAmount} class="w-full gap-1.5">
				{creating ? m.transactions_create_submit() + '...' : m.transactions_create_submit()}
				<ArrowRight size={14} />
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

<!-- Refund Dialog -->
<Dialog.Root open={!!refundPayment} onOpenChange={(o) => { if (!o) closeRefundDialog(); }}>
	<Dialog.Content class="sm:max-w-sm">
		<Dialog.Header>
			<Dialog.Title>{m.transactions_refund_title()}</Dialog.Title>
		</Dialog.Header>
		{#if refundResult}
			<div class="space-y-3 py-2">
				<p class="text-xs text-muted-foreground">{m.transactions_refund_success()}</p>
				<div class="space-y-1 text-xs">
					<div class="flex justify-between">
						<span class="text-muted-foreground">{m.transactions_refund_to_label()}</span>
						<span class="font-mono">{refundResult.to.slice(0, 12)}…</span>
					</div>
					<div class="flex justify-between">
						<span class="text-muted-foreground">{m.transactions_refund_tx_label()}</span>
						<span class="font-mono">{refundResult.tx_hash.slice(0, 12)}…</span>
					</div>
				</div>
			</div>
			<Dialog.Footer>
				<Button size="sm" class="text-xs" onclick={closeRefundDialog}>{m.transactions_close()}</Button>
			</Dialog.Footer>
		{:else}
			<div class="space-y-3 py-2">
				<p class="text-xs text-muted-foreground">
					{m.transactions_refund_message({ amount: refundPayment?.amount_crypto, symbol: refundPayment?.currency_crypto })}
				</p>
				{#if refundError}
					<p class="text-xs text-destructive">{refundError}</p>
				{/if}
			</div>
			<Dialog.Footer class="gap-2">
				<Button variant="outline" size="sm" class="text-xs" onclick={closeRefundDialog} disabled={refunding}>{m.transactions_cancel()}</Button>
				<Button size="sm" class="gap-1.5 text-xs" onclick={doRefund} disabled={refunding}>
					{#if refunding}
						<Spinner size={12} class="animate-spin" />
						{m.transactions_refund_sending()}
					{:else}
						<ArrowUUpLeft size={12} />
						{m.transactions_refund_confirm()}
					{/if}
				</Button>
			</Dialog.Footer>
		{/if}
	</Dialog.Content>
</Dialog.Root>
