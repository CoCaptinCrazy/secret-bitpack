$ErrorActionPreference = 'Stop'

$goCommand = Get-Command go.exe -ErrorAction SilentlyContinue
if (-not $goCommand) {
    throw 'Go was not found. Install Go or add go.exe to PATH.'
}

$go = $goCommand.Source
Push-Location $PSScriptRoot
try {
    & $go install ./cmd/bitpack
    if ($LASTEXITCODE -ne 0) {
        throw "go install failed with exit code $LASTEXITCODE."
    }

    $installDirectory = (& $go env GOBIN).Trim()
    if ($LASTEXITCODE -ne 0) {
        throw 'Could not determine the Go binary install directory.'
    }
    if (-not $installDirectory) {
        $goPath = (& $go env GOPATH).Trim()
        if ($LASTEXITCODE -ne 0 -or -not $goPath) {
            throw 'Could not determine GOPATH.'
        }
        $goPathEntry = ($goPath -split [IO.Path]::PathSeparator)[0]
        $installDirectory = Join-Path $goPathEntry 'bin'
    }
    $installDirectory = [IO.Path]::GetFullPath($installDirectory)
    $executable = Join-Path $installDirectory 'bitpack.exe'
    if (-not (Test-Path $executable)) {
        throw "Bitpack was not created at $executable."
    }

    $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
    $userEntries = @()
    if ($userPath) {
        $userEntries = @($userPath -split ';' | Where-Object { $_.Trim() })
    }
    $normalizedInstallDirectory = $installDirectory.TrimEnd('\')
    $userHasInstallDirectory = $false
    foreach ($entry in $userEntries) {
        if ([string]::Equals($entry.TrimEnd('\'), $normalizedInstallDirectory, [StringComparison]::OrdinalIgnoreCase)) {
            $userHasInstallDirectory = $true
            break
        }
    }
    if (-not $userHasInstallDirectory) {
        $userEntries += $installDirectory
        [Environment]::SetEnvironmentVariable('Path', ($userEntries -join ';'), 'User')
    }

    $processEntries = @($env:Path -split ';' | Where-Object { $_.Trim() })
    $processHasInstallDirectory = $false
    foreach ($entry in $processEntries) {
        if ([string]::Equals($entry.TrimEnd('\'), $normalizedInstallDirectory, [StringComparison]::OrdinalIgnoreCase)) {
            $processHasInstallDirectory = $true
            break
        }
    }
    if (-not $processHasInstallDirectory) {
        $env:Path = (@($installDirectory) + $processEntries) -join ';'
    }

    Write-Host "Installed Bitpack to $executable"
    Write-Host 'Run: bitpack help'
    Write-Host 'Restart your terminal host to load the updated user PATH.'
}
finally {
    Pop-Location
}
