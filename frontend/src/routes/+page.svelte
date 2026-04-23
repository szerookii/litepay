<script lang="ts">
	import { Motion } from 'svelte-motion';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';

	let allowRegister = $state<boolean>(true);

	$effect(() => {
		fetch('/api/config')
			.then((r) => r.json())
			.then((d) => { allowRegister = d.allow_register ?? true; })
			.catch(() => {});
	});

	const features = [
		{
			icon: '⬡',
			title: 'Self-Hosted',
			desc: 'Run on your own infrastructure. Your keys, your coins, your rules. No third-party holds your funds.',
		},
		{
			icon: '#',
			title: 'Open Source',
			desc: 'Full source code on GitHub. Inspect every line, fork it, extend it. No black boxes.',
		},
		{
			icon: '⚡',
			title: 'API-First',
			desc: 'Simple JSON API. Create a payment in one HTTP request. Integrate with any stack in minutes.',
		},
		{
			icon: '↻',
			title: 'Auto-Verify',
			desc: 'Background cron checks on-chain confirmations automatically. No webhooks to configure manually.',
		},
	];

	const coins = [
		{
			symbol: 'BTC',
			name: 'Bitcoin',
			color: 'text-[#F7931A]',
			border: 'border-[#F7931A]/20',
			bg: 'bg-[#F7931A]/5',
			desc: 'Non-custodial HD wallet. Unique address per payment via xpub derivation.',
		},
		{
			symbol: 'LTC',
			name: 'Litecoin',
			color: 'text-[#BEBEBE]',
			border: 'border-[#BEBEBE]/20',
			bg: 'bg-[#BEBEBE]/5',
			desc: 'Fast, low-fee payments. Battle-tested since 2011. BIP84 xpub support.',
		},
		{
			symbol: 'SOL',
			name: 'Solana',
			color: 'text-[#9945FF]',
			border: 'border-[#9945FF]/20',
			bg: 'bg-[#9945FF]/5',
			desc: 'Solana Pay protocol. Unique reference key per payment, no xpub needed.',
		},
	];

	const steps = [
		{ n: '01', title: 'Deploy', desc: 'Clone, add your xpub/wallet address in the dashboard, run the binary.' },
		{ n: '02', title: 'Create Payment', desc: 'POST /api/payment with symbol, amount, currency. Get a deposit address back.' },
		{ n: '03', title: 'Get Paid', desc: 'LitePay polls the chain. Status updates from pending → confirmed automatically.' },
	];

	const apiExample = `POST /payments
Content-Type: application/json
X-Api-Key: YOUR_API_KEY

{
  "symbol": "LTC",
  "amount": 49.99,
  "currency": "USD"
}

// Response
{
  "id": "0196f2a1-...",
  "wallet_address": "LbYkLqV...",
  "amount_crypto": 0.82341,
  "currency_crypto": "LTC",
  "amount_fiat": 49.99,
  "currency_fiat": "USD",
  "status": "pending",
  "expires_at": "2025-01-01T13:00:00Z"
}`;

	function fade(delay = 0) {
		return {
			initial: { opacity: 0, y: 24 },
			animate: { opacity: 1, y: 0 },
			transition: { duration: 0.5, delay },
		};
	}
</script>

