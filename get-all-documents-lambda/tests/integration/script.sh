#!/bin/bash

# Get the absolute path of the current script
scriptPath="$(readlink -f "$0")"

################################################
# VARIABLES
################################################
projectName="lambda-golang-example"
lambdaName="get-all-documents-lambda"
dockerNetworkName="dynamodb-golang-lambda-network"
dockerImageName="amazon/dynamodb-local"
dynamodbDockerName="dynamodb-local-integration-test"
dynamodbPort="8001"
dynamodbTableName="users"
samFunctionName="getAllDocumentsFunction"

export AWS_DEFAULT_PROFILE="default"


################################################
# DOCKER IS RUNNING
################################################
if docker info > /dev/null 2>&1; then
  echo "docker is running"
else
  echo "docker is not running"
  exit 1
fi


################################################
# FUNCTIONS DEFINITION
################################################
function IsWordInCurrentPath() {
  workdir=$1
  keyword=$2

  if echo "$workdir" | grep -q "\<$keyword\>"; then
    echo "TRUE"
  else
    echo "FALSE"
  fi
}

function GetWorkDir() {
  workdir=$1
  keyword=$2

  # shellcheck disable=SC2016
  substring=$(echo "$workdir" | awk -F"$keyword" '{print $1}')
  echo "$substring$keyword"
}

function CleanContainer() {
  docker stop "$dynamodbDockerName"
  docker rm "$dynamodbDockerName"
  docker network rm "$dockerNetworkName"
}

function RunCMD() {
  typeset cmd="$1"
  typeset ret_code

  # shellcheck disable=SC2154
  eval "$cmd"
  ret_code=$?

  if [ $ret_code == 0 ]
  then
    echo "command executed successfully"
  else
    CleanContainer
    exit 1
  fi
}

function GetArchitecture() {
  arch=$(docker version --format '{{.Server.Arch}}')

  # return architecture
  echo "$arch"
}

function GetRuntime() {
  arch=$(docker version --format '{{.Server.Arch}}')
  runtime="go1.x"

  if [ "$arch" == "arm64" ]; then
      runtime="provided.al2"
  fi

  echo "$runtime"
}

function GetHandler() {
  arch=$(docker version --format '{{.Server.Arch}}')
  handler="main"

  if [ "$arch" == "arm64" ]; then
      handler="bootstrap"
  fi

  echo "$handler"
}


################################################
# Configure Workdir
################################################
echo "configuring workdir..."
lambdaWorkdir="$(IsWordInCurrentPath "$scriptPath" $lambdaName)"
projectWorkdir="$(IsWordInCurrentPath "$scriptPath" $projectName)"
workdir=""

if [ "$lambdaWorkdir" == "TRUE" ]; then
  echo "working on lambda workdir"
  workdir=$(GetWorkDir "$scriptPath" "$lambdaName")
elif [ "$projectWorkdir" == "TRUE" ]; then
  echo "working on project workdir"
  workdir=$(GetWorkDir "$scriptPath" "$projectName")
else
  echo "Invalid workdir"
  exit 1
fi

# validate last char
lastChar="${workdir: -1}"

if [ "$lastChar" == "/" ]; then
  workdir="${workdir%?}"
fi

echo "workdir is $workdir"


################################################
# Configure Container
################################################
echo "creating network"
docker network create -d bridge "$dockerNetworkName"
echo "creating container"
docker run -d -p "$dynamodbPort:8000" --name "$dynamodbDockerName" --network "$dockerNetworkName" "$dockerImageName"
dynamodbURL="http://localhost:$dynamodbPort"
echo "container created! - URL -> $dynamodbURL"

# Define the maximum number of retries
maxRetries=5
retryInterval=3

# Wait for the container to be ready
retryCount=0
while [ $retryCount -lt $maxRetries ]; do
  if docker inspect --format '{{.State.Status}}' "$dynamodbDockerName" | grep -q "running"; then
    echo "Container is ready."
    break
  else
    echo "Container is not ready yet. Retrying in $retryInterval seconds..."
    sleep $retryInterval
    ((retryCount++))
  fi
done

# Check if the maximum number of retries reached
if [ "$retryCount" -eq $maxRetries ]; then
  echo "Max retries reached. Container may not be ready."
fi


################################################
# Configure Table
################################################
echo "creating table $dynamodbTableName"
RunCMD "aws dynamodb create-table --endpoint-url $dynamodbURL --table-name $dynamodbTableName --attribute-definitions AttributeName=Id,AttributeType=S AttributeName=CreatedAt,AttributeType=S --key-schema AttributeName=Id,KeyType=HASH AttributeName=CreatedAt,KeyType=RANGE --provisioned-throughput ReadCapacityUnits=10,WriteCapacityUnits=10 --tags Key=Owner,Value=RicardoRomero --no-cli-pager"
echo "table created with name $dynamodbTableName"

