<script lang="ts">
	import { Copy, CopyCheck } from "lucide-svelte";
	import Button from "./ui/button/button.svelte";
	import { Input } from "$lib/components/ui/input";

	let copied = $state(false);

	interface Props {
		value: string;
	}

	const { value }: Props = $props();
</script>

<div class="relative flex items-center justify-center">
	<Input class="w-full truncate px-4 py-2.5 pr-12" {value} readonly />
	<Button
		variant="ghost"
		size="icon"
		class="absolute right-1 transition-colors duration-300"
		onclick={() => {
			navigator.clipboard.writeText(value);
			copied = true;
			setTimeout(() => {
				copied = false;
			}, 2000);
		}}
	>
		{#if copied}
			<CopyCheck  />
		{:else}
			<Copy />
		{/if}
	</Button>
</div>
