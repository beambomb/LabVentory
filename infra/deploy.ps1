param(
    [string]$StackName = "lab-component-inventory",
    [string]$Region = "ap-southeast-1"
)

$ErrorActionPreference = "Stop"

Write-Host "=== LabComponentInventory AWS Deployment ===" -ForegroundColor Cyan

# 1. Build bootstrap.zip
& (Join-Path $PSScriptRoot "build.ps1")

# 2. Verify AWS credentials & resolve Account ID
$AccountId = aws sts get-caller-identity --query Account --output text --region $Region
if ($LASTEXITCODE -ne 0 -or [string]::IsNullOrWhiteSpace($AccountId)) {
    throw "AWS CLI is not configured. Run 'aws configure' first."
}

$BucketName = "labventroy-artifacts-$AccountId-$Region"
$ZipPath = Join-Path $PSScriptRoot "bootstrap.zip"
$S3Key = "lambda/bootstrap-$(Get-Date -Format 'yyyyMMddHHmmss').zip"

# 3. Ensure artifact bucket exists
Write-Host "Ensuring S3 artifact bucket s3://$BucketName exists..." -ForegroundColor Cyan
aws s3api head-bucket --bucket $BucketName --region $Region 2>$null
if ($LASTEXITCODE -ne 0) {
    if ($Region -eq "us-east-1") {
        aws s3api create-bucket --bucket $BucketName --region $Region | Out-Null
    } else {
        aws s3api create-bucket --bucket $BucketName --region $Region --create-bucket-configuration LocationConstraint=$Region | Out-Null
    }
}

# 4. Upload bootstrap.zip
Write-Host "Uploading package to s3://$BucketName/$S3Key..." -ForegroundColor Cyan
aws s3 cp $ZipPath "s3://$BucketName/$S3Key" --region $Region

# 5. Deploy CloudFormation stack
Write-Host "Deploying CloudFormation stack '$StackName'..." -ForegroundColor Cyan
aws cloudformation deploy `
    --template-file (Join-Path $PSScriptRoot "template.yaml") `
    --stack-name $StackName `
    --region $Region `
    --capabilities CAPABILITY_NAMED_IAM `
    --parameter-overrides CodeS3Bucket=$BucketName CodeS3Key=$S3Key

if ($LASTEXITCODE -ne 0) {
    throw "CloudFormation deployment failed"
}

Write-Host "`nDeployment Outputs:" -ForegroundColor Green
aws cloudformation describe-stacks `
    --stack-name $StackName `
    --region $Region `
    --query "Stacks[0].Outputs" `
    --output table
