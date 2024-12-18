<script>
	import Button from "$lib/components/ui/button/button.svelte";
	import { LoaderCircle } from "lucide-svelte";

	let isLoading = $state(false);

	async function submit() {
		isLoading = true;

		try {
			const res = await fetch("/api/payment", {
				method: "POST",
				headers: {
					"Content-Type": "application/json",
					Authorization: `Bearer a09855b8ed60619ab1cd68ce9980e224`
				},
				body: JSON.stringify({
					symbol: "ltc",
					amount: 10,
					currency: "usd"
				})
			});

			if (!res.ok) {
				throw new Error("Failed to create payment");
			}

			const { id } = await res.json();
			window.location.href = `/payment/${id}`;
		} catch (error) {
			console.error(error);

			alert("Failed to create payment");
		} finally {
			isLoading = false;
		}
	}
</script>

<div class="flex h-screen w-full items-center justify-center">
	<Button size="default" onclick={submit} disabled={isLoading}>
		{#if isLoading}
			<LoaderCircle class="animate-spin" />
		{:else}
			<img src="/imgs/svg/litecoin.svg" alt="Litecoin" class="h-6 w-6" />
		{/if}
		Pay with Litecoin
	</Button>
</div>
