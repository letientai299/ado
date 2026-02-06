#Requires -Version 5.1
<#
.SYNOPSIS
    Download and install ado from GitHub.

.DESCRIPTION
    This script downloads the latest release or main branch build of ado
    and installs it to the specified directory.

.PARAMETER FromMain
    Download the latest successful build from the main branch instead of
    the latest release. Requires GitHub CLI (gh) to be installed and authenticated.

.PARAMETER InstallDir
    Installation directory. Defaults to ~/bin (or ~\bin on Windows).

.EXAMPLE
    .\install.ps1
    Install the latest release to ~/bin

.EXAMPLE
    .\install.ps1 -FromMain
    Install the latest build from main branch

.EXAMPLE
    .\install.ps1 -InstallDir "C:\Tools"
    Install to a custom directory
#>

[CmdletBinding()]
param(
    [Alias("m")]
    [switch]$FromMain,

    [Alias("d")]
    [string]$InstallDir
)

$ErrorActionPreference = "Stop"

$Repo = "letientai299/ado"
$BinaryName = "ado"

function Get-Platform {
    # Detect OS - handle both PowerShell 5.x (Windows only) and 6.x+ (cross-platform)
    $os = $null
    $arch = $null

    # PowerShell 5.x is Windows-only and doesn't have $IsWindows
    if ($PSVersionTable.PSEdition -eq "Desktop") {
        $os = "windows"
    }
    elseif ($IsWindows) {
        $os = "windows"
    }
    elseif ($IsLinux) {
        $os = "linux"
    }
    elseif ($IsMacOS) {
        $os = "darwin"
    }
    else {
        throw "Unsupported operating system"
    }

    # Detect architecture
    if ($os -eq "windows") {
        # Use environment variable for Windows (works on 5.x and 6.x+)
        switch ($env:PROCESSOR_ARCHITECTURE) {
            "AMD64" { $arch = "amd64" }
            "ARM64" { throw "ARM64 Windows is not currently supported" }
            default { throw "Unsupported architecture: $env:PROCESSOR_ARCHITECTURE" }
        }
    }
    else {
        # For non-Windows (PowerShell 6.x+ only), use RuntimeInformation
        $osArch = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture
        switch ($osArch) {
            "X64" { $arch = "amd64" }
            "Arm64" {
                if ($os -eq "darwin") {
                    $arch = "arm64"
                }
                else {
                    throw "ARM64 is only supported on macOS"
                }
            }
            default { throw "Unsupported architecture: $osArch" }
        }
    }

    return "$os-$arch"
}

$Platform = Get-Platform
$IsWindowsPlatform = $Platform -like "windows-*"
$ArtifactName = "$BinaryName-$Platform"
$BinaryExt = if ($IsWindowsPlatform) { ".exe" } else { "" }

# Set default install dir based on platform
if (-not $InstallDir) {
    $InstallDir = if ($IsWindowsPlatform) {
        "$env:USERPROFILE\bin"
    }
    else {
        "$env:HOME/.local/bin"
    }
}

function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor Green
}

function Write-Warn {
    param([string]$Message)
    Write-Host "[WARN] $Message" -ForegroundColor Yellow
}

function Write-Err {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
    exit 1
}

function Get-LatestRelease {
    Write-Info "Fetching latest release..."

    $releaseUrl = "https://api.github.com/repos/$Repo/releases/latest"

    # Use different parameters based on PowerShell version for compatibility
    $params = @{
        Uri    = $releaseUrl
        Method = "Get"
    }

    # UseBasicParsing is deprecated in 6.x+ but required for 5.x without IE
    if ($PSVersionTable.PSVersion.Major -le 5) {
        $params.UseBasicParsing = $true
    }

    $release = Invoke-RestMethod @params

    # Find the artifact for our platform
    $assetName = "$ArtifactName$BinaryExt"
    $asset = $release.assets | Where-Object { $_.name -eq $assetName } | Select-Object -First 1

    if (-not $asset) {
        Write-Err "Could not find release artifact for platform: $Platform"
    }

    $downloadUrl = $asset.browser_download_url
    $tempFile = Join-Path ([System.IO.Path]::GetTempPath()) $assetName

    Write-Info "Downloading $assetName..."

    $webParams = @{
        Uri     = $downloadUrl
        OutFile = $tempFile
    }

    if ($PSVersionTable.PSVersion.Major -le 5) {
        $webParams.UseBasicParsing = $true
    }

    Invoke-WebRequest @webParams

    return $tempFile
}

