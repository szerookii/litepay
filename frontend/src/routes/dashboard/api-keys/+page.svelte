<script lang="ts">
	import { onMount } from 'svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Copy, Check, ArrowsClockwise, Warning, Eye, EyeSlash } from 'phosphor-svelte';
	import * as m from '$lib/paraglide/messages.js';

	let apiKey = $state('');
	let revealed = $state(false);
	let copied = $state(false);
	let regenerating = $state(false);
	let confirmRegen = $state(false);
	let error = $state('');

	const maskedKey = $derived(
		apiKey ? `${apiKey.slice(0, 8)}${'•'.repeat(24)}${apiKey.slice(-8)}` : ''
	);
	const displayKey = $derived(revealed ? apiKey : maskedKey);

	async function load(): Promise<void> {
		const token = localStorage.getItem('token');
		if (!token) return;
		const res = await fetch('/api/user/me', {
			headers: { Authorization: `Bearer ${token}` }
		});
		if (res.ok) {
			const me = await res.json();
			apiKey = me.api_key;
		}
	}

	function copyKey(): void {
		navigator.clipboard.writeText(apiKey);
		copied = true;
		setTimeout(() => (copied = false), 2000);
	}

	async function regenerate(): Promise<void> {
		if (!confirmRegen) {
			confirmRegen = true;
			return;
		}
		regenerating = true;
		error = '';
		confirmRegen = false;
		const token = localStorage.getItem('token');
		try {
			const res = await fetch('/api/user/api-key', {
				method: 'POST',
				headers: { Authorization: `Bearer ${token}` }
			});
			if (res.ok) {
				const data = await res.json();
				apiKey = data.api_key;
				revealed = true;
			} else {
				error = 'Failed to regenerate key';
			}
		} catch {
			error = 'Network error';
		} finally {
			regenerating = false;
		}
	}

	onMount(load);
</script>

<div class="max-w-2xl space-y-8">
	<div class="space-y-1">
		<p class="text-xs uppercase tracking-widest text-muted-foreground">{m.dash_nav_api_keys()}</p>
		<h1 class="text-xl font-medium text-foreground">{m.api_keys_title()}</h1>
		<p class="text-xs text-muted-foreground">{m.api_keys_subtitle()}</p>
	</div>

	<div class="ring-border space-y-5 p-5 ring-1">
		<div class="flex items-center justify-between">
			<div>
				<p class="text-xs font-medium text-foreground">{m.api_keys_label()}</p>
				<p class="mt-0.5 text-xs text-muted-foreground">
					Send as <code class="bg-muted px-1 py-0.5 font-mono text-[10px]"
						>Authorization: Bearer &lt;key&gt;</code
					>
				</p>
			</div>
			<Badge variant="outline">Active</Badge>
		</div>

		<Separator />

		<!-- Key display -->
		<div class="ring-border flex items-center gap-2 bg-muted/20 px-3 py-2.5 ring-1">
			<span class="flex-1 truncate font-mono text-xs text-foreground">{displayKey || '…'}</span>
			<div class="flex shrink-0 items-center gap-1">
				<button
					onclick={() => (revealed = !revealed)}
					class="p-1 text-muted-foreground transition-colors hover:text-foreground"
					aria-label={revealed ? m.api_keys_hide() : m.api_keys_reveal()}
				>
					{#if revealed}
						<EyeSlash size={13} />
					{:else}
						<Eye size={13} />
					{/if}
				</button>
				<button
					onclick={copyKey}
					class="p-1 text-muted-foreground transition-colors hover:text-foreground"
					aria-label={m.api_keys_copy()}
				>
					{#if copied}
						<Check size={13} color="#4ade80" />
					{:else}
						<Copy size={13} />
					{/if}
				</button>
			</div>
		</div>

		{#if copied}
			<p class="text-xs text-green-400">{m.api_keys_copied()}</p>
		{/if}
	</div>

	<!-- Usage example -->
	<div class="space-y-2">
		<p class="text-xs uppercase tracking-widest text-muted-foreground">Usage</p>
		<div class="ring-border overflow-hidden ring-1">
			<div class="border-border flex items-center gap-2 border-b bg-muted/50 px-4 py-2">
				<Badge variant="secondary" class="text-[10px]">POST</Badge>
				<span class="text-xs text-muted-foreground">/api/payment</span>
			</div>
			<pre
				class="overflow-x-auto bg-card p-4 text-xs leading-relaxed text-foreground"><code><span class="text-muted-foreground"># Create a payment</span>
<span class="text-green-400">$</span> curl -X POST https://your-litepay.com/api/payment \
    -H <span class="text-yellow-400">"Authorization: Bearer {displayKey.slice(0, 12)}…"</span> \
    -H <span class="text-yellow-400">"Content-Type: application/json"</span> \
    -d <span class="text-yellow-400">'&#123;"symbol":"BTC","amount":49.99,"currency":"USD"&#125;'</span></code></pre>
		</div>
	</div>

	<Separator />

	<!-- Regenerate -->
	<div class="space-y-3">
		<div>
			<p class="text-xs font-medium text-foreground">{m.api_keys_regenerate()}</p>
			<p class="mt-0.5 text-xs text-muted-foreground">
				{m.api_keys_warning()}
			</p>
		</div>

		{#if error}
			<p class="flex items-center gap-1.5 text-xs text-destructive">
				<Warning size={12} />
				{error}
			</p>
		{/if}

		{#if confirmRegen}
			<div class="ring-border bg-destructive/5 p-3 ring-1 ring-destructive/20">
				<p class="mb-3 text-xs text-foreground">
					{m.api_keys_warning()}
				</p>
				<div class="flex gap-2">
					<Button
						variant="destructive"
						size="sm"
						onclick={regenerate}
						disabled={regenerating}
						class="gap-1.5"
					>
						<ArrowsClockwise size={12} />
						{regenerating ? m.api_keys_regenerating() : m.api_keys_regenerate()}
					</Button>
					<Button variant="ghost" size="sm" onclick={() => (confirmRegen = false)}>Cancel</Button>
				</div>
			</div>
		{:else}
			<Button variant="outline" size="sm" onclick={regenerate} class="gap-1.5">
				<ArrowsClockwise size={12} />
				{m.api_keys_regenerate()}
			</Button>
		{/if}
	</div>
</div>
