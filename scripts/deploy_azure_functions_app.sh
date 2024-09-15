#!/bin/sh

set -eux

if [ -z "$FUNCTION_APP_NAME" ]; then
    echo "FUNCTION_APP_NAME is not set"
    echo "Please set the FUNCTION_APP_NAME environment variable e.g. export FUNCTION_APP_NAME=my-function-app"
    exit 1
fi

# Build the function app
make build GOOS=linux GOARCH=amd64 OUTPUT=assets/azurefunctions/dist/sandbox-go

# Deploy the function app
cd assets/azurefunctions
func azure functionapp publish "$FUNCTION_APP_NAME"
