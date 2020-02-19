#!/bin/bash

source scripts/diff.sh

if [ ! "$(diff_includes backend)" ]; then
  echo "Skipping backend step because backend was not changed."
  exit
fi

cd backend || exit

case "$1" in
"install")
  ;;
"build")
  go build -o bin/backend src/*.go
  ;;
"deploy")
  export HEROKU_API_KEY=$2
  HEROKU_APP_NAME=$3
  HEROKU_PROCESS_NAME=$4

  # Install heroku if it doesn't exist
  if [ ! -f "$(command -v heroku)" ]; then
      echo "Installing Heroku CLI..."
      curl https://cli-assets.heroku.com/install.sh | sh
  fi

  heroku container:login
  heroku container:push "$HEROKU_PROCESS_NAME" -a "$HEROKU_APP_NAME"
  heroku container:release "$HEROKU_PROCESS_NAME" -a "$HEROKU_APP_NAME"
  ;;
*)
  echo "Usage: ./ci.sh <command>"
  echo "Where <command> is \"install\", \"build\", or \"deploy\"."
  ;;
esac