export type SEOConfig = {
	title: string;
	description: string;
	keywords?: string;
	ogImage?: string;
	ogType?: 'website' | 'product' | 'article';
	canonical?: string;
	noindex?: boolean;
	themeColor?: string;
	twitterHandle?: string;
};

export const BASE_URL = 'https://litepay.io';
export const TWITTER_HANDLE = '@litepayio';
export const THEME_COLOR_LIGHT = '#ffffff';
export const THEME_COLOR_DARK = '#000000';

export const SEO_CONFIG: Record<string, SEOConfig> = {
	'/': {
		title: 'LitePay - Accept Crypto Payments | No Middlemen',
		description: 'LitePay is an open-source crypto payment processor. Accept BTC, LTC, SOL on your own server. Self-hosted, API-first, fully automated.',
		keywords: 'crypto payments, bitcoin, litecoin, solana, payment processor, open source, self-hosted',
		ogImage: `${BASE_URL}/og-home.svg`,
		ogType: 'website',
		themeColor: THEME_COLOR_DARK,
		twitterHandle: TWITTER_HANDLE
	},
	'/auth/login': {
		title: 'Login - LitePay Merchant Dashboard',
		description: 'Sign in to your LitePay merchant dashboard to manage payments and wallets.',
		ogType: 'website',
		noindex: true,
		themeColor: THEME_COLOR_DARK,
		twitterHandle: TWITTER_HANDLE
	},
	'/auth/register': {
		title: 'Create Account - LitePay',
		description: 'Start accepting crypto payments with LitePay. Create your merchant account in minutes.',
		ogType: 'website',
		noindex: true,
		themeColor: THEME_COLOR_DARK,
		twitterHandle: TWITTER_HANDLE
	},
	'/dashboard': {
		title: 'Dashboard - LitePay',
		description: 'Your LitePay merchant dashboard. View account overview and quick links.',
		ogType: 'website',
		noindex: true,
		themeColor: THEME_COLOR_DARK,
		twitterHandle: TWITTER_HANDLE
	},
	'/dashboard/transactions': {
		title: 'Transactions - LitePay Dashboard',
		description: 'View and manage your crypto payment transactions.',
		ogType: 'website',
		noindex: true,
		themeColor: THEME_COLOR_DARK,
		twitterHandle: TWITTER_HANDLE
	},
	'/dashboard/wallets': {
		title: 'Wallets - LitePay Dashboard',
		description: 'Manage your wallets and webhook configuration.',
		ogType: 'website',
		noindex: true,
		themeColor: THEME_COLOR_DARK,
		twitterHandle: TWITTER_HANDLE
	},
	'/dashboard/api-keys': {
		title: 'API Keys - LitePay Dashboard',
		description: 'Manage your LitePay API credentials.',
		ogType: 'website',
		noindex: true,
		themeColor: THEME_COLOR_DARK,
		twitterHandle: TWITTER_HANDLE
	},
	'/dashboard/cashout': {
		title: 'Cashout - LitePay Dashboard',
		description: 'Withdraw your funds to external wallets.',
		ogType: 'website',
		noindex: true,
		themeColor: THEME_COLOR_DARK,
		twitterHandle: TWITTER_HANDLE
	},
	'/pay': {
		title: 'Payment - LitePay',
		description: 'Complete your payment securely with LitePay.',
		ogType: 'website',
		noindex: true,
		themeColor: THEME_COLOR_DARK,
		twitterHandle: TWITTER_HANDLE
	}
};

export function getSEOConfig(path: string): SEOConfig {
	return SEO_CONFIG[path] ?? SEO_CONFIG['/'];
}
