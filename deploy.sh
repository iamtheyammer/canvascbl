#!/bin/bash

export HEROKU_API_KEY=$1
HEROKU_APP_NAME=$2
HEROKU_PROCESS_NAME=$3

# Install heroku if it doesn't exist
if [ ! -f "$(command -v heroku)" ]; then
    echo "Installing Heroku CLI..."
    curl https://cli-assets.heroku.com/install.sh | sh
fi

cd backend || exit 1

heroku container:login
heroku container:push "$HEROKU_PROCESS_NAME" -a "$HEROKU_APP_NAME"
heroku container:release "$HEROKU_PROCESS_NAME" -a "$HEROKU_APP_NAME"

cd ..