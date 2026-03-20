$root = Split-Path -Parent $PSScriptRoot
$env:STREAM_PORT = "8081"
$env:PLATFORM_MAP = "platformA=http://localhost:8080,platformB=http://localhost:8084"
$env:P2P_MAP = "platformA=http://localhost:8090,platformB=http://localhost:8091"
Set-Location "$root\\streaming-service"
go run .