function Get-FromMain {
    Write-Info "Fetching latest successful build from main branch..."

    # Check if gh CLI is available
    $ghPath = Get-Command gh -ErrorAction SilentlyContinue
    if (-not $ghPath) {
        Write-Err "GitHub CLI (gh) is required for -FromMain. Install from: https://cli.github.com/"
    }

    # Check if authenticated
    $null = gh auth status 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Err "GitHub CLI not authenticated. Run: gh auth login"
    }

    # Get latest successful workflow run
    $runJson = gh run list --repo $Repo --branch main --workflow ci.yml --status success --limit 1 --json databaseId 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Err "Failed to fetch workflow runs: $runJson"
    }

    $runs = $runJson | ConvertFrom-Json
    if (-not $runs -or $runs.Count -eq 0) {
        Write-Err "No successful workflow runs found on main branch"
    }

    $runId = $runs[0].databaseId
    Write-Info "Found workflow run: $runId"

    $tempDir = Join-Path ([System.IO.Path]::GetTempPath()) "ado-download-$PID"
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null

    Write-Info "Downloading artifact: $ArtifactName..."
    gh run download $runId --repo $Repo --name $ArtifactName --dir $tempDir

    if ($LASTEXITCODE -ne 0) {
        Write-Err "Failed to download artifact"
    }

    # Find the downloaded binary
    $binaryPattern = "$BinaryName*"
    $binary = Get-ChildItem -Path $tempDir -Filter $binaryPattern -Recurse | Select-Object -First 1

    if (-not $binary) {
        Write-Err "Could not find downloaded binary in $tempDir"
    }

    return $binary.FullName
}

function Install-Binary {
    param(
        [string]$BinaryPath,
        [string]$DestDir
    )

    # Create install directory if needed
    if (-not (Test-Path $DestDir)) {
        Write-Info "Creating directory: $DestDir"
        New-Item -ItemType Directory -Path $DestDir -Force | Out-Null
    }

    $destPath = Join-Path $DestDir "$BinaryName$BinaryExt"
    Write-Info "Installing to $destPath..."
    Copy-Item -Path $BinaryPath -Destination $destPath -Force

    # Make executable on Unix (PowerShell 6.x+ only)
    if (-not $IsWindowsPlatform -and (Get-Command chmod -ErrorAction SilentlyContinue)) {
        chmod +x $destPath
    }

    # Cleanup
    $parentDir = Split-Path -Parent $BinaryPath
    if ($parentDir -like "*ado-download-*") {
        Remove-Item -Path $parentDir -Recurse -Force -ErrorAction SilentlyContinue
    }
    else {
        Remove-Item -Path $BinaryPath -Force -ErrorAction SilentlyContinue
    }

    Write-Info "Successfully installed $BinaryName to $destPath"

    # Check if directory is in PATH
    $pathSeparator = if ($IsWindowsPlatform) { ";" } else { ":" }
    $pathDirs = $env:PATH -split [regex]::Escape($pathSeparator)
    if ($pathDirs -notcontains $DestDir) {
        Write-Warn "$DestDir is not in your PATH"
        if ($IsWindowsPlatform) {
            Write-Warn "Add it with: `$env:PATH += `";$DestDir`""
            Write-Warn "Or permanently: [Environment]::SetEnvironmentVariable('PATH', `$env:PATH + ';$DestDir', 'User')"
        }
        else {
            Write-Warn "Add it with: export PATH=`"`$PATH:$DestDir`""
        }
    }
}

# Main
try {
    Write-Info "Platform: $Platform"

    $binary = if ($FromMain) {
        Get-FromMain
    }
    else {
        Get-LatestRelease
    }

    Install-Binary -BinaryPath $binary -DestDir $InstallDir
}
catch {
    Write-Err $_.Exception.Message
}
