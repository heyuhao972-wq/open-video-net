$root = Split-Path -Parent $PSScriptRoot
$env:PLATFORM_ID = "platformA"
$env:PLATFORM_PORT = "8080"
$env:INDEX_BASE = "http://localhost:8083"
$env:ACCEPT_TAGS = "tech,ai,technology"
Set-Location "$root\\video-platform"
go run .
