$root = Split-Path -Parent $PSScriptRoot
Set-Location "$root\\client"
python -m http.server 5173
