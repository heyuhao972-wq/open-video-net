$root = Split-Path -Parent $PSScriptRoot
Set-Location "$root\\video-index"
go run .
