<script lang="ts">
	import Button from "$lib/components/ui/button/button.svelte";
	import { LoaderCircle } from "lucide-svelte";

	let loadingStates = $state({
		ltc: false,
		bch: false
	});

	async function submit(symbol: string, amount: number = 10) {
		loadingStates[symbol] = true;

		try {
			const res = await fetch("/api/payment", {
				method: "POST",
				headers: {
					"Content-Type": "application/json",
					Authorization: `Bearer a09855b8ed60619ab1cd68ce9980e224`
				},
				body: JSON.stringify({
					symbol: symbol,
					amount: amount,
					currency: "usd"
				})
			});

			if (!res.ok) {
				throw new Error("Failed to create payment");
			}

			const { id } = await res.json();
			window.location.href = `/payment?id=${id}`;
		} catch (error) {
			console.error(error);

			alert("Failed to create payment");
		} finally {
			loadingStates[symbol] = false;
		}
	}
</script>

<div class="flex gap-4 h-screen w-full items-center justify-center">
	<Button size="default" onclick={() => submit("ltc")} disabled={loadingStates.ltc}>
		{#if loadingStates.ltc}
			<LoaderCircle class="animate-spin" />
		{:else}
			<img src="/imgs/svg/litecoin.svg" alt="Litecoin" class="h-6 w-6" />
		{/if}
		Pay with Litecoin
	</Button>
	<Button size="default" onclick={() => submit("bch", 0.5)} disabled={loadingStates.bch}>
		{#if loadingStates.bch}
			<LoaderCircle class="animate-spin" />
		{:else}
			<img src="/imgs/svg/bitcoin-cash.svg" alt="BitcoinCash" class="h-6 w-6" />
		{/if}
		Pay with Bitcoin Cash
	</Button>
</div>
