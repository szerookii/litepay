<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { Motion } from 'svelte-motion';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import {
		Copy, Check, ArrowCounterClockwise,
		Warning, CheckCircle, ShieldCheck, ArrowRight,
	} from 'phosphor-svelte';
	import { qr } from '@svelte-put/qr/svg';
	import { m } from '$lib/paraglide/messages.js';
	import SEO from '$lib/components/seo.svelte';
	import { getLocale } from '$lib/paraglide/runtime.js';

	type Status = 'PENDING' | 'CONFIRMING' | 'PAID' | 'EXPIRED';

	interface Payment {
		id: string;
		wallet_address: string;
		sol_reference?: string;
		amount_crypto: number;
		currency_crypto_name: string;
		currency_crypto_symbol: string;
		amount_fiat: number;
		currency_fiat: string;
		status: Status;
		expires_at: string;
		last_transaction_hash?: string;
		confirmations: number | null;
		required_confirmations: number;
	}

	let payment = $state<Payment | null>(null);
	let loadError = $state('');
	let loading = $state(true);

	const paymentId = $derived(page.url.searchParams.get('id') ?? '');

	const coinColor: Record<string, string> = {
		LTC: '#BEBEBE', BCH: '#8DC351', ETH: '#627EEA', BTC: '#F7931A', SOL: '#9945FF',
	};
	const uriScheme: Record<string, string> = {
		LTC: 'litecoin', BCH: 'bitcoincash', ETH: 'ethereum', BTC: 'bitcoin', SOL: 'solana',
	};

	const color = $derived(payment ? (coinColor[payment.currency_crypto_symbol] ?? '#888') : '#888');
	const status = $derived<Status>(payment?.status ?? 'PENDING');
	const expiresAt = $derived(payment ? new Date(payment.expires_at) : new Date());
	const paidAt = $state<Date | null>(null);

	// For Solana, the QR encodes solana:<address>?amount=<sol>&reference=<ref>
	const qrData = $derived(() => {
		if (!payment) return '';
		const sym = payment.currency_crypto_symbol;
		const scheme = uriScheme[sym] ?? sym.toLowerCase();
		if (sym === 'SOL' && payment.sol_reference) {
			return `${scheme}:${payment.wallet_address}?amount=${payment.amount_crypto}&reference=${payment.sol_reference}`;
		}
		return `${scheme}:${payment.wallet_address}?amount=${payment.amount_crypto}`;
	});

	// Countdown timer
	let now = $state(Date.now());
	$effect(() => {
		const id = setInterval(() => (now = Date.now()), 1000);
		return () => clearInterval(id);
	});

	const remaining = $derived(Math.max(0, expiresAt.getTime() - now));
	const mins      = $derived(Math.floor(remaining / 60000));
	const secs      = $derived(Math.floor((remaining % 60000) / 1000));
	const timerStr  = $derived(`${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`);
	const progress  = $derived(remaining / (60 * 60 * 1000));
	const isUrgent  = $derived(remaining > 0 && remaining < 300_000);
	const isDead    = $derived(remaining === 0);

	// Copy address
	let copied = $state(false);
	function copyAddress() {
		if (!payment) return;
		navigator.clipboard.writeText(payment.wallet_address);
		copied = true;
		setTimeout(() => (copied = false), 2000);
	}

	async function fetchPayment(): Promise<void> {
		if (!paymentId) return;
		try {
			const res = await fetch(`/api/payment/${paymentId}`);
			if (!res.ok) {
				loadError = 'Payment not found';
				return;
			}
			const data: Payment = await res.json();
			payment = data;
			loading = false;
		} catch {
			loadError = 'Network error';
			loading = false;
		}
	}

	onMount(() => {
		if (!paymentId) {
			loadError = 'No payment ID in URL';
			loading = false;
			return;
		}
		fetchPayment();
		// Poll every 5s until terminal state
		const interval = setInterval(() => {
			if (payment?.status === 'PAID' || payment?.status === 'EXPIRED') {
				clearInterval(interval);
				return;
			}
			fetchPayment();
		}, 5000);
		return () => clearInterval(interval);
	});
</script>

<SEO
	config={{
		title: m.seo_pay_title(),
		description: m.seo_pay_desc(),
		ogType: 'website'
	}}
	url={page.url.pathname}
	locale={getLocale()}
