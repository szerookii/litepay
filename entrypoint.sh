#!/bin/sh
set -e

echo "Applying migrations..."
for i in 1 2 3 4 5; do
  if atlas migrate apply --url "$DATABASE_URL" 2>/dev/null; then
    break
  fi
  echo "Waiting for DB..."
  sleep 5
done

echo "Starting LitePay..."
exec /litepay