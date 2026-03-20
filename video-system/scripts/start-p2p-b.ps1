$root = Split-Path -Parent $PSScriptRoot
$env:CHUNK_DIR = "$root\\video-platform\\data\\video-storage\\chunks"
$env:P2P_HTTP_PORT = "8091"
Set-Location "$root\\p2p-node"
go run .
