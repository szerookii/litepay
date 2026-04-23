#!/bin/sh
# Run after each Vault container restart to unseal.
set -e

VAULT_ADDR=${VAULT_ADDR:-http://localhost:8200}
export VAULT_ADDR

printf "Unseal key: "
read -rs KEY
echo ""
vault operator unseal "$KEY"
