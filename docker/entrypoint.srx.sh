#!/bin/sh

# Start /app/server and redirect its logs
/app/backend 2>&1 | tee /var/log/backend.log &

# Start frontend and redirect its logs
node /app/frontend 2>&1 | tee /var/log/frontend.log &

# Wait for all background processes to finish
# effectively blocks the CMD from finishing
wait