<!-- NAV -->
<nav class="border-border/50 sticky top-0 z-50 border-b bg-background/80 backdrop-blur-sm">
	<div class="mx-auto flex max-w-6xl items-center justify-between px-6 py-3">
		<div class="flex items-center gap-2">
			<span class="text-sm font-medium tracking-tight text-foreground">lite<span class="text-muted-foreground">pay</span></span>
			<Badge variant="outline">v0.1</Badge>
		</div>
		<div class="flex items-center gap-2">
			<Button variant="ghost" size="sm" href="https://github.com/szerookii/litepay">GitHub</Button>
			<Button variant="ghost" size="sm" href="#api">Docs</Button>
			{#if allowRegister}
				<Button variant="ghost" size="sm" href="/auth/register">Register</Button>
			{/if}
			<Button variant="outline" size="sm" href="/auth/login">Login</Button>
		</div>
	</div>
</nav>

<!-- HERO -->
<section class="relative overflow-hidden px-6 pb-24 pt-24">
	<div class="pointer-events-none absolute inset-0 -z-10">
		<div class="absolute left-1/2 top-0 h-[500px] w-[800px] -translate-x-1/2 rounded-full bg-primary/5 blur-3xl"></div>
	</div>

	<div class="mx-auto max-w-6xl">
		<div class="grid items-center gap-16 lg:grid-cols-2">
			<div class="space-y-6">
				<Motion {...fade(0)} let:motion>
					<div use:motion>
						<Badge variant="secondary" class="mb-2">Open Source · Self-Hosted · MIT</Badge>
					</div>
				</Motion>

				<Motion {...fade(0.1)} let:motion>
					<h1 use:motion class="text-4xl font-medium leading-tight tracking-tight text-foreground lg:text-5xl">
						Accept crypto payments.<br />
						<span class="text-muted-foreground">No middlemen.</span>
					</h1>
				</Motion>

				<Motion {...fade(0.2)} let:motion>
					<p use:motion class="text-sm leading-relaxed text-muted-foreground">
						LitePay is an open-source crypto payment processor you run on your own server.
						One API call to create a payment, on-chain verification handled automatically.
						BTC, LTC, and SOL out of the box.
					</p>
				</Motion>

				<Motion {...fade(0.3)} let:motion>
					<div use:motion class="flex flex-wrap gap-2">
						<Button href="https://github.com/szerookii/litepay" size="lg">
							View on GitHub
						</Button>
						<Button variant="outline" href="#how-it-works" size="lg">
							How it works
						</Button>
					</div>
				</Motion>
			</div>

			<!-- Terminal -->
			<Motion {...fade(0.2)} let:motion>
				<div use:motion class="ring-border overflow-hidden rounded-none ring-1">
					<div class="border-border flex items-center gap-1.5 border-b bg-muted/50 px-4 py-2">
						<span class="size-2.5 rounded-full bg-destructive/60"></span>
						<span class="size-2.5 rounded-full bg-yellow-400/60"></span>
						<span class="size-2.5 rounded-full bg-green-400/60"></span>
						<span class="ml-2 text-xs text-muted-foreground">litepay — api</span>
					</div>
					<pre class="overflow-x-auto bg-card p-5 text-xs leading-relaxed"><code><span class="text-muted-foreground"># create a payment</span>
<span class="text-green-400">$</span> curl -X POST https://pay.example.com/payments \
    -H <span class="text-yellow-400">"X-Api-Key: sk_live_..."</span> \
    -d <span class="text-yellow-400">'&#123;"symbol":"LTC","amount":49.99,"currency":"USD"&#125;'</span>

<span class="text-muted-foreground">&#123;</span>
  <span class="text-blue-400">"wallet_address"</span>: <span class="text-yellow-400">"LbYkLqV8Ur3..."</span>,
  <span class="text-blue-400">"amount_crypto"</span>: <span class="text-green-400">0.82341</span>,
  <span class="text-blue-400">"currency_crypto"</span>: <span class="text-yellow-400">"LTC"</span>,
  <span class="text-blue-400">"status"</span>: <span class="text-yellow-400">"pending"</span>,
  <span class="text-blue-400">"expires_at"</span>: <span class="text-yellow-400">"2025-01-01T13:00:00Z"</span>
<span class="text-muted-foreground">&#125;</span></code></pre>
				</div>
			</Motion>
		</div>
	</div>
</section>

<Separator />

<!-- FEATURES -->
<section class="px-6 py-20">
	<div class="mx-auto max-w-6xl">
		<Motion {...fade(0)} let:motion>
			<div use:motion class="mb-12 space-y-2">
				<p class="text-xs uppercase tracking-widest text-muted-foreground">Why LitePay</p>
				<h2 class="text-2xl font-medium text-foreground">Built for developers who own their stack.</h2>
			</div>
		</Motion>

		<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
			{#each features as f, i}
				<Motion {...fade(i * 0.08)} let:motion>
					<Card>
						<CardHeader>
							<span class="text-xl text-muted-foreground">{f.icon}</span>
							<CardTitle>{f.title}</CardTitle>
						</CardHeader>
						<CardContent>
							<p class="text-muted-foreground">{f.desc}</p>
						</CardContent>
					</Card>
				</Motion>
			{/each}
		</div>
	</div>
</section>

<Separator />

<!-- COINS -->
<section class="px-6 py-20">
	<div class="mx-auto max-w-6xl">
		<Motion {...fade(0)} let:motion>
			<div use:motion class="mb-12 space-y-2">
				<p class="text-xs uppercase tracking-widest text-muted-foreground">Supported coins</p>
				<h2 class="text-2xl font-medium text-foreground">Three chains, one API. Non-custodial.</h2>
			</div>
		</Motion>

		<div class="grid gap-4 sm:grid-cols-3">
			{#each coins as coin, i}
				<Motion {...fade(i * 0.1)} let:motion>
					<div use:motion class="ring-border {coin.bg} rounded-none p-6 ring-1 transition-all hover:ring-2">
						<div class="mb-4 flex items-center gap-3">
							<span class="text-2xl font-bold {coin.color}">{coin.symbol}</span>
							<span class="text-sm text-muted-foreground">{coin.name}</span>
						</div>
						<p class="text-xs text-muted-foreground">{coin.desc}</p>
					</div>
				</Motion>
			{/each}
		</div>
	</div>
</section>

<Separator />

<!-- HOW IT WORKS -->
<section id="how-it-works" class="px-6 py-20">
	<div class="mx-auto max-w-6xl">
		<Motion {...fade(0)} let:motion>
			<div use:motion class="mb-12 space-y-2">
				<p class="text-xs uppercase tracking-widest text-muted-foreground">Integration</p>
				<h2 class="text-2xl font-medium text-foreground">Up and running in minutes.</h2>
			</div>
		</Motion>

		<div class="grid gap-8 sm:grid-cols-3">
			{#each steps as step, i}
				<Motion {...fade(i * 0.1)} let:motion>
					<div use:motion class="space-y-3">
						<span class="text-4xl font-medium text-border">{step.n}</span>
						<h3 class="text-sm font-medium text-foreground">{step.title}</h3>
						<p class="text-xs leading-relaxed text-muted-foreground">{step.desc}</p>
					</div>
				</Motion>
			{/each}
		</div>
	</div>
</section>

<Separator />

<!-- API REFERENCE -->
<section id="api" class="px-6 py-20">
	<div class="mx-auto max-w-6xl">
		<Motion {...fade(0)} let:motion>
			<div use:motion class="mb-12 space-y-2">
				<p class="text-xs uppercase tracking-widest text-muted-foreground">API</p>
				<h2 class="text-2xl font-medium text-foreground">One endpoint. That's it.</h2>
			</div>
		</Motion>

		<Motion {...fade(0.1)} let:motion>
			<div use:motion class="ring-border overflow-hidden rounded-none ring-1">
				<div class="border-border flex items-center justify-between border-b bg-muted/50 px-4 py-2">
					<div class="flex items-center gap-2">
						<Badge variant="secondary">POST</Badge>
						<span class="text-xs text-muted-foreground">/payments</span>
					</div>
					<span class="text-xs text-muted-foreground">application/json</span>
				</div>
				<pre class="overflow-x-auto bg-card p-6 text-xs leading-relaxed text-foreground"><code>{apiExample}</code></pre>
			</div>
		</Motion>
	</div>
</section>

<Separator />

<!-- CTA -->
<section class="px-6 py-24 text-center">
	<div class="mx-auto max-w-xl space-y-6">
		<Motion {...fade(0)} let:motion>
			<h2 use:motion class="text-3xl font-medium text-foreground">
				Own your payment stack.
			</h2>
		</Motion>
		<Motion {...fade(0.1)} let:motion>
			<p use:motion class="text-sm text-muted-foreground">
				No fees. No vendor lock-in. No KYC. Deploy LitePay and start accepting crypto payments today.
			</p>
		</Motion>
		<Motion {...fade(0.2)} let:motion>
			<div use:motion class="flex justify-center gap-2">
				<Button href="https://github.com/szerookii/litepay" size="lg">
					Get Started →
				</Button>
				<Button variant="outline" size="lg" href="https://github.com/szerookii/litepay">
					Star on GitHub
				</Button>
			</div>
		</Motion>
	</div>
</section>

<!-- FOOTER -->
<footer class="border-border border-t px-6 py-8">
	<div class="mx-auto flex max-w-6xl items-center justify-between">
		<span class="text-xs text-muted-foreground">litepay · MIT License</span>
		<div class="flex items-center gap-4">
			<a href="https://github.com/szerookii/litepay" class="text-xs text-muted-foreground hover:text-foreground transition-colors">GitHub</a>
		</div>
	</div>
</footer>
