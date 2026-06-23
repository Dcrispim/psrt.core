# Installs Chrome stable next to the Wails binary when none is bundled yet.
# Usage: install-chrome.ps1 <path-to-psrt-gui.exe>
param(
    [Parameter(Mandatory = $true, Position = 0)]
    [string]$BinPath
)

$ErrorActionPreference = 'Stop'
$Dir = Split-Path -Parent (Resolve-Path $BinPath)

function Test-BundledChrome {
    param([string]$Root)
    $names = @(
        'chrome.exe',
        'chromium.exe',
        'chrome-headless-shell.exe',
        'msedge.exe',
        'Google Chrome.exe',
        'browser\chrome.exe',
        'browser\chrome-headless-shell.exe',
        'chromium\chrome.exe'
    )
    foreach ($name in $names) {
        if (Test-Path -LiteralPath (Join-Path $Root $name)) {
            return $true
        }
    }
    $chromeDir = Join-Path $Root 'chrome'
    if (Test-Path -LiteralPath $chromeDir) {
        $found = Get-ChildItem -LiteralPath $chromeDir -Recurse -Filter 'chrome.exe' -ErrorAction SilentlyContinue |
            Select-Object -First 1
        if ($found) {
            return $true
        }
    }
    return $false
}

if (Test-BundledChrome -Root $Dir) {
    Write-Host "Chrome already present in $Dir; skipping @puppeteer/browsers install."
    exit 0
}

Write-Host "Installing Chrome stable into $Dir ..."
Push-Location -LiteralPath $Dir
try {
    & npx --yes @puppeteer/browsers install chrome@stable --path .
    if ($LASTEXITCODE -ne 0) {
        throw "npx @puppeteer/browsers install failed with exit code $LASTEXITCODE"
    }
    Write-Host 'Chrome installed.'
}
finally {
    Pop-Location
}
