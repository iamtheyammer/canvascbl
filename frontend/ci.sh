#!/bin/bash

cd frontend || exit

case "$1" in
"install")
  nvm install "$(cat .nvmrc)"
  nvm use
  yarn install
  ;;
"build")
  yarn run formatcheck
	yarn run build
	echo "Built frontend. Output is in ./build"
  ;;
"before_deploy")
  racv=$(NODE_DISABLE_COLORS=1 node -e 'console.log(Date.now())')
  export REACT_APP_CURRENT_VERSION=$racv
  echo "$CNAME" > build/CNAME
  ;;
*)
  echo "Usage: ./ci.sh <command>"
  echo "Where <command> is \"install\", \"build\", or \"before_deploy\"."
  ;;
esac