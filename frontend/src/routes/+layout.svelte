<script lang="ts">
	import "../app.css";

	import { ModeWatcher } from "mode-watcher";
	import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";

	const queryClient = new QueryClient();
	const { children } = $props();
</script>

<svelte:head>
	<meta name="darkreader-lock">
</svelte:head>

<ModeWatcher defaultMode={"dark"} />
<QueryClientProvider client={queryClient}>
	<main class="min-h-screen scroll-smooth">
		{@render children()}
	</main>
</QueryClientProvider>
