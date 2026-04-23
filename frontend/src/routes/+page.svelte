<script lang="ts">
	import { Motion } from 'svelte-motion';
	import { page } from '$app/state';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import * as m from '$lib/paraglide/messages.js';
	import { locales, getLocale, localizeHref, setLocale } from '$lib/paraglide/runtime.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { TranslateIcon } from 'phosphor-svelte';
	import SEO from '$lib/components/seo.svelte';

	let allowRegister = $state<boolean>(true);
	let version = $state<string>('v0.1');

	$effect(() => {
		fetch('/api/config')
			.then((r) => r.json())
			.then((d) => {
				allowRegister = d.allow_register ?? true;
				version = d.version ? `v${d.version}` : 'v0.1';
			})
			.catch(() => {});
	});

	const features = [
		{
			icon: '⬡',
			title: m.landing_feature_1_title(),
			desc: m.landing_feature_1_desc(),
		},
		{
			icon: '#',
			title: m.landing_feature_2_title(),
			desc: m.landing_feature_2_desc(),
		},
		{
			icon: '⚡',
			title: m.landing_feature_3_title(),
			desc: m.landing_feature_3_desc(),
		},
		{
			icon: '↻',
			title: m.landing_feature_4_title(),
			desc: m.landing_feature_4_desc(),
		},
	];

	const coins = [
		{
			symbol: 'BTC',
			name: 'Bitcoin',
			color: 'text-[#F7931A]',
			border: 'border-[#F7931A]/20',
			bg: 'bg-[#F7931A]/5',
			desc: m.landing_coin_btc_desc(),
		},
		{
			symbol: 'LTC',
			name: 'Litecoin',
			color: 'text-[#BEBEBE]',
			border: 'border-[#BEBEBE]/20',
			bg: 'bg-[#BEBEBE]/5',
			desc: m.landing_coin_ltc_desc(),
		},
		{
			symbol: 'SOL',
			name: 'Solana',
			color: 'text-[#9945FF]',
			border: 'border-[#9945FF]/20',
			bg: 'bg-[#9945FF]/5',
			desc: m.landing_coin_sol_desc(),
		},
	];

	const steps = [
		{ n: '01', title: m.landing_step_1_title(), desc: m.landing_step_1_desc() },
		{ n: '02', title: m.landing_step_2_title(), desc: m.landing_step_2_desc() },
		{ n: '03', title: m.landing_step_3_title(), desc: m.landing_step_3_desc() },
	];

	const apiExample = `POST /payments
Content-Type: application/json
Authorization: Bearer YOUR_API_KEY

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

<SEO
	config={{
		title: m.seo_home_title(),
		description: m.seo_home_desc(),
		keywords: m.seo_home_keywords(),
		ogType: 'website'
	}}
	url={page.url.pathname}
	locale={getLocale()}
/>

<!-- NAV -->
<nav class="border-border/50 sticky top-0 z-50 border-b bg-background/80 backdrop-blur-sm">
	<div class="mx-auto flex max-w-6xl items-center justify-between px-3 py-3 sm:px-6">
		<div class="flex items-center gap-1 sm:gap-2 min-w-0">
			<span class="text-xs sm:text-sm font-medium tracking-tight text-foreground whitespace-nowrap">lite<span class="text-muted-foreground">pay</span></span>
			<Badge variant="outline" class="hidden sm:flex text-xs">{version}</Badge>
		</div>
		<div class="flex items-center gap-1 sm:gap-2">
			<DropdownMenu.Root>
				<DropdownMenu.Trigger>
					{#snippet child({ props })}
						<Button variant="ghost" size="sm" {...props} class="gap-1 sm:gap-2 px-2 sm:px-3">
							<TranslateIcon size={16} />
							<span class="uppercase hidden sm:inline">{getLocale()}</span>
						</Button>
					{/snippet}
				</DropdownMenu.Trigger>
				<DropdownMenu.Content align="end">
					{#each locales as locale}
						<DropdownMenu.Item>
							<button
								onclick={(e) => {
									e.preventDefault();
									setLocale(locale);
								}}
								class="flex w-full items-center justify-between gap-2"
							>
								<span class="capitalize">
									{locale === 'en' ? 'English' : locale === 'fr' ? 'Français' : locale}
								</span>
								{#if locale === getLocale()}
									<div class="size-1.5 rounded-full bg-primary"></div>
								{/if}
							</button>
						</DropdownMenu.Item>
					{/each}
				</DropdownMenu.Content>
			</DropdownMenu.Root>

			<Separator orientation="vertical" class="mx-0.5 sm:mx-1 h-4 hidden sm:block" />

			<Button variant="ghost" size="sm" href="https://github.com/szerookii/litepay" class="px-2 sm:px-3 hidden sm:inline-flex">{m.landing_nav_github()}</Button>
			<Button variant="ghost" size="sm" href="#api" class="px-2 sm:px-3 hidden sm:inline-flex">{m.landing_nav_docs()}</Button>
			{#if allowRegister}
				<Button variant="ghost" size="sm" href="/auth/register" class="px-2 sm:px-3 hidden sm:inline-flex">{m.landing_nav_register()}</Button>
			{/if}
			<Button variant="outline" size="sm" href="/auth/login" class="px-2 sm:px-3 text-xs sm:text-sm">{m.landing_nav_login()}</Button>
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
						<Badge variant="secondary" class="mb-2">{m.landing_hero_badge()}</Badge>
					</div>
				</Motion>

				<Motion {...fade(0.1)} let:motion>
					<h1 use:motion class="text-4xl font-medium leading-tight tracking-tight text-foreground lg:text-5xl">
						{@html m.landing_hero_title()}
					</h1>
				</Motion>

				<Motion {...fade(0.2)} let:motion>
					<p use:motion class="text-sm leading-relaxed text-muted-foreground">
						{m.landing_hero_subtitle()}
					</p>
				</Motion>

				<Motion {...fade(0.3)} let:motion>
					<div use:motion class="flex flex-wrap gap-2">
						<Button href="https://github.com/szerookii/litepay" size="lg">
							{m.landing_hero_cta_github()}
						</Button>
						<Button variant="outline" href="#how-it-works" size="lg">
							{m.landing_hero_cta_how()}
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
					<pre class="overflow-x-auto bg-card p-5 text-xs leading-relaxed"><code><span class="text-muted-foreground">{m.landing_hero_terminal_comment()}</span>
<span class="text-green-400">$</span> curl -X POST https://pay.example.com/payments \
    -H <span class="text-yellow-400">"Authorization: Bearer sk_live_..."</span> \
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
				<p class="text-xs uppercase tracking-widest text-muted-foreground">{m.landing_features_label()}</p>
				<h2 class="text-2xl font-medium text-foreground">{m.landing_features_title()}</h2>
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
				<p class="text-xs uppercase tracking-widest text-muted-foreground">{m.landing_coins_label()}</p>
				<h2 class="text-2xl font-medium text-foreground">{m.landing_coins_title()}</h2>
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
				<p class="text-xs uppercase tracking-widest text-muted-foreground">{m.landing_steps_label()}</p>
				<h2 class="text-2xl font-medium text-foreground">{m.landing_steps_title()}</h2>
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
				<p class="text-xs uppercase tracking-widest text-muted-foreground">{m.landing_api_label()}</p>
				<h2 class="text-2xl font-medium text-foreground">{m.landing_api_title()}</h2>
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
				{m.landing_cta_title()}
			</h2>
		</Motion>
		<Motion {...fade(0.1)} let:motion>
			<p use:motion class="text-sm text-muted-foreground">
				{m.landing_cta_desc()}
			</p>
		</Motion>
		<Motion {...fade(0.2)} let:motion>
			<div use:motion class="flex justify-center gap-2">
				<Button href="https://github.com/szerookii/litepay" size="lg">
					{m.landing_cta_start()}
				</Button>
				<Button variant="outline" size="lg" href="https://github.com/szerookii/litepay">
					{m.landing_cta_star()}
				</Button>
			</div>
		</Motion>
	</div>
</section>

<!-- FOOTER -->
<footer class="border-border border-t px-6 py-8">
	<div class="mx-auto flex max-w-6xl items-center justify-between">
		<span class="text-xs text-muted-foreground">{m.landing_footer_license()}</span>
		<div class="flex items-center gap-4">
			<a href="https://github.com/szerookii/litepay" class="text-xs text-muted-foreground hover:text-foreground transition-colors">{m.landing_nav_github()}</a>
		</div>
	</div>
</footer>
