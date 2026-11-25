$ErrorActionPreference = "Stop"

$RootDir = Split-Path -Parent $PSScriptRoot
$BackendDir = Join-Path $RootDir "backend"
$OutputBin = Join-Path $PSScriptRoot "bootstrap"
$OutputZip = Join-Path $PSScriptRoot "bootstrap.zip"

Write-Host "[1/2] Cross-compiling Go Lambda binary (linux/amd64)..." -ForegroundColor Cyan
Push-Location $BackendDir
try {
    $env:GOOS = "linux"
    $env:GOARCH = "amd64"
    $env:CGO_ENABLED = "0"
    go build -tags lambda.norpc -ldflags="-s -w" -o $OutputBin ./cmd/lambda
    if ($LASTEXITCODE -ne 0) {
        throw "go build failed"
    }
}
finally {
    Remove-Item Env:GOOS, Env:GOARCH, Env:CGO_ENABLED -ErrorAction SilentlyContinue
    Pop-Location
}

Write-Host "[2/2] Packaging bootstrap.zip..." -ForegroundColor Cyan
if (Test-Path $OutputZip) {
    Remove-Item $OutputZip -Force
}
Compress-Archive -Path $OutputBin -DestinationPath $OutputZip -CompressionLevel Optimal
Remove-Item $OutputBin -Force

Write-Host "Build complete: $OutputZip" -ForegroundColor Green
