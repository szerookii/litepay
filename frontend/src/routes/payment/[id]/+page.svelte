<script lang="ts">
	import CopyInput from "$lib/components/copy-input.svelte";
	import * as Card from "$lib/components/ui/card";
	import type { Payment } from "$lib/types";
	import QrCode from "@castlenine/svelte-qrcode";
	import { createQuery } from "@tanstack/svelte-query";

	const { data } = $props();

	let timeLeft = $state("");

	const paymentDataQuery = $state(
		createQuery<Payment>({
			queryKey: ["payment", data.payment.id],
			queryFn: async () => {
				const res = await fetch(`/api/payment/${data.payment.id}`);
				return res.json();
			},
			initialData: data.payment,
			refetchInterval: data.payment.status === "PENDING" ? 10000 : false
		})
	);

	const expiryDate = new Date($paymentDataQuery.data.expires_at).getTime();

	function updateTimer() {
		if (!expiryDate) {
			timeLeft = "Transaction expired";
			return;
		}

		const diff = expiryDate - Date.now();
		if (diff <= 0) {
			timeLeft = "Transaction expired";
			return;
		}

		const minutes = Math.floor(diff / 1000 / 60);
		const seconds = Math.floor((diff / 1000) % 60);
		timeLeft = `${minutes}:${seconds.toString().padStart(2, "0")}`;
	}

	setInterval(updateTimer, 1000);
	updateTimer();
</script>

<div class="flex h-screen w-full items-center justify-center">
	<Card.Root class="p-6 md:w-2/5">
		<Card.Header>
			<Card.Title class="text-center font-bold">LitePay</Card.Title>
			{#if $paymentDataQuery.data.status === "PENDING"}
				<Card.Description class="text-center"
					>Please complete your payment to continue</Card.Description
				>
			{:else if $paymentDataQuery.data.status === "PAID"}
				<Card.Description class="text-center">Your payment has been confirmed</Card.Description>
			{:else}
				<Card.Description class="text-center">Your payment expired</Card.Description>
			{/if}
		</Card.Header>
		<Card.Content class="justify-center text-center">
			{#if $paymentDataQuery.data.status === "PENDING"}
				<div class="space-y-0.5">
					<h3 class="text-sm text-muted-foreground">Payment Details</h3>
					<p class="text-sm">{$paymentDataQuery.data.id}</p>
				</div>
				<div class="mt-4">
					<h3 class="text-2xl font-extrabold">
						Send {$paymentDataQuery.data.amount_crypto}
						{$paymentDataQuery.data.currency_crypto_symbol}
					</h3>
					<p class="text-sm text-muted-foreground">
						≈ {$paymentDataQuery.data.amount_fiat.toFixed(2)}
						{$paymentDataQuery.data.currency_fiat}
					</p>
				</div>
				<div class="mt-4 md:mx-4 space-y-2">
					<p class="text-sm text-muted-foreground">To this recipient address</p>
					<CopyInput value={$paymentDataQuery.data.wallet_address} />
				</div>
				<div class="mt-4 flex flex-col justify-center space-y-2">
					<p class="text-sm text-muted-foreground">or scan the QR code to pay</p>
					<div class="flex w-full flex-wrap justify-center">
						<QrCode
							shape="circle"
							logoPath={`/imgs/svg/${$paymentDataQuery.data.currency_crypto_name.toLocaleLowerCase()}.svg`}
							data={`${$paymentDataQuery.data.currency_crypto_name.toLocaleLowerCase()}:${$paymentDataQuery.data.wallet_address}`}
						/>
					</div>
					{#if $paymentDataQuery.data.last_transaction_hash}
						<a
							href={`https://blockchair.com/${$paymentDataQuery.data.currency_crypto_name.toLocaleLowerCase()}/transaction/${$paymentDataQuery.data.last_transaction_hash}`}
							target="_blank"
							rel="noopener noreferrer"
							class="text-xs text-muted-foreground hover:underline"
						>
							View transaction details
						</a>
					{/if}
				</div>
				<Card.Root class="mt-8 border border-border bg-border/75 md:mx-4">
					<Card.Content class="flex items-center justify-between px-4 py-2">
						<div class="flex flex-col items-start space-y-1">
							<p class="text-xs text-muted-foreground">Expires in</p>
							<p class="text-sm font-medium text-white">
								{timeLeft}
							</p>
						</div>
						<div class="flex flex-col items-end space-y-1">
							<p class="text-xs text-muted-foreground">
								{$paymentDataQuery.data.confirmations !== null ? "Confirmations" : "Status"}
							</p>
							<p class="text-sm font-medium text-white">
								{$paymentDataQuery.data.confirmations !== null
									? `${$paymentDataQuery.data.confirmations} / ${$paymentDataQuery.data.required_confirmations}`
									: "Waiting..."}
							</p>
						</div>
					</Card.Content>	
				</Card.Root>
			{:else if $paymentDataQuery.data.status === "PAID"}
				<!-- Paid -->
				<div class="mt-4">
					<h3 class="text-2xl font-extrabold">Payment Confirmed</h3>
				</div>
			{:else}
				<!-- TODO: make it cleaner -->
				<div class="mt-4">
					<h3 class="text-2xl font-extrabold">Payment Expired</h3>
				</div>
			{/if}
		</Card.Content>
	</Card.Root>
</div>
