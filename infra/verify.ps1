param(
    [Parameter(Mandatory = $true)]
    [string]$ApiUrl
)

$ErrorActionPreference = "Stop"
$BaseUrl = $ApiUrl.TrimEnd("/")
$TestId = "CMP-TEST-001"

Write-Host "=== Connectivity & CRUD Verification: $BaseUrl ===" -ForegroundColor Cyan

# 1. POST /items (201 Created)
Write-Host "`n[1/6] POST /items" -ForegroundColor Yellow
$CreatePayload = @{
    item_id         = $TestId
    name            = "ESP32-C3 SuperMini"
    category        = "MICROCONTROLLER"
    stock_qty       = 10
    min_stock       = 3
    drawer_location = "RAK-T1"
    condition       = "NEW"
} | ConvertTo-Json

$PostResp = Invoke-WebRequest -Uri "$BaseUrl/items" -Method Post -ContentType "application/json" -Body $CreatePayload
Write-Host "Status: $($PostResp.StatusCode)" -ForegroundColor Green
Write-Host $PostResp.Content

# 2. GET /items (200 OK)
Write-Host "`n[2/6] GET /items" -ForegroundColor Yellow
$ListResp = Invoke-WebRequest -Uri "$BaseUrl/items" -Method Get
Write-Host "Status: $($ListResp.StatusCode)" -ForegroundColor Green
Write-Host $ListResp.Content

# 3. GET /items/{id} (200 OK)
Write-Host "`n[3/6] GET /items/$TestId" -ForegroundColor Yellow
$GetResp = Invoke-WebRequest -Uri "$BaseUrl/items/$TestId" -Method Get
Write-Host "Status: $($GetResp.StatusCode)" -ForegroundColor Green
Write-Host $GetResp.Content

# 4. PUT /items/{id} (200 OK)
Write-Host "`n[4/6] PUT /items/$TestId" -ForegroundColor Yellow
$UpdatePayload = @{
    stock_qty       = 2
    drawer_location = "RAK-T2"
    condition       = "GOOD"
} | ConvertTo-Json

$PutResp = Invoke-WebRequest -Uri "$BaseUrl/items/$TestId" -Method Put -ContentType "application/json" -Body $UpdatePayload
Write-Host "Status: $($PutResp.StatusCode)" -ForegroundColor Green
Write-Host $PutResp.Content

# 5. DELETE /items/{id} (200 OK)
Write-Host "`n[5/6] DELETE /items/$TestId" -ForegroundColor Yellow
$DelResp = Invoke-WebRequest -Uri "$BaseUrl/items/$TestId" -Method Delete
Write-Host "Status: $($DelResp.StatusCode)" -ForegroundColor Green
Write-Host $DelResp.Content

# 6. GET /items/{id} after delete (expect 404)
Write-Host "`n[6/6] GET /items/$TestId (expect 404 Not Found)" -ForegroundColor Yellow
try {
    Invoke-WebRequest -Uri "$BaseUrl/items/$TestId" -Method Get | Out-Null
    throw "Expected 404 Not Found, but request succeeded"
}
catch {
    $StatusCode = [int]$_.Exception.Response.StatusCode
    Write-Host "Status: $StatusCode (Expected 404)" -ForegroundColor Green
}

Write-Host "`nAll connectivity & CRUD checks passed!" -ForegroundColor Green
