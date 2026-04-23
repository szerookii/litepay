![LitePay Cover](imgs/cover.png)

# LitePay

Self-hosted crypto payment processor for Bitcoin, Litecoin, and Solana. No middlemen, no fees, no third-party custody.

**Features:**
- Multiple chains (BTC, LTC, SOL)
- Automatic wallet generation
- API-first design
- Blockchain verification workers
- Embedded dashboard & payment UI
- Self-hosted infrastructure

## Prerequisites

- Go 1.21+
- PostgreSQL or SQLite
- Node.js (for building frontend)

## Quick Start

Copy `.env-example` to `.env`:
```env
HOST=0.0.0.0
PORT=8080
DATABASE_URL=postgres://user:pass@localhost:5432/litepay
JWT_SECRET=your_secret_key
MASTER_SEED=your_seed_phrase
ALLOW_REGISTER=true
BTC_RPC_URL=https://...
LTC_RPC_URL=https://...
SOL_RPC_URL=https://...
```

Build and run:
```bash
cd frontend && pnpm install && pnpm build && cd ..
go run main.go
```

## Integration

```bash
curl -X POST https://your-litepay.com/api/payment \
    -H "Authorization: Bearer YOUR_API_KEY" \
    -H "Content-Type: application/json" \
    -d '{"symbol":"BTC","amount":49.99,"currency":"USD"}'
```

## Security

LitePay supports optional KMS integration for key management. See documentation for setup details. Always secure your deployment and protect your `MASTER_SEED`.
