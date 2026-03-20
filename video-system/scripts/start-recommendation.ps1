$root = Split-Path -Parent $PSScriptRoot
$env:RECOMMEND_PORT = "8082"
$env:INDEX_BASE = "http://localhost:8083"
Set-Location "$root\\recommendation-platform"
go run .
