$projectName = "lambda-golang-example"
$lambdaName = "get-all-documents-lambda"
$dockerNetworkName = "dynamodb-golang-lambda-network"
$dynamodbDockerName = "dynamodb-local-integration-test"
$dynamodbImageName = "amazon/dynamodb-local"
$dynamodbPort = 8080
$dynamodbTableName = "users"
$samFunctionName = "getAllDocumentsFunction"

$env:AWS_DEFAULT_PROFILE = "default"


################################################
# DOCKER IS RUNNING
################################################
$dockerProcess = Get-Process -Name "Docker Desktop"
if ($null -ne $dockerProcess) {
    Write-Host "Docker Desktop process is running."
} else {
    Write-Host "Docker Desktop process is not running."
    exit 1
}


################################################
# FUNCTIONS DEFINITION
################################################
$currentPath = $PWD.Path

function IsWordInCurrentPath {
    param (
        [string]$word
    )
    $keywordIndex = $currentPath.IndexOf($word)

    if ($keywordIndex -ge 0) {
        return $true
    } else {
        return $false
    }
}

function GetWorkDir {
    param (
        [string]$keyword
    )
    $keywordIndex = $currentPath.IndexOf($keyword)
    $workdir = $currentPath.Substring(0, $keywordIndex)
    return $workdir + $keyword
}

function clean_container {
    docker stop $dynamodbDockerName
    docker rm $dynamodbDockerName
    docker network rm $dockerNetworkName
}

function run_cmd {
    param (
        [string]$cmd
    )

    try {
        $process = Start-Process -FilePath "cmd.exe" -ArgumentList "/c $cmd" -NoNewWindow -PassThru -Wait
        $exitCode = $process.ExitCode

        if ($exitCode -eq 0) {
            Write-Host "command executed successfully"
        } else {
            Write-Host "invalid response $exitCode"
            clean_container
            exit 1
        }
    }
    catch {
        Write-Host "Error occurred while running the command $_"
        clean_container
        exit 1
    }
}

function GetArchitecture {
    $architecture = ""
    $runtime = "go1.x"
    $handler = "main"

    # Get Docker architecture based on environment variable
    $arch = $env:PROCESSOR_ARCHITECTURE

    switch -Regex ($arch) {
        "AMD64" {
            $architecture = "x86_64"
        }
        "ARM64" {
            $architecture = "arm64"
            $runtime = "provided.al2"
            $handler = "bootstrap"
        }
        "X86" { $architecture = "386" }
        default { $architecture = "unknown" }
    }

    return $architecture, $runtime, $handler
}


################################################
# Configure Workdir
################################################

$lambdaWorkDir = IsWordInCurrentPath -word $lambdaName
$projectWorkDir = IsWordInCurrentPath -word $projectName
$workdir = ""
if ($lambdaWorkDir) {
    Write-Host "working on lambda workdir"
    $workdir = GetWorkDir -word $lambdaName
} elseif ($projectWorkDir) {
    Write-Host "working on project workdir"
    $temp = GetWorkDir -keyword $projectName
    $workdir = Join-Path -Path $temp -ChildPath "$lambdaName"
} else {
    Write-Host "invalid workdir"
    exit 1
}

Write-Host "work-dir is $workdir"


################################################
# Configure Container
################################################

Write-Host "creating network"
docker network create -d bridge $dockerNetworkName
Write-Host "creating container"
$containerPort = $dynamodbPort.ToString() + ":8000"
docker run -d -p $containerPort --name $dynamodbDockerName --network $dockerNetworkName $dynamodbImageName
$dynamodbUrl = "http://localhost:${dynamodbPort}"
Write-Host "container created, url $dynamodbUrl"

