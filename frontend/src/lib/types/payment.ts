export interface Payment {
	id: string;
	wallet_address: string;
	amount_crypto: number;
	currency_crypto_name: string;
	currency_crypto_symbol: string;
	amount_fiat: number;
	currency_fiat: string;
	status: "PENDING" | "PAID" | "EXPIRED";
	expires_at: string;
	last_transaction_hash?: string | null;
	confirmations?: number | null;
	required_confirmations: number;
}
