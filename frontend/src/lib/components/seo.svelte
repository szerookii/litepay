<script lang="ts">
	import { getLocale, locales } from '$lib/paraglide/runtime.js';
	import { BASE_URL, TWITTER_HANDLE, type SEOConfig } from '$lib/seo/config.js';

	interface Props {
		config: SEOConfig;
		url: string;
		locale?: string;
	}

	let { config, url, locale = getLocale() }: Props = $props();

	const currentUrl = url === '/' ? BASE_URL : `${BASE_URL}${url}`;
	const canonicalUrl = url === '/' ? BASE_URL : currentUrl;
	const ogImage = config.ogImage ?? undefined;
	const themeColor = config.themeColor ?? '#000000';
	const twitterHandle = config.twitterHandle ?? TWITTER_HANDLE;
</script>

<svelte:head>
	<!-- Primary Meta Tags -->
	<title>{config.title}</title>
	<meta name="title" content={config.title} />
	<meta name="description" content={config.description} />
	{#if config.keywords}
		<meta name="keywords" content={config.keywords} />
	{/if}

	<!-- General Meta -->
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
	<meta name="language" content={locale} />
	<meta name="author" content="LitePay" />
	<meta name="creator" content="LitePay" />
	<meta name="robots" content={config.noindex ? 'noindex, nofollow' : 'index, follow'} />

	<!-- Canonical -->
	<link rel="canonical" href={canonicalUrl} />

	<!-- Open Graph / Facebook -->
	<meta property="og:type" content={config.ogType ?? 'website'} />
	<meta property="og:url" content={currentUrl} />
	<meta property="og:title" content={config.title} />
	<meta property="og:description" content={config.description} />
	{#if ogImage}
		<meta property="og:image" content={ogImage} />
		<meta property="og:image:width" content="1200" />
		<meta property="og:image:height" content="630" />
	{/if}
	<meta property="og:locale" content={locale === 'fr' ? 'fr_FR' : 'en_US'} />
	<meta property="og:site_name" content="LitePay" />

	<!-- Twitter / X -->
	<meta name="twitter:card" content={ogImage ? 'summary_large_image' : 'summary'} />
	<meta name="twitter:url" content={currentUrl} />
	<meta name="twitter:title" content={config.title} />
	<meta name="twitter:description" content={config.description} />
	{#if ogImage}
		<meta name="twitter:image" content={ogImage} />
	{/if}
	<meta name="twitter:creator" content={twitterHandle} />
	<meta name="twitter:site" content={twitterHandle} />

	<!-- Hreflang for multilingue -->
	{#each locales as lang}
		{@const hrefPath = lang === 'en' ? url : `/fr${url}`}
		<link rel="alternate" hreflang={lang} href={`${BASE_URL}${hrefPath}`} />
	{/each}
	<link rel="alternate" hreflang="x-default" href={`${BASE_URL}${url}`} />

	<!-- Theme Colors -->
	<meta name="theme-color" content={themeColor} media="(prefers-color-scheme: dark)" />
	<meta name="theme-color" content="#ffffff" media="(prefers-color-scheme: light)" />

	<!-- Additional SEO Meta -->
	<meta name="format-detection" content="telephone=no" />
	<meta property="article:author" content="LitePay" />
	<meta property="article:published_time" content={new Date().toISOString()} />

	<!-- Noindex if needed -->
	{#if config.noindex}
		<meta name="robots" content="noindex, nofollow" />
	{/if}
</svelte:head>