# create table
Write-Host "creating table $dynamodbTableName"
run_cmd "aws dynamodb create-table --endpoint-url $dynamodbUrl --table-name $dynamodbTableName --attribute-definitions AttributeName=Id,AttributeType=S AttributeName=CreatedAt,AttributeType=S --key-schema AttributeName=Id,KeyType=HASH AttributeName=CreatedAt,KeyType=RANGE --provisioned-throughput ReadCapacityUnits=10,WriteCapacityUnits=10 --tags Key=Owner,Value=RicardoRomero"
Write-Host "table created with name $dynamodbTableName"

# insert items
Write-Host "inserting items"
$firstFile = Join-Path -Path $workdir -ChildPath "tests\integration\items\first.json"
$secondFile = Join-Path -Path $workdir -ChildPath "tests\integration\items\second.json"
$thirdFile = Join-Path -Path $workdir -ChildPath "tests\integration\items\third.json"
$fourthFile = Join-Path -Path $workdir -ChildPath "tests\integration\items\fourth.json"
$fifthFile = Join-Path -Path $workdir -ChildPath "tests\integration\items\fifth.json"

#run_cmd $firstItem
run_cmd "aws dynamodb put-item --endpoint-url $dynamodbUrl --table-name $dynamodbTableName --cli-input-json file://$firstFile"
run_cmd "aws dynamodb put-item --endpoint-url $dynamodbUrl --table-name $dynamodbTableName --cli-input-json file://$secondFile"
run_cmd "aws dynamodb put-item --endpoint-url $dynamodbUrl --table-name $dynamodbTableName --cli-input-json file://$thirdFile"
run_cmd "aws dynamodb put-item --endpoint-url $dynamodbUrl --table-name $dynamodbTableName --cli-input-json file://$fourthFile"
run_cmd "aws dynamodb put-item --endpoint-url $dynamodbUrl --table-name $dynamodbTableName --cli-input-json file://$fifthFile"
Write-Host "items inserted"


################################################
# Run SAM CLI
################################################

# define variables to use
$codeURI = Join-Path -Path $workdir -ChildPath "cmd"
$integrationDir = Join-Path -Path $workdir -ChildPath "tests\integration"
$buildPath = Join-Path -Path $integrationDir -ChildPath ".aws-sam\build"
$samTemplatePath = Join-Path -Path $integrationDir -ChildPath "template.tmpl"
$samFilePath = Join-Path -Path $integrationDir -ChildPath "template.yaml"
$architecture, $runtime, $handler = GetArchitecture

# create directory
New-Item -Path $buildPath -ItemType "Directory" -Force
Write-Host "directory created: $buildPath"

# Preprocess the SAM template
Write-Host "Pre-process the SAM template $samTemplatePath"
Get-Content $samTemplatePath | ForEach-Object {
    $_ -replace "{{ CODE_URI }}", "$codeURI" `
       -replace "{{ DYNAMODB_URL }}", "http://host.docker.internal:$dynamodbPort" `
       -replace "{{ ARCHITECTURE }}", "$architecture" `
       -replace "{{ RUNTIME }}", "$runtime" `
       -replace "{{ HANDLER }}", "$handler" `
       -replace "{{ DYNAMODB_TABLE_NAME }}", "$dynamodbTableName"
} | Set-Content $samFilePath

Write-Host "template generated " + $samFilePath

# sam cli vars
$templateBuildPath = Join-Path -Path $buildPath -ChildPath "template.yaml"
$eventFilePath = Join-Path -Path $integrationDir -ChildPath "event.json"

Write-Host "running sam cli for function" + $samFunctionName
# build sam cli
run_cmd "sam build $samFunctionName --template $samFilePath --build-dir $buildPath --no-cached"
Write-Host "program built"
# run sam cli
Write-Host "start lambda"
run_cmd "sam local start-lambda --template $templateBuildPath"

Write-Host "invoke lambda"
run_cmd "sam local invoke $samFunctionName --template $templateBuildPath --event $eventFilePath"

Write-Host ""
Write-Host "cleaning environment"
Start-Sleep -Seconds 1
Remove-Item -Path "$integrationDir\.aws-sam" -Recurse -Force
Remove-Item -Path "$samFilePath" -Force

clean_container
