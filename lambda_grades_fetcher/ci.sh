#!/bin/bash

cd lambda_grades_fetcher || exit

case "$1" in
"install")
  ;;
"build")
  GOOS=linux go build -o bin/lambda_grades_fetcher ./*.go
  zip bin/lambda_grades_fetcher.zip bin/lambda_grades_fetcher
  ;;
"deploy")
  mkdir aws
  cd aws || exit
  curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "aws/awscliv2.zip"
  unzip aws/awscliv2.zip -d aws/
  ./aws/install -i ./aws/cli
  newpath="$(pwd)/aws/cli:$PATH"
  export PATH=$newpath
  aws lambda update-function-code \
    --function-name "$LAMBDA_GRADES_FETCHER_FUNCTION_NAME" \
    --zip-file fileb://bin/lambda_grades_fetcher.zip
  ;;
*)
  echo "Usage: ./ci.sh <command>"
  echo "Where <command> is \"install\", \"build\", or \"deploy\"."
  ;;
esac