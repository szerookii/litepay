#!/bin/sh
# Run once after first `docker compose up vault`.
# Initializes Vault, enables KV v2, and stores the master seed.
set -e

VAULT_ADDR=${VAULT_ADDR:-http://localhost:8200}
export VAULT_ADDR

echo "=== LitePay — Vault first-time setup ==="
echo ""

OUTPUT=$(vault operator init -key-shares=1 -key-threshold=1)
echo "$OUTPUT"
echo ""

UNSEAL_KEY=$(echo "$OUTPUT" | grep "Unseal Key 1" | awk '{print $NF}')
ROOT_TOKEN=$(echo "$OUTPUT" | grep "Initial Root Token" | awk '{print $NF}')

echo ">>> SAVE vault-keys.json SECURELY AND OFFLINE <<<"
echo "{\"unseal_key\":\"$UNSEAL_KEY\",\"root_token\":\"$ROOT_TOKEN\"}" > vault-keys.json
chmod 600 vault-keys.json
echo ""

vault operator unseal "$UNSEAL_KEY"
vault login "$ROOT_TOKEN"
vault secrets enable -path=secret kv-v2 2>/dev/null || true

echo ""
printf "Enter your BIP39 master seed (12 or 24 words): "
read -r MASTER_SEED
vault kv put secret/litepay master_seed="$MASTER_SEED"

echo ""
echo "=== Add to .env ==="
echo "SECRET_PROVIDER=vault"
echo "VAULT_ADDR=http://vault:8200"
echo "VAULT_TOKEN=$ROOT_TOKEN"
echo "VAULT_MOUNT=secret"
echo "VAULT_PATH=litepay"
echo "VAULT_KEY=master_seed"
