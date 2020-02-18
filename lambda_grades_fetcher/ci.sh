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
  export AWS_ACCESS_KEY_ID=$2
  export AWS_SECRET_ACCESS_KEY=$3
  LAMBDA_GRADES_FETCHER_FUNCTION_NAME=$4

  mkdir aws
  curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o aws/awscliv2.zip
  # pipe output to null
  unzip aws/awscliv2.zip -d aws/ > /dev/null
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