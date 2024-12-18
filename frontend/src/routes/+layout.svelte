<script lang="ts">
	import "../app.css";

	import { ModeWatcher } from "mode-watcher";
	import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";

	const queryClient = new QueryClient();
	const { children } = $props();
</script>

<ModeWatcher defaultMode={"dark"} />

<QueryClientProvider client={queryClient}>
	<main class="min-h-screen scroll-smooth">
		{@render children()}
	</main>
</QueryClientProvider>