## insert items
echo "inserting items ..."
RunCMD "aws dynamodb put-item --endpoint-url $dynamodbURL --table-name $dynamodbTableName --item '{\"Age\": {\"N\": \"30\"},\"CreatedAt\": {\"S\": \"2020-05-15T14:44:37.609166-06:00\"},\"Email\": {\"S\": \"alex.cruz@test.com\"},\"Id\": {\"S\": \"7d76ceb5-58f9-4263-8cf8-1f85664214c4\"},\"Lastname\": {\"S\": \"cruz\"},\"Name\": {\"S\": \"alejandro\"},\"UpdatedAt\": {\"S\": \"2024-04-14T13:44:37.609169-06:00\"}}'"
RunCMD "aws dynamodb put-item --endpoint-url $dynamodbURL --table-name $dynamodbTableName --item '{\"Age\": {\"N\": \"24\"},\"CreatedAt\": {\"S\": \"2021-05-15T14:44:37.609166-06:00\"},\"Email\": {\"S\": \"josep.cruz@test.com\"},\"Id\": {\"S\": \"7d76ceb5-58f9-4263-8cf8-1f85664214c5\"},\"Lastname\": {\"S\": \"cruz\"},\"Name\": {\"S\": \"josep\"},\"UpdatedAt\": {\"S\": \"2024-04-14T13:44:37.609169-06:00\"}}'"
RunCMD "aws dynamodb put-item --endpoint-url $dynamodbURL --table-name $dynamodbTableName --item '{\"Age\": {\"N\": \"18\"},\"CreatedAt\": {\"S\": \"2022-05-15T14:44:37.609166-06:00\"},\"Email\": {\"S\": \"miguel.cruz@test.com\"},\"Id\": {\"S\": \"7d76ceb5-58f9-4263-8cf8-1f85664214c6\"},\"Lastname\": {\"S\": \"cruz\"},\"Name\": {\"S\": \"miguel\"},\"UpdatedAt\": {\"S\": \"2024-04-14T13:44:37.609169-06:00\"}}'"
RunCMD "aws dynamodb put-item --endpoint-url $dynamodbURL --table-name $dynamodbTableName --item '{\"Age\": {\"N\": \"19\"},\"CreatedAt\": {\"S\": \"2023-05-15T14:44:37.609166-06:00\"},\"Email\": {\"S\": \"roberto.cruz@test.com\"},\"Id\": {\"S\": \"7d76ceb5-58f9-4263-8cf8-1f85664214c7\"},\"Lastname\": {\"S\": \"cruz\"},\"Name\": {\"S\": \"roberto\"},\"UpdatedAt\": {\"S\": \"2024-04-14T13:44:37.609169-06:00\"}}'"
RunCMD "aws dynamodb put-item --endpoint-url $dynamodbURL --table-name $dynamodbTableName --item '{\"Age\": {\"N\": \"17\"},\"CreatedAt\": {\"S\": \"2024-05-15T14:44:37.609166-06:00\"},\"Email\": {\"S\": \"john.cruz@test.com\"},\"Id\": {\"S\": \"7d76ceb5-58f9-4263-8cf8-1f85664214c8\"},\"Lastname\": {\"S\": \"cruz\"},\"Name\": {\"S\": \"john\"},\"UpdatedAt\": {\"S\": \"2024-04-14T13:44:37.609169-06:00\"}}'"
echo "items inserted"


################################################
# Run SAM CLI
################################################

# define sam cli variables
codeURI="$workdir/cmd"
integrationDir="$workdir/tests/integration"
buildPath="$integrationDir/.aws-sam/build"
samTemplatePath="$integrationDir/template.tmpl"
samFilePath="$integrationDir/template.yaml"
eventFilePath="$integrationDir/event.json"
architecture="$(GetArchitecture)"
runtime="$(GetRuntime)"
handler="$(GetHandler)"

mkdir -p "$buildPath"

## generate template file
echo "generating template, $buildPath"
echo "architecture $architecture, runtime $runtime, handler $handler"

dynamodbInternal="http://host.docker.internal:$dynamodbPort"
# Preprocess the SAM template
sed -e "s|{{ CODE_URI }}|$codeURI|g" \
    -e "s|{{ DYNAMODB_URL }}|$dynamodbInternal|g" \
    -e "s|{{ ARCHITECTURE }}|$architecture|g" \
    -e "s|{{ RUNTIME }}|$runtime|g" \
    -e "s|{{ HANDLER }}|$handler|g" \
    -e "s|{{ DYNAMODB_TABLE_NAME }}|$dynamodbTableName|g" \
    "$samTemplatePath" > "$samFilePath"
echo "template generated"

# sam cli build
samBuildPath="$buildPath/template.yaml"

# building program
echo "running sam cli"
## build sam cli
RunCMD "sam build $samFunctionName --template $samFilePath --build-dir $buildPath --no-cached"

## run sam cli
echo "invoke lambda"
RunCMD "sam local invoke $samFunctionName --template $samBuildPath --event $eventFilePath --docker-network $dockerNetworkName"

echo
echo "cleaning environment"
sleep 2
rm -rf "$integrationDir/.aws-sam"
rm -rf "$integrationDir/.sam-cli"
rm "$samFilePath"

CleanContainer
