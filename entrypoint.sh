#!/bin/sh
set -e

# Wait for database to be ready
until nc -z postgres 5432; do
  echo "Waiting for postgres..."
  sleep 1
done

# Apply migrations
echo "Applying database migrations with Atlas..."
atlas migrate apply --url "$DATABASE_URL"

# Start the application
echo "Starting LitePay..."
exec /litepay