/>

<!-- ambient glow -->
<div class="pointer-events-none fixed inset-0 -z-10 overflow-hidden" aria-hidden="true">
	<div
		class="absolute left-1/2 top-0 h-[500px] w-[600px] -translate-x-1/2 rounded-full blur-[120px]
			{status === 'PAID' ? 'opacity-[0.08]' : 'opacity-[0.05]'}"
		style="background:{status === 'PAID' ? '#4ade80' : status === 'EXPIRED' ? '#ef4444' : color}"
	></div>
</div>

<div class="flex min-h-svh flex-col bg-background">
	<!-- NAV -->
	<nav class="border-border/40 sticky top-0 z-50 border-b bg-background/80 backdrop-blur-md">
		<div class="mx-auto flex max-w-lg items-center justify-between px-4 py-3 sm:px-6">
			<span class="text-sm font-medium text-foreground">
				lite<span class="text-muted-foreground">pay</span>
			</span>
			{#if payment}
				<span class="font-mono text-xs text-muted-foreground">#{payment.id.slice(0, 8)}</span>
			{/if}
		</div>
	</nav>

	<main class="flex flex-1 flex-col items-center px-0 py-0 sm:px-4 sm:py-10">
		<div class="w-full max-w-lg">

			{#if loading}
				<div class="flex items-center justify-center py-24">
					<div class="space-y-2 text-center">
						<div class="mx-auto size-8 animate-spin rounded-full border-2 border-border border-t-foreground"></div>
						<p class="text-xs text-muted-foreground">Loading payment…</p>
					</div>
				</div>
			{:else if loadError || !payment}
				<div class="ring-border mx-4 mt-10 p-6 ring-1 sm:mx-0">
					<p class="text-sm font-medium text-destructive">{loadError || 'Payment not found'}</p>
					<p class="mt-1 text-xs text-muted-foreground">Check the URL and try again.</p>
				</div>
			{:else}
				<!-- CARD -->
				<Motion
					initial={{ opacity: 0, y: 14 }}
					animate={{ opacity: 1, y: 0 }}
					transition={{ duration: 0.35 }}
					let:motion
				>
					<div
						use:motion
						class="ring-border/50 overflow-hidden bg-card sm:ring-1
							{status === 'PAID'    ? 'sm:ring-green-500/25' : ''}
							{status === 'EXPIRED' ? 'sm:ring-destructive/25' : ''}"
					>
						<!-- HEADER -->
						<div class="px-5 py-5 sm:px-6">
							<div class="mb-5 flex items-center justify-between">
								<div class="flex items-center gap-2">
									<span
										class="inline-flex size-5 items-center justify-center rounded-full text-[9px] font-bold"
										style="background:{color}18;color:{color};outline:1px solid {color}40"
									>{payment.currency_crypto_symbol[0]}</span>
									<span class="text-xs font-medium">{payment.currency_crypto_name}</span>
									<Badge variant="outline" class="hidden sm:inline-flex">{payment.currency_crypto_symbol}</Badge>
								</div>

								{#if status === 'PENDING' || status === 'CONFIRMING'}
									<div class="flex items-center gap-1.5 text-xs text-muted-foreground">
										<span class="relative flex size-2">
											<span class="absolute inline-flex size-full animate-ping rounded-full opacity-60" style="background:{color}"></span>
											<span class="relative inline-flex size-2 rounded-full" style="background:{color}"></span>
										</span>
										{status === 'CONFIRMING' ? `Confirming… ${payment.confirmations ?? 0}/${payment.required_confirmations}` : m.pay_awaiting()}
									</div>
								{:else if status === 'PAID'}
									<div class="flex items-center gap-1.5 text-xs text-green-400">
										<CheckCircle size={12} weight="fill" />
										{m.pay_confirmed_status()}
									</div>
								{:else}
									<div class="flex items-center gap-1.5 text-xs text-destructive">
										<Warning size={12} weight="fill" />
										{m.pay_expired_status()}
									</div>
								{/if}
							</div>

							<!-- amount -->
							<div class="space-y-1">
								<p class="text-xs text-muted-foreground">
									{status === 'PAID' ? m.pay_amount_received() : m.pay_amount_due()}
								</p>
								<div class="flex items-baseline gap-2">
									<span
										class="text-[2rem] font-semibold leading-none tracking-tight sm:text-[2.4rem]"
										style="color:{status === 'PAID' ? '#4ade80' : status === 'EXPIRED' ? 'var(--muted-foreground)' : color}"
									>{payment.amount_crypto}</span>
									<span class="text-base text-muted-foreground">{payment.currency_crypto_symbol}</span>
								</div>
								<p class="text-xs text-muted-foreground">
									{m.pay_approx({ currency: payment.currency_fiat, amount: payment.amount_fiat.toFixed(2) })}
								</p>
							</div>
						</div>

						<Separator />

						<!-- ── PENDING / CONFIRMING body ── -->
						{#if status === 'PENDING' || status === 'CONFIRMING'}

							<!-- QR -->
							<div class="flex flex-col items-center gap-4 px-5 py-7 sm:px-6">
								<p class="text-xs text-muted-foreground">{m.pay_scan_wallet()}</p>
								<div class="relative">
									<div class="absolute inset-0 rounded-full opacity-20 blur-2xl" style="background:{color}"></div>
									<div class="ring-border relative p-3 ring-1" style="background:#ffffff;color:#000000">
										<svg
											use:qr={{
												data: qrData(),
												moduleFill: '#000000',
												anchorOuterFill: '#000000',
												anchorInnerFill: '#000000',
											}}
											style="width:clamp(160px,40vw,196px);height:clamp(160px,40vw,196px);display:block"
										></svg>
									</div>
								</div>
								<p class="max-w-[200px] text-center text-xs leading-relaxed text-muted-foreground">
									{m.pay_send_exactly()} <span class="font-medium text-foreground">{payment.amount_crypto} {payment.currency_crypto_symbol}</span> {m.pay_to_address_below()}
								</p>
							</div>

							<Separator />

							<!-- address -->
							<div class="px-5 py-4 sm:px-6">
								<p class="mb-2 text-xs text-muted-foreground">{m.pay_wallet_address()}</p>
								<button
									onclick={copyAddress}
									class="ring-border group flex w-full items-center justify-between gap-3 bg-muted/30 px-3 py-3 ring-1 transition-colors hover:bg-muted/50 active:scale-[0.99]"
									aria-label="Copy address"
								>
									<span class="truncate font-mono text-xs text-foreground">{payment.wallet_address}</span>
									<span class="shrink-0 text-muted-foreground transition-colors group-hover:text-foreground">
										{#if copied}
											<Check size={14} weight="bold" color="#4ade80" />
										{:else}
											<Copy size={14} />
										{/if}
									</span>
								</button>
								{#if copied}
									<p class="mt-1.5 text-xs text-green-400">{m.pay_address_copied()}</p>
								{/if}
							</div>

							<Separator />

							<!-- timer -->
							<div class="px-5 py-4 sm:px-6">
								<div class="mb-2 flex items-center justify-between">
									<p class="text-xs text-muted-foreground">{m.pay_time_remaining()}</p>
									<span class="font-mono text-sm tabular-nums {isUrgent ? 'animate-pulse text-destructive' : isDead ? 'text-destructive/60' : 'text-foreground'}">
										{isDead ? '00:00' : timerStr}
									</span>
								</div>
								<div class="bg-muted h-1 w-full overflow-hidden">
									<div
										class="h-full transition-[width] duration-1000 ease-linear"
										style="width:{(progress * 100).toFixed(3)}%;background:{isUrgent ? 'oklch(0.704 0.191 22.216)' : color}"
									></div>
								</div>
								<p class="mt-2 text-xs text-muted-foreground">
									{isDead
										? m.pay_payment_expired_timer()
										: m.pay_expires_at({ time: expiresAt.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) })}
								</p>
							</div>

							<Separator />

							<div class="flex items-center justify-between bg-muted/10 px-5 py-3 sm:px-6">
								<p class="text-xs text-muted-foreground">
									{m.pay_need_help()} <a href="/" class="text-foreground underline-offset-2 hover:underline">{m.pay_contact_us()}</a>
								</p>
								<Button variant="ghost" size="sm" class="gap-1.5 active:scale-95" onclick={fetchPayment}>
									<ArrowCounterClockwise size={12} />
									{m.pay_refresh()}
								</Button>
							</div>

						<!-- ── PAID body ── -->
						{:else if status === 'PAID'}
							<Motion
								initial={{ opacity: 0, scale: 0.96 }}
								animate={{ opacity: 1, scale: 1 }}
								transition={{ duration: 0.3 }}
								let:motion
							>
								<div use:motion class="flex flex-col items-center gap-6 px-5 py-10 sm:px-6">
									<div class="relative flex items-center justify-center">
										<div class="absolute size-20 animate-ping rounded-full bg-green-400/10"></div>
										<div class="relative flex size-16 items-center justify-center rounded-full bg-green-400/10 ring-1 ring-green-400/30">
											<CheckCircle size={32} weight="fill" color="#4ade80" />
										</div>
									</div>

									<div class="space-y-1 text-center">
										<p class="text-sm font-medium text-foreground">{m.pay_payment_received()}</p>
										<p class="text-xs text-muted-foreground">
											{m.pay_confirmed_onchain({ amount: payment.amount_crypto, symbol: payment.currency_crypto_symbol })}
										</p>
									</div>

									<div class="ring-border w-full space-y-3 bg-muted/20 px-4 py-3 ring-1">
										{#if payment.last_transaction_hash}
											<div class="flex items-center justify-between">
												<span class="text-xs text-muted-foreground">{m.pay_transaction()}</span>
												<span class="font-mono text-xs text-foreground">{payment.last_transaction_hash.slice(0, 8)}…{payment.last_transaction_hash.slice(-6)}</span>
											</div>
											<Separator />
										{/if}
										<div class="flex items-center justify-between">
											<span class="text-xs text-muted-foreground">{m.pay_network()}</span>
											<span class="text-xs text-foreground">{payment.currency_crypto_name}</span>
										</div>
									</div>

									<Button href="/" class="w-full gap-1.5">
										{m.pay_back_to_store()}
										<ArrowRight size={12} />
									</Button>
								</div>
							</Motion>

						<!-- ── EXPIRED body ── -->
						{:else}
							<Motion
								initial={{ opacity: 0, scale: 0.96 }}
								animate={{ opacity: 1, scale: 1 }}
								transition={{ duration: 0.3 }}
								let:motion
							>
								<div use:motion class="flex flex-col items-center gap-6 px-5 py-10 sm:px-6">
									<div class="flex size-16 items-center justify-center rounded-full bg-destructive/10 ring-1 ring-destructive/30">
										<Warning size={32} weight="fill" color="oklch(0.704 0.191 22.216)" />
									</div>

									<div class="space-y-1 text-center">
										<p class="text-sm font-medium text-foreground">{m.pay_payment_expired_title()}</p>
										<p class="text-xs leading-relaxed text-muted-foreground">{m.pay_payment_expired_desc()}</p>
									</div>

									<div class="ring-border w-full space-y-3 bg-muted/20 px-4 py-3 ring-1">
										<div class="flex items-center justify-between">
											<span class="text-xs text-muted-foreground">{m.pay_amount()}</span>
											<span class="text-xs text-muted-foreground line-through">
												{payment.amount_crypto} {payment.currency_crypto_symbol}
											</span>
										</div>
										<Separator />
										<div class="flex items-center justify-between">
											<span class="text-xs text-muted-foreground">{m.pay_expired_at_label()}</span>
											<span class="text-xs text-foreground">
												{expiresAt.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
											</span>
										</div>
									</div>

									<div class="flex w-full flex-col gap-2">
										<Button href="/" class="w-full gap-1.5">
											{m.pay_create_new_payment()}
											<ArrowRight size={12} />
										</Button>
										<Button variant="ghost" href="/" class="w-full text-muted-foreground">
											{m.pay_back_to_store()}
										</Button>
									</div>
								</div>
							</Motion>
						{/if}
					</div>
				</Motion>

				<div class="flex items-center justify-center gap-1.5 px-4 py-5 text-xs text-muted-foreground/50">
					<ShieldCheck size={12} />
					{m.pay_trust_line()}
				</div>
			{/if}
		</div>
	</main>
</div>
