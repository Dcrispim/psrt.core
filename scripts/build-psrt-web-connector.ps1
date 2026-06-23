# Builds psrt-web-connector.exe without a console window on double-click.
# Logs appear only when launched from an existing terminal (cmd/PowerShell).
$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
$icons = Join-Path $root "internal\webconnector\icons"
$connector = Join-Path $root "cmd\psrt-web-connector"

Push-Location $root
try {
    if (-not (Test-Path (Join-Path $icons "logo1-16.png"))) {
        throw "Missing internal/webconnector/icons/logo1-16.png"
    }

    Add-Type -AssemblyName System.Drawing
    $bmp = [System.Drawing.Bitmap]::FromFile((Join-Path $icons "logo1-16.png"))
    $icon = [System.Drawing.Icon]::FromHandle($bmp.GetHicon())
    $icoPath = Join-Path $icons "tray.ico"
    $fs = [System.IO.File]::Create($icoPath)
    $icon.Save($fs)
    $fs.Close()
    $bmp.Dispose()

    if (-not (Get-Command go-winres -ErrorAction SilentlyContinue)) {
        go install github.com/tc-hib/go-winres@v0.3.3
    }
    Push-Location $connector
    go-winres make --in winres/winres.json --out resource.syso --arch amd64
    Pop-Location

    go build -ldflags="-H windowsgui" -o psrt-web-connector.exe ./cmd/psrt-web-connector
    Write-Host "Built: $root\psrt-web-connector.exe"
} finally {
    Pop-Location
}
