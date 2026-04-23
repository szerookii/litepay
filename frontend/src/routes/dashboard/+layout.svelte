<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import { House, Wallet, Key, SignOut, List, X, Receipt, ArrowUpRight } from 'phosphor-svelte';
	import * as m from '$lib/paraglide/messages.js';
	import SEO from '$lib/components/seo.svelte';
	import { getLocale } from '$lib/paraglide/runtime.js';
	import { getSEOConfig } from '$lib/seo/config.js';

	let { children } = $props();

	let mobileOpen = $state(false);
	let version = $state('v0.1');

	$effect(() => {
		fetch('/api/config')
			.then((r) => r.json())
			.then((d) => {
				version = d.version ? `v${d.version}` : 'v0.1';
			})
			.catch(() => {});
	});

	function logout() {
		localStorage.removeItem('token');
		goto('/auth/login');
	}

	const nav = [
		{ href: '/dashboard', label: m.dash_nav_dashboard(), icon: House },
		{ href: '/dashboard/transactions', label: m.dash_nav_transactions(), icon: Receipt },
		{ href: '/dashboard/cashout', label: 'Cashout', icon: ArrowUpRight },
		{ href: '/dashboard/wallets', label: m.dash_nav_wallets(), icon: Wallet },
		{ href: '/dashboard/api-keys', label: m.dash_nav_api_keys(), icon: Key }
	];

	const currentPath = $derived(page.url.pathname);
</script>

<SEO
	config={getSEOConfig(currentPath)}
	url={page.url.pathname}
	locale={getLocale()}
/>

<div class="flex min-h-svh bg-background">
	<!-- Sidebar (desktop) -->
	<aside class="border-border hidden w-56 shrink-0 flex-col border-r lg:flex">
		<div class="flex items-center gap-2 px-5 py-4">
			<span class="text-sm font-medium tracking-tight text-foreground"
				>lite<span class="text-muted-foreground">pay</span></span
			>
			<Badge variant="outline" class="text-[10px]">{version}</Badge>
		</div>
		<Separator />
		<nav class="flex flex-1 flex-col gap-0.5 p-3">
			{#each nav as item}
				{@const active = currentPath === item.href}
				<a
					href={item.href}
					class="flex items-center gap-2.5 px-3 py-2 text-xs transition-colors
						{active
						? 'bg-muted text-foreground'
						: 'text-muted-foreground hover:bg-muted/50 hover:text-foreground'}"
				>
					<item.icon size={14} weight={active ? 'fill' : 'regular'} />
					{item.label}
				</a>
			{/each}
		</nav>
		<Separator />
		<div class="p-3">
			<button
				onclick={logout}
				class="flex w-full items-center gap-2.5 px-3 py-2 text-xs text-muted-foreground transition-colors hover:bg-muted/50 hover:text-destructive"
			>
				<SignOut size={14} />
				{m.dash_nav_logout()}
			</button>
		</div>
	</aside>

	<!-- Mobile overlay -->
	{#if mobileOpen}
		<div
			class="fixed inset-0 z-40 bg-background/80 backdrop-blur-sm lg:hidden"
			onclick={() => (mobileOpen = false)}
			aria-hidden="true"
		></div>
		<aside
			class="border-border fixed inset-y-0 left-0 z-50 flex w-56 flex-col border-r bg-background lg:hidden"
		>
			<div class="flex items-center justify-between px-5 py-4">
				<span class="text-sm font-medium text-foreground"
					>lite<span class="text-muted-foreground">pay</span></span
				>
				<button onclick={() => (mobileOpen = false)} class="text-muted-foreground hover:text-foreground">
					<X size={16} />
				</button>
			</div>
			<Separator />
			<nav class="flex flex-1 flex-col gap-0.5 p-3">
				{#each nav as item}
					{@const active = currentPath === item.href}
					<a
						href={item.href}
						onclick={() => (mobileOpen = false)}
						class="flex items-center gap-2.5 px-3 py-2 text-xs transition-colors
							{active
							? 'bg-muted text-foreground'
							: 'text-muted-foreground hover:bg-muted/50 hover:text-foreground'}"
					>
						<item.icon size={14} weight={active ? 'fill' : 'regular'} />
						{item.label}
					</a>
				{/each}
			</nav>
			<Separator />
			<div class="p-3">
				<button
					onclick={logout}
					class="flex w-full items-center gap-2.5 px-3 py-2 text-xs text-muted-foreground transition-colors hover:bg-muted/50 hover:text-destructive"
				>
					<SignOut size={14} />
					{m.dash_nav_logout()}
				</button>
			</div>
		</aside>
	{/if}

	<!-- Main content -->
	<div class="flex min-w-0 flex-1 flex-col">
		<!-- Mobile top bar -->
		<header class="border-border flex items-center gap-3 border-b px-4 py-3 lg:hidden">
			<button onclick={() => (mobileOpen = true)} class="text-muted-foreground hover:text-foreground">
				<List size={18} />
			</button>
			<span class="text-sm font-medium text-foreground"
				>lite<span class="text-muted-foreground">pay</span></span
			>
		</header>

		<main class="flex-1 overflow-auto p-6">
			{@render children()}
		</main>
	</div>
</div>
