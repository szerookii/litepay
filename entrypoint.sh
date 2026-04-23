#!/bin/sh
set -e

# Apply database migrations using Atlas
echo "Applying database migrations..."
/atlas migrate apply --url "$DATABASE_URL"

# Start the application
exec /litepay
