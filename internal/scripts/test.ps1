$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
Set-Location -Path $repoRoot

Write-Host "Goファイルを検索しています。"
$goFiles = Get-ChildItem -Path $repoRoot -Recurse -Filter *.go -File -ErrorAction SilentlyContinue |
    Select-Object -First 1
if (-not $goFiles) {
    Write-Host "Goファイルが見つからないため、テストをスキップします。"
    exit 0
}

Write-Host "go test ./... を実行します。"
go test ./...