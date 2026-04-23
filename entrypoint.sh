#!/bin/sh
set -e

echo "Checking Atlas..."
atlas version

echo "Applying migrations..."
atlas migrate apply --url "$DATABASE_URL"

echo "Starting LitePay..."
exec /litepay
