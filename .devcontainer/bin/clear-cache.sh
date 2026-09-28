#!/bin/bash

# Configuration: Set your DB URL or rely on .env
# export DATABASE_URL="mysql://root:password@localhost:3306/my_db"

pushd /workspaces

echo "Starting to roll back all migrations..."

# Loop until dbmate down has nothing left to do
while true; do
    # Run dbmate down and capture output
    OUTPUT=$(dbmate down 2>&1)
    
    # Print the output
    echo "$OUTPUT"
    
    # Check if the output indicates "no migrations to rollback"
    if echo "$OUTPUT" | grep -q "Error: can't rollback: no migrations have been applied"; then
        echo "All migrations rolled back."
        break
    fi
    
    # Check for potential errors that should stop the loop
    if echo "$OUTPUT" | grep -q "Error"; then
        echo "Error detected. Stopping."
        exit 1
    fi
    
    # Optional: Small sleep to prevent hammering the database
    sleep 0.5
done

mc rm --recursive --force garage/${S3_BUCKET}/

dbmate up

bun run --cwd /workspaces/frontend better-auth-migrate

echo "Done."
popd