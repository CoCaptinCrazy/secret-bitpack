$ErrorActionPreference = 'Stop'
$launcher = Join-Path $PSScriptRoot 'bp.ps1'
& $launcher @args
exit $LASTEXITCODE
