#!/bin/bash

cd teacher_frontend || exit

case "$1" in
"install")
#  frontend has installed nvm
#  curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.35.2/install.sh | bash
  export NVM_DIR="$HOME/.nvm"
  [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"  # This loads nvm
  nvm install "$(cat .nvmrc)"
  nvm use
  yarn install
  ;;
"build")
  nvm use
  yarn run formatcheck || exit 2
	yarn run build || exit 2
	echo "Built frontend. Output is in ./build"
  ;;
"before_deploy")
  racv=$(NODE_DISABLE_COLORS=1 node -e 'console.log(Date.now())')
  export REACT_APP_CURRENT_VERSION=$racv
  echo "$TEACHER_CNAME" > build/CNAME
  ;;
*)
  echo "Usage: ./ci.sh <command>"
  echo "Where <command> is \"install\", \"build\", or \"before_deploy\"."
  ;;
esac