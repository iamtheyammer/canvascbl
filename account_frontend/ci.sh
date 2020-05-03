#!/bin/bash

cd account_frontend || exit

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
#  racv=$(NODE_DISABLE_COLORS=1 node -e 'console.log(Date.now())')
#  export REACT_APP_CURRENT_VERSION=$racv
#  echo "$ACCOUNT_CNAME" > build/CNAME
  ;;
"deploy")
  export AWS_ACCESS_KEY_ID=$2
  export AWS_SECRET_ACCESS_KEY=$3
  BUCKET_NAME=$4

  aws s3 sync ./build s3://"$BUCKET_NAME" --delete
  ;;
*)
  echo "Usage: ./ci.sh <command>"
  echo "Where <command> is \"install\", \"build\", \"before_deploy\", or \"deploy\"."
  ;;
esac