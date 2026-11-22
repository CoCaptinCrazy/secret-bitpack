$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $MyInvocation.MyCommand.Path
$goCommand = Get-Command go.exe -ErrorAction SilentlyContinue
if ($goCommand) {
    $go = $goCommand.Source
}
elseif (Test-Path 'C:\Program Files\Go\bin\go.exe') {
    $go = 'C:\Program Files\Go\bin\go.exe'
}
else {
    throw 'Go was not found. Install Go or add go.exe to PATH.'
}
Push-Location $root
try {
    & $go run .\cmd\bitpack @args
    $result = $LASTEXITCODE
}
finally {
    Pop-Location
}
exit $result
