$root = Split-Path -Parent $PSScriptRoot
$env:PLATFORM_ID = "platformB"
$env:PLATFORM_PORT = "8084"
$env:INDEX_BASE = "http://localhost:8083"
$env:ACCEPT_TAGS = ""
Set-Location "$root\\video-platform"
go run .
