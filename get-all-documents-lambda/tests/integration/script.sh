#!/bin/bash

OWNER="ricardojonathanromero"
DOCKER_NETWORK_NAME="dynamodb-golang-lambda-network"
DYNAMODB_DOCKER_NAME="dynamodb-local-integration-test"
DYNAMODB_IMAGE="amazon/dynamodb-local"
DYNAMODB_PORT=8001
DYNAMODB_TABLE_NAME="users"
SAM_FUNCTION_NAME="getAllDocumentsFunction"

## clean gopath
ABS_GOLANG_PATH=""
GOLANG_PATH="$(go env GOPATH)"
paths=$(echo "$GOLANG_PATH" | tr ":" "\n")
for path in $paths
do
    # shellcheck disable=SC2034
    ABS_GOLANG_PATH="$path"
done

## get GOPATH
projectPath="$ABS_GOLANG_PATH/src/github.com/$OWNER/lambda-golang-example/get-all-documents-lambda"

## configure db

## create docker network
echo "creating network"
docker network create -d bridge "$DOCKER_NETWORK_NAME"
echo "creating container"
docker run -d --expose "8000" -p "$DYNAMODB_PORT:8000" --name "$DYNAMODB_DOCKER_NAME" --network "$DOCKER_NETWORK_NAME" "$DYNAMODB_IMAGE"
DYNAMODB_URL="http://localhost:$DYNAMODB_PORT"
echo "container created, url $DYNAMODB_URL"

function clean_container() {
    docker stop "$DYNAMODB_DOCKER_NAME"
    docker rm "$DYNAMODB_DOCKER_NAME"
    docker network rm "$DOCKER_NETWORK_NAME"
}

function run_cmd() {
    typeset cmd="$1"
    typeset ret_code

    eval "$cmd"
    ret_code=$?

    if [ $ret_code == 0 ]
    then
        echo "command executed successfully"
    else
        clean_container
        exit 1
    fi
}

# dummy credentials
export AWS_DEFAULT_PROFILE="default"

## create table
echo "creating table $DYNAMODB_TABLE_NAME"
run_cmd "aws dynamodb create-table --endpoint-url $DYNAMODB_URL --table-name $DYNAMODB_TABLE_NAME --attribute-definitions AttributeName=Id,AttributeType=S AttributeName=CreatedAt,AttributeType=S --key-schema AttributeName=Id,KeyType=HASH AttributeName=CreatedAt,KeyType=RANGE --provisioned-throughput ReadCapacityUnits=10,WriteCapacityUnits=10 --tags Key=Owner,Value=RicardoRomero --no-cli-pager"
echo "table created with name $DYNAMODB_TABLE_NAME"

## insert items
echo "inserting items"
run_cmd "aws dynamodb put-item --endpoint-url $DYNAMODB_URL --table-name $DYNAMODB_TABLE_NAME --item '{\"Age\": {\"N\": \"30\"},\"CreatedAt\": {\"S\": \"2020-05-15T14:44:37.609166-06:00\"},\"Email\": {\"S\": \"alex.cruz@test.com\"},\"Id\": {\"S\": \"7d76ceb5-58f9-4263-8cf8-1f85664214c4\"},\"Lastname\": {\"S\": \"cruz\"},\"Name\": {\"S\": \"alejandro\"},\"UpdatedAt\": {\"S\": \"2024-04-14T13:44:37.609169-06:00\"}}'"
run_cmd "aws dynamodb put-item --endpoint-url $DYNAMODB_URL --table-name $DYNAMODB_TABLE_NAME --item '{\"Age\": {\"N\": \"24\"},\"CreatedAt\": {\"S\": \"2021-05-15T14:44:37.609166-06:00\"},\"Email\": {\"S\": \"josep.cruz@test.com\"},\"Id\": {\"S\": \"7d76ceb5-58f9-4263-8cf8-1f85664214c5\"},\"Lastname\": {\"S\": \"cruz\"},\"Name\": {\"S\": \"josep\"},\"UpdatedAt\": {\"S\": \"2024-04-14T13:44:37.609169-06:00\"}}'"
run_cmd "aws dynamodb put-item --endpoint-url $DYNAMODB_URL --table-name $DYNAMODB_TABLE_NAME --item '{\"Age\": {\"N\": \"18\"},\"CreatedAt\": {\"S\": \"2022-05-15T14:44:37.609166-06:00\"},\"Email\": {\"S\": \"miguel.cruz@test.com\"},\"Id\": {\"S\": \"7d76ceb5-58f9-4263-8cf8-1f85664214c6\"},\"Lastname\": {\"S\": \"cruz\"},\"Name\": {\"S\": \"miguel\"},\"UpdatedAt\": {\"S\": \"2024-04-14T13:44:37.609169-06:00\"}}'"
run_cmd "aws dynamodb put-item --endpoint-url $DYNAMODB_URL --table-name $DYNAMODB_TABLE_NAME --item '{\"Age\": {\"N\": \"19\"},\"CreatedAt\": {\"S\": \"2023-05-15T14:44:37.609166-06:00\"},\"Email\": {\"S\": \"roberto.cruz@test.com\"},\"Id\": {\"S\": \"7d76ceb5-58f9-4263-8cf8-1f85664214c7\"},\"Lastname\": {\"S\": \"cruz\"},\"Name\": {\"S\": \"roberto\"},\"UpdatedAt\": {\"S\": \"2024-04-14T13:44:37.609169-06:00\"}}'"
run_cmd "aws dynamodb put-item --endpoint-url $DYNAMODB_URL --table-name $DYNAMODB_TABLE_NAME --item '{\"Age\": {\"N\": \"17\"},\"CreatedAt\": {\"S\": \"2024-05-15T14:44:37.609166-06:00\"},\"Email\": {\"S\": \"john.cruz@test.com\"},\"Id\": {\"S\": \"7d76ceb5-58f9-4263-8cf8-1f85664214c8\"},\"Lastname\": {\"S\": \"cruz\"},\"Name\": {\"S\": \"john\"},\"UpdatedAt\": {\"S\": \"2024-04-14T13:44:37.609169-06:00\"}}'"
echo "items inserted"

sleep 120

## initialize sam cli
dir="$projectPath/tests/integration"

mkdir -p "$dir/.sam-cli/build"

## generate template file
echo "generating template, $dir"
CODE_URI="$projectPath/cmd"
# Preprocess the SAM template
sed -e "s|{{ CODE_URI }}|$CODE_URI|g" \
    -e "s|{{ DYNAMODB_URL }}|$DOCKER_URL|g" \
    -e "s|{{ DYNAMODB_TABLE_NAME }}|$DYNAMODB_TABLE_NAME|g" \
    "$dir/template.tmpl" > "$dir/template.yaml"
echo "template generated"

echo ""
echo "running sam cli"
## build sam cli
run_cmd "sam build $SAM_FUNCTION_NAME --template $dir/template.yaml --build-dir $dir/.aws-sam/build --network $DOCKER_NETWORK_NAME"
## run sam cli
run_cmd "sam local invoke $SAM_FUNCTION_NAME --template $dir/.aws-sam/build/template.yaml --event $dir/event.json --network $DOCKER_NETWORK_NAME"

echo
echo "cleaning environment"
sleep 1
rm -rf "$dir/.aws-sam"
rm -rf "$dir/.sam-cli"
rm "$dir/template.yaml"

clean_container
