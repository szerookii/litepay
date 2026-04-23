![LitePay Cover](imgs/cover.png)

# LitePay

Self-hosted crypto payment processor for Bitcoin, Litecoin, and Solana. Accept payments from anywhere in the world without middlemen, fees, or third-party custody of funds.

**Key Features:**
- ✅ Multiple blockchains (BTC, LTC, SOL)
- ✅ Automatic HD wallet generation (BIP32/BIP39)
- ✅ API-first design with JWT + API Key auth
- ✅ Real-time blockchain verification (background worker)
- ✅ Webhooks with HMAC-SHA256 signing and automatic retry
- ✅ Embedded merchant dashboard & payment UI
- ✅ Multi-language support (EN/FR)
- ✅ Self-hosted infrastructure, full control

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Quick Deploy (Docker)](#quick-deploy-docker)
3. [Local Development](#local-development)
4. [Configuration](#configuration)
5. [API Documentation](#api-documentation)
6. [Payment Flow](#payment-flow)
7. [Webhooks](#webhooks)
8. [Database Schema](#database-schema)
9. [Architecture](#architecture)
10. [Troubleshooting](#troubleshooting)

---

## Prerequisites

- **Go** 1.25+ (for backend)
- **PostgreSQL** 17+ (for database)
- **Node.js** 20+ (for frontend builds)
- **pnpm** (frontend package manager)
- **Docker & Docker Compose** (recommended for deployment)

---

## Quick Deploy (Docker)

### Step 1 — Configure Environment

```bash
cp .env-example .env
```

Edit `.env` and set at minimum:

```env
POSTGRES_PASSWORD=your_strong_password
JWT_SECRET=random_string_at_least_32_chars_long_here
BTC_RPC_URL=https://lb.drpc.live/bitcoin/YOUR_KEY
LTC_RPC_URL=https://lb.drpc.live/litecoin/YOUR_KEY
SOL_RPC_URL=https://api.mainnet.solana.com
```

### Step 2 — Start Database & Vault

```bash
docker compose up -d postgres vault
```

### Step 3 — Initialize Vault (First Time Only)

```bash
VAULT_ADDR=http://localhost:8200 sh scripts/vault-init.sh
```

This script will:
1. Initialize Vault and save an **unseal key** to `vault-keys.json` — **KEEP THIS SAFE OFFLINE**
2. Prompt you to enter your BIP39 master seed (12 or 24 words)
3. Print values to add to `.env` (VAULT_TOKEN, etc.)

Copy the printed values into `.env`:

```env
SECRET_PROVIDER=vault
VAULT_ADDR=http://vault:8200
VAULT_TOKEN=<printed_by_script>
VAULT_MOUNT=secret
VAULT_PATH=litepay
VAULT_KEY=master_seed
```

### Step 4 — Start the App

```bash
docker compose up -d app
```

Access the app at `http://localhost:8080`

**After each server reboot:**

```bash
VAULT_ADDR=http://localhost:8200 sh scripts/vault-unseal.sh
# Enter the unseal key from vault-keys.json when prompted
docker compose up -d app
```

---

## Local Development

### 1. Install Dependencies

```bash
# Backend
go mod download

# Frontend
cd frontend
pnpm install
cd ..
```

### 2. Configure Environment

```bash
cp .env-example .env
# Edit .env with local settings
```

### 3. Start Database

```bash
docker compose -f docker-compose.dev.yml up -d postgres postgres_dev
```

### 4. Build Frontend (one-time)

```bash
cd frontend
pnpm build
cd ..
```

### 5. Run Backend

```bash
go run main.go
```

The app will be available at `http://localhost:8080`

---

## Configuration

### Required Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | - | PostgreSQL connection string (required) |
| `JWT_SECRET` | - | JWT signing secret, min 32 characters (required) |
| `SECRET_PROVIDER` | `env` | Where to load master seed: `env`, `vault`, `bitwarden`, `aws`, `gcp` |

### Optional Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `HOST` | `0.0.0.0` | Bind address (0.0.0.0 = public) |
| `PORT` | `8080` | Bind port |
| `ALLOW_REGISTER` | `true` | Enable new user registration via `/api/auth/register` |
| `ALLOWED_ORIGINS` | `http://localhost:5173` | CORS origins (comma-separated) |

### Blockchain RPC Endpoints

| Variable | Description | Example |
|----------|-------------|---------|
| `BTC_RPC_URL` | Bitcoin RPC endpoint | `https://lb.drpc.live/bitcoin/YOUR_KEY` |
| `LTC_RPC_URL` | Litecoin RPC endpoint | `https://lb.drpc.live/litecoin/YOUR_KEY` |
| `SOL_RPC_URL` | Solana RPC endpoint | `https://api.mainnet.solana.com` |

**Free Public RPC Services:**
- Bitcoin/Litecoin: [dRPC](https://drpc.live/)
- Solana: [Solana Labs](https://api.mainnet.solana.com) or [Helius](https://www.helius.dev/)

### Master Seed Management

The master seed is a BIP39 mnemonic (12 or 24 words) that derives all wallet addresses. **Losing it = losing all funds. Leaking it = draining all funds.**

Set `SECRET_PROVIDER` to choose where it lives:

#### `env` (default, simple)

The seed lives in your `.env` file. Protect it with `chmod 600` and encrypt your backups.

```env
SECRET_PROVIDER=env
MASTER_SEED=word1 word2 word3 ... word12
```

#### `vault` (recommended for servers)

HashiCorp Vault stores the seed AES-256 encrypted on disk. Requires an unseal key on startup.

```env
SECRET_PROVIDER=vault
VAULT_ADDR=http://your-vault:8200
VAULT_TOKEN=s.xxxxx
VAULT_MOUNT=secret
VAULT_PATH=litepay
VAULT_KEY=master_seed
```

See `scripts/vault-init.sh` and `scripts/vault-unseal.sh`

#### `bitwarden` (Bitwarden Secrets Manager)

```env
SECRET_PROVIDER=bitwarden
BITWARDEN_CLIENT_ID=...
BITWARDEN_CLIENT_SECRET=...
BITWARDEN_SECRET_ID=<uuid>
BITWARDEN_IDENTITY_URL=https://identity.bitwarden.com
BITWARDEN_API_URL=https://api.bitwarden.com
```

#### `aws` (AWS Secrets Manager)

```env
SECRET_PROVIDER=aws
AWS_REGION=us-east-1
AWS_SECRET_ID=litepay/master_seed
AWS_SECRET_KEY=master_seed  # if secret is JSON, extract this field
# AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY only needed outside AWS
```

#### `gcp` (GCP Secret Manager)

```env
SECRET_PROVIDER=gcp
GCP_PROJECT_ID=my-project
GCP_SECRET_NAME=litepay-master-seed
GCP_SECRET_VERSION=latest
# GOOGLE_APPLICATION_CREDENTIALS only needed outside GCP
```

---

## API Documentation

### Authentication

Two auth methods are supported:

1. **JWT (User Dashboard)**
   - Login via `/api/auth/login` to get JWT token
   - Use in header: `Authorization: Bearer <JWT_TOKEN>`
   - Valid for dashboard and user API endpoints

2. **API Key (Merchant API)**
   - Retrieve from user dashboard
   - Use in header: `Authorization: Bearer <API_KEY>`
   - Only for `/api/payment` endpoints

### Rate Limiting

- Registration: 5 per minute per IP
- Login: 10 per minute per IP
- Payment creation: 100 per minute per key

### Endpoints

#### **Authentication**

##### Register User
```http
POST /api/auth/register
Content-Type: application/json

{
  "email": "merchant@example.com",
  "password": "secure_password_min_8_chars"
}
```

**Response (201):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "merchant@example.com",
  "api_key": "sk_live_abcd1234..."
}
```

**Errors:**
- `400` - Invalid email or weak password
- `409` - Email already registered
- `503` - Registration disabled (ALLOW_REGISTER=false)

---

##### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "merchant@example.com",
  "password": "secure_password_min_8_chars"
}
```

**Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Errors:**
- `401` - Invalid email or password

---

##### Logout
```http
POST /api/auth/logout
Authorization: Bearer <JWT_TOKEN>
```

---

#### **Payment API (Merchant)**

##### Create Payment
```http
POST /api/payment
Authorization: Bearer <API_KEY>
Content-Type: application/json

{
  "symbol": "BTC",      # BTC, LTC, or SOL
  "amount": 49.99,      # Fiat amount
  "currency": "USD"     # USD or EUR
}
```

**Response (201):**
```json
{
  "id": "pay_550e8400e29b41d4a716446655440000",
  "wallet_address": "bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4",
  "amount_crypto": 0.00150000,
  "currency_crypto": "BTC",
  "amount_fiat": 49.99,
  "currency_fiat": "USD",
  "status": "PENDING",
  "expires_at": "2025-04-23T16:30:00Z"
}
```

**Errors:**
- `400` - Invalid symbol or currency
- `401` - Invalid API key
- `429` - Rate limit exceeded

---

##### Get Payment Status
```http
GET /api/payment/:id
```

**Response (200):**
```json
{
  "id": "pay_550e8400e29b41d4a716446655440000",
  "wallet_address": "bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4",
  "amount_crypto": 0.00150000,
  "currency_crypto": "BTC",
  "amount_fiat": 49.99,
  "currency_fiat": "USD",
  "received_amount": 0.00150000,
  "status": "PAID",
  "transaction_hash": "d64e0c1d8acfe1...",
  "expires_at": "2025-04-23T16:30:00Z",
  "created_at": "2025-04-23T15:30:00Z"
}
```

**Possible Statuses:**
- `PENDING` - Waiting for payment
- `CONFIRMING` - Payment received, awaiting confirmations
- `PAID` - Payment confirmed on-chain
- `EXPIRED` - Payment window closed, no funds received
- `REFUNDED` - Payment was refunded by merchant
- `CASHED_OUT` - Funds withdrawn by merchant

---

#### **Dashboard User API**

##### Get User Profile
```http
GET /api/user/me
Authorization: Bearer <JWT_TOKEN>
```

**Response (200):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "merchant@example.com",
  "api_key": "sk_live_abcd1234...",
  "supported_coins": ["BTC", "LTC", "SOL"],
  "webhook_url": "https://your-api.example.com/webhooks/litepay"
}
```

---

##### Get Account Balance
```http
GET /api/user/balance
Authorization: Bearer <JWT_TOKEN>
```

**Response (200):**
```json
{
  "BTC": 0.05300000,
  "LTC": 2.50000000,
  "SOL": 15.25000000
}
```

(Only confirmed PAID payments count toward balance)

---

##### Update Webhook URL
```http
PUT /api/user/wallets
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "webhook_url": "https://your-api.example.com/webhooks/litepay"
}
```

**Response (200):**
```json
{
  "webhook_url": "https://your-api.example.com/webhooks/litepay"
}
```

---

##### Get User Payments
```http
GET /api/user/payments
Authorization: Bearer <JWT_TOKEN>
```

**Response (200):**
```json
[
  {
    "id": "pay_550e8400e29b41d4a716446655440000",
    "wallet_address": "bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4",
    "amount_crypto": 0.00150000,
    "currency_crypto": "BTC",
    "amount_fiat": 49.99,
    "currency_fiat": "USD",
    "status": "PAID",
    "transaction_hash": "d64e0c1d8acfe1...",
    "created_at": "2025-04-23T15:30:00Z",
    "expires_at": "2025-04-23T16:30:00Z"
  }
]
```

---

##### Refund Payment
```http
POST /api/payment/:id/refund
Authorization: Bearer <JWT_TOKEN>
```

**Response (200):**
```json
{
  "tx_hash": "d64e0c1d8acfe1...",
  "to": "bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4"
}
```

**Errors:**
- `404` - Payment not found
- `400` - Payment cannot be refunded (not PAID, or already refunded)
- `500` - Blockchain error

---

##### Withdraw Funds (Cashout)
```http
POST /api/user/cashout
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "symbol": "BTC",
  "destination": "bc1q...external_address"
}
```

**Response (200):**
```json
{
  "transactions": [
    {
      "from_address": "bc1q...internal_address",
      "tx_hash": "d64e0c1d8acfe1...",
      "amount": 0.05000000
    }
  ]
}
```

---

##### Regenerate API Key
```http
POST /api/user/api-key
Authorization: Bearer <JWT_TOKEN>
```

**Response (200):**
```json
{
  "api_key": "sk_live_newkey123456789..."
}
```

---

#### **Public API**

##### Health Check
```http
GET /api/status
```

**Response (200):**
```json
{
  "status": "ok",
  "timestamp": "2025-04-23T15:30:00Z"
}
```

---

##### Get Public Config
```http
GET /api/config
```

**Response (200):**
```json
{
  "version": "0.1.0",
  "allow_register": true
}
```

---

## Payment Flow

### Step 1: Merchant Creates Payment

```bash
curl -X POST https://your-litepay.com/api/payment \
  -H "Authorization: Bearer sk_live_abc123..." \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTC",
    "amount": 100.00,
    "currency": "USD"
  }'
```

**Response:**
- Unique wallet address (HD-derived from master seed)
- Crypto amount (auto-converted from fiat)
- Payment ID

### Step 2: Customer Sends Payment

Customer sees payment page at `/pay?id=pay_xxx` with:
- Wallet address to send to
- Exact crypto amount needed
- QR code for mobile wallets
- Countdown timer (1 hour expiry)

### Step 3: Background Verifier Checks Blockchain

Every 30 seconds, background worker:
1. Queries all PENDING and CONFIRMING payments
2. Checks blockchain for transactions to each address
3. Updates status:
   - PENDING → CONFIRMING (transaction found)
   - CONFIRMING → PAID (required confirmations met)
   - PENDING → EXPIRED (1 hour passed)

### Step 4: Webhook Notification (Future)

⚠️ **Currently not implemented**, but the database schema is ready:

```json
POST https://merchant-webhook-url/litepay
Content-Type: application/json
X-Signature: hmac-sha256=...

{
  "event": "payment.paid",
  "id": "pay_550e8400e29b41d4a716446655440000",
  "status": "PAID",
  "amount_crypto": 0.00150000,
  "currency_crypto": "BTC",
  "transaction_hash": "d64e0c1d8acfe1...",
  "confirmed_at": "2025-04-23T15:35:00Z"
}
```

### Step 5: Merchant Polls or Waits

Merchant can:
- **Poll** `/api/payment/:id` to check status
- **Wait for webhook** (not yet implemented)
- Show customer success page once status = PAID

### Step 6: Optional Refund

If needed, merchant can refund via:

```bash
curl -X POST https://your-litepay.com/api/payment/pay_xxx/refund \
  -H "Authorization: Bearer JWT_TOKEN" \
```

This sends funds back to the original sender address.

---

## Webhooks

### Overview

Webhooks allow your application to receive real-time notifications when payment statuses change. When a payment transitions between states (PENDING → CONFIRMING → PAID), LitePay will POST an event to your configured webhook URL with signature verification.

### Current Status: ⚠️ Partial Implementation

**What's Done:**
- ✅ Webhook URL storage in user profile
- ✅ Webhook Secret generation and rotation
- ✅ UI to set/update webhook endpoint (`/dashboard/wallets`)
- ✅ UI to view and rotate webhook secrets
- ✅ Database schema ready for webhook events and delivery logs

**What's Missing:**
- ❌ Webhook sending on payment state changes
- ❌ Retry logic with exponential backoff
- ❌ Webhook event history/logs in UI
- ❌ Webhook test/ping endpoint

### Setting Up Webhooks (Dashboard)

1. **Login** to your merchant dashboard
2. Navigate to **Wallets** (`/dashboard/wallets`)
3. Enter your webhook URL in the "Webhook URL" field (e.g., `https://yourstore.com/api/webhooks/litepay`)
4. Click **Save Changes**
5. Copy your **Webhook Secret** — you'll need it to verify signatures
6. **Rotate Secret** anytime using the rotate button (invalidates the current secret immediately)

### Webhook Events

When implemented, webhooks will be sent as **POST requests** with the following payload:

```json
{
  "event": "payment.status_changed",
  "id": "pay_550e8400e29b41d4a716446655440000",
  "status": "PAID",
  "amount_crypto": 0.00150000,
  "currency_crypto": "BTC",
  "amount_fiat": 49.99,
  "currency_fiat": "USD",
  "received_amount": 0.00150000,
  "transaction_hash": "d64e0c1d8acfe1...",
  "confirmed_at": "2025-04-23T15:35:00Z",
  "created_at": "2025-04-23T15:30:00Z"
}
```

**Possible Status Values:**
- `PENDING` - Payment created, waiting for funds
- `CONFIRMING` - Payment received, awaiting blockchain confirmations
- `PAID` - Payment confirmed with sufficient confirmations
- `EXPIRED` - Payment window closed, no funds received
- `REFUNDED` - Merchant initiated refund
- `CASHED_OUT` - Merchant withdrew funds

### Webhook Signature Verification

All webhooks include an `X-LitePay-Signature` header with an **HMAC-SHA256 signature**. Verify it to ensure requests are from LitePay:

**Python Example:**
```python
import hmac
import hashlib

def verify_webhook(body: bytes, signature: str, secret: str) -> bool:
    expected = hmac.new(
        secret.encode(),
        body,
        hashlib.sha256
    ).hexdigest()
    return hmac.compare_digest(expected, signature)

# In your webhook handler:
signature = request.headers.get('X-LitePay-Signature')
if not verify_webhook(request.data, signature, YOUR_WEBHOOK_SECRET):
    return {"error": "Invalid signature"}, 401

payload = request.json
print(f"Payment {payload['id']} is now {payload['status']}")
```

**Node.js/Express Example:**
```javascript
const crypto = require('crypto');

function verifyWebhook(body, signature, secret) {
  const expected = crypto
    .createHmac('sha256', secret)
    .update(body)
    .digest('hex');
  return crypto.timingSafeEqual(expected, signature);
}

app.post('/api/webhooks/litepay', (req, res) => {
  const signature = req.headers['x-litepay-signature'];
  if (!verifyWebhook(req.rawBody, signature, process.env.LITEPAY_WEBHOOK_SECRET)) {
    return res.status(401).json({ error: 'Invalid signature' });
  }

  const { id, status, amount_crypto } = req.body;
  console.log(`Payment ${id} is now ${status} (${amount_crypto} crypto)`);
  
  // Update your order/invoice status here
  res.json({ success: true });
});
```

### Implementing Webhooks (Backend Development)

**Location:** `backend/cron/verify_transactions.go` line 109

The TODO comment marks where webhook sending should be added after a payment status changes.

**Required Implementation Steps:**

1. **Create Webhook Dispatcher Function**
   - After payment status update, marshal the payment data to JSON
   - Compute HMAC-SHA256 signature of the raw JSON body
   - Add `X-LitePay-Signature` header

2. **Add Webhook Delivery Table (Database)**
   - Track all webhook attempts for audit/debugging
   - Store: webhook_id, user_id, payment_id, event_type, status_code, response, attempts, next_retry_at
   - Indexed on next_retry_at for efficient polling

3. **Implement Retry Logic**
   - Max 5 retry attempts per webhook
   - Exponential backoff: 30s, 1min, 5min, 15min, 1hour
   - Only retry on transient failures (5xx, timeout, network error)
   - Skip retry on permanent failures (4xx except 408, 429)

4. **Create Webhook Worker**
   - Run every 5 seconds to process failed deliveries
   - Check deliveries with next_retry_at <= now
   - POST to merchant endpoint with signature
   - Update delivery log with response

5. **Optional: Add Test Endpoint**
   - `POST /api/user/webhook-test` - sends a sample webhook to test URL
   - Useful for merchants to validate their endpoint before go-live

### Webhook Best Practices

**For Merchants:**
- ✅ Always verify the signature before processing
- ✅ Idempotency: Ignore duplicate webhook events (use payment ID as key)
- ✅ Respond with 2xx status within 5 seconds
- ✅ Don't block on external API calls — queue the work
- ✅ Log all webhook events for debugging
- ❌ Don't use webhooks as your only payment status source — also poll `/api/payment/:id`

**For LitePay (Implementation):**
- Send webhooks for all status changes, not just PAID
- Include timestamp and payment details in every event
- Use persistent delivery queue to handle network failures
- Provide webhook history/logs in the dashboard
- Consider rate limiting: max 100 webhooks/user/minute

---

## Database Schema

### Users Table

| Column | Type | Notes |
|--------|------|-------|
| `id` | UUID | Primary key |
| `email` | String | Unique, normalized |
| `password_hash` | String | bcrypt hash |
| `api_key` | String | Unique, secret |
| `account_index` | Int | Unique HD wallet account |
| `webhook_url` | String | Nullable, for future use |
| `created_at` | Timestamp | |
| `updated_at` | Timestamp | |

### Payments Table

| Column | Type | Notes |
|--------|------|-------|
| `id` | UUID | Primary key, payment ID |
| `user_id` | UUID | Foreign key to users |
| `wallet_address` | String | Deposit address |
| `sol_reference` | String | Solana memo field (nullable) |
| `address_index` | Int | HD wallet derivation index |
| `amount_crypto` | Float | Requested crypto amount |
| `currency_crypto` | Enum | BTC, LTC, SOL |
| `amount_fiat` | Float | Requested fiat amount |
| `currency_fiat` | String | USD, EUR, etc |
| `received_amount` | Float | Actual received (nullable) |
| `transaction_hash` | String | On-chain tx hash (nullable) |
| `status` | Enum | PENDING, CONFIRMING, PAID, EXPIRED, REFUNDED, CASHED_OUT |
| `expires_at` | Timestamp | 1 hour after creation |
| `created_at` | Timestamp | |
| `updated_at` | Timestamp | |

### Status Lifecycle

```
PENDING ──(payment received)──> CONFIRMING ──(confirmations met)──> PAID
         \                                                         /
          ────────(1 hour pass)────> EXPIRED
                                                                /
         (merchant calls /refund) ──────────────────────> REFUNDED

(merchant calls /cashout) ──────────────────────> CASHED_OUT
```

---

## Architecture

### High-Level Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    Frontend (SvelteKit)                      │
│  Dashboard | Payment UI | i18n (EN/FR) | Tailwind CSS      │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTP/REST
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    Backend (Go/Gin)                          │
│  Auth (JWT+ApiKey) | Routes | Middleware | Rate Limiting   │
└────────┬──────────────────────┬──────────────────────┬───────┘
         │                      │                      │
    Database              Blockchain Cron         Secret Manager
      (PG)             (verify_transactions)       (Vault/AWS/GCP)
```

### Components

#### **Frontend** (`frontend/`)

- **SvelteKit** - Server-side rendering + client components
- **Paraglide.js** - Internationalization (EN/FR)
- **Tailwind CSS** - Utility-first styling
- **ShadCN-Svelte** - UI component library
- **Vite** - Build tool

**Key Pages:**
- `/auth/login` - User login
- `/auth/register` - User registration
- `/dashboard` - Main dashboard
- `/dashboard/transactions` - Payment history
- `/dashboard/wallets` - Webhook configuration
- `/dashboard/cashout` - Withdraw funds
- `/dashboard/api-keys` - API key management
- `/pay` - Public payment page

#### **Backend** (`backend/`)

- **Gin** - Web framework
- **pgx/v5** - PostgreSQL driver
- **Ent** - Entity framework (database abstraction)
- **JWT** - `golang-jwt/jwt/v5`
- **Validation** - `go-playground/validator`

**Structure:**
```
backend/
├── db/          # Database functions
├── crypto/      # Blockchain RPC clients
├── cron/        # Background workers (verify_transactions)
├── router/      # HTTP routes
├── ent/         # Database schema (Ent generated)
├── config/      # Configuration
├── models/      # Data structures
└── middleware/  # Auth, CORS, rate limit
```

**Background Worker: Transaction Verifier**
- Runs every 30 seconds
- Queries all PENDING/CONFIRMING payments
- Calls blockchain RPC to check for transactions
- Updates payment status in database
- **TODO: Send webhooks here** (line 109)

#### **Database** (`PostgreSQL`)

- **Migrations**: Automatic via Atlas
- **Schema**: Defined in `backend/ent/schema/`
- **Backup**: Use `pg_dump` or volume backups

#### **Secret Management**

Pluggable secret providers:
- `env` - .env file (simple)
- `vault` - HashiCorp Vault (enterprise)
- `bitwarden` - Bitwarden Secrets Manager
- `aws` - AWS Secrets Manager
- `gcp` - GCP Secret Manager

### Language Support

Frontend fully internationalized via Inlang Paraglide:

**Supported Languages:**
- English (en)
- Français (fr)

**To add a language:**
1. Add translation JSON in `frontend/messages/{lang}.json`
2. Paraglide auto-generates language functions
3. Use `m.key_name()` in components

---

## Troubleshooting

### Common Issues

#### "required variable POSTGRES_PASSWORD is missing a value"

**Fix:** Use correct Docker Compose flag:
```bash
docker compose -f docker-compose.dev.yml up -d
```

---

#### "Connection refused" when accessing app

**Check if containers are running:**
```bash
docker compose ps
```

**Check logs:**
```bash
docker compose logs app
```

---

#### "Invalid master seed"

**Ensure:**
1. Seed is 12 or 24 words
2. Words are separated by spaces
3. Words are valid BIP39 words

---

#### "Payment not appearing"

**Causes:**
1. RPC endpoint down - check `BTC_RPC_URL`, `LTC_RPC_URL`, `SOL_RPC_URL`
2. Background worker not running - check `docker compose logs app`
3. Transaction not sent to correct address

**Debug:**
```bash
# Check payment status
curl https://your-litepay.com/api/payment/pay_xxx

# Check background worker logs
docker compose logs app | grep verify
```

---

#### "Vault won't unseal"

**Steps:**
1. Find unseal key in `vault-keys.json`
2. Run: `VAULT_ADDR=http://localhost:8200 sh scripts/vault-unseal.sh`
3. Paste key when prompted

---

#### "CORS errors when calling API from browser"

**Fix:** Set `ALLOWED_ORIGINS` in `.env`:

```env
ALLOWED_ORIGINS=http://localhost:3000,https://your-frontend.com
```

---

## Development

### Building Frontend

```bash
cd frontend
pnpm install
pnpm build
cd ..
```

### Running Tests

```bash
# Backend
go test ./...

# Frontend
cd frontend
pnpm test
```

### Database Migrations

Migrations are automatic via Atlas. To create a new migration:

```bash
atlas migrate diff --env local
```

---

## Security Considerations

1. **Master Seed** - Never hardcode, use secret manager
2. **API Keys** - Rotate regularly from dashboard
3. **HTTPS** - Always use in production
4. **CORS** - Whitelist only needed origins
5. **Rate Limiting** - Built-in on auth endpoints
6. **JWT Secret** - Min 32 chars, cryptographically random

---

## License

MIT - See LICENSE file

---

## Support & Community

- **Issues**: Report bugs on [GitHub Issues](https://github.com/szerookii/litepay/issues)
- **Discussions**: [GitHub Discussions](https://github.com/szerookii/litepay/discussions)
- **Documentation**: Full API docs in this README

---

## Roadmap

- [ ] Webhook execution with retry logic
- [ ] More payment status events
- [ ] Payment splitting/invoicing
- [ ] Advanced analytics & reports
- [ ] More languages (ES, DE, IT)
- [ ] Mobile app for payment confirmation
- [ ] Lightning Network support

---

**Made with ❤️ by the LitePay team**
