#!/bin/bash

cd teacher_frontend || exit

[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"  # This loads nvm
nvm install "$(cat .nvmrc)"
nvm use

case "$1" in
"install")
  yarn install
  ;;
"build")
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