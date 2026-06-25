#Requires -Version 5.1
# ============================================================
#  bgscan installer (Windows)
#  https://github.com/MohsenBg/bgscan
# ============================================================
Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$REPO  = "MohsenBg/bgscan"
$API   = "https://api.github.com/repos/$REPO/releases/latest"

# ── colour helpers ───────────────────────────────────────────
function Write-Banner {
    Write-Host ""
    Write-Host "  ██████╗  ██████╗ ███████╗ ██████╗ █████╗ ███╗   ██╗" -ForegroundColor Cyan
    Write-Host "  ██╔══██╗██╔════╝ ██╔════╝██╔════╝██╔══██╗████╗  ██║" -ForegroundColor Cyan
    Write-Host "  ██████╔╝██║  ███╗███████╗██║     ███████║██╔██╗ ██║" -ForegroundColor Cyan
    Write-Host "  ██╔══██╗██║   ██║╚════██║██║     ██╔══██║██║╚██╗██║" -ForegroundColor Cyan
    Write-Host "  ██████╔╝╚██████╔╝███████║╚██████╗██║  ██║██║ ╚████║" -ForegroundColor Cyan
    Write-Host "  ╚═════╝  ╚═════╝ ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═══╝" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "  Installer  •  github.com/$REPO" -ForegroundColor DarkGray
    Write-Host "  ────────────────────────────────────────────────────" -ForegroundColor DarkGray
    Write-Host ""
}

function Write-Section([string]$msg) {
    Write-Host ""
    Write-Host "  $msg" -ForegroundColor Magenta
    Write-Host "  ────────────────────────────────────────────────" -ForegroundColor DarkGray
}

function Write-Step([string]$msg)  { Write-Host "  → " -ForegroundColor Blue   -NoNewline; Write-Host $msg -ForegroundColor White }
function Write-Ok([string]$msg)    { Write-Host "  ✔ " -ForegroundColor Green  -NoNewline; Write-Host $msg -ForegroundColor Green }
function Write-Warn([string]$msg)  { Write-Host "  ⚠ " -ForegroundColor Yellow -NoNewline; Write-Host $msg -ForegroundColor Yellow }
function Write-Info([string]$msg)  { Write-Host "     $msg"                    -ForegroundColor DarkGray }
function Write-Fatal([string]$msg) {
    Write-Host ""
    Write-Host "  ✖  Error: " -ForegroundColor Red -NoNewline
    Write-Host $msg            -ForegroundColor Red
    Write-Host ""
    exit 1
}

# ── detect platform ──────────────────────────────────────────
function Get-Asset {
    $arch = $env:PROCESSOR_ARCHITECTURE   # AMD64 | ARM64 | x86

    switch ($arch) {
        'AMD64' { return "bgscan-windows-64.zip"    }
        'ARM64' { return "bgscan-windows-arm64.zip" }
        'x86'   { return "bgscan-windows-32.zip"    }
        default { Write-Fatal "Unsupported architecture: $arch`n`n  Please open an issue: https://github.com/$REPO/issues" }
    }
}

# ── fetch latest release download URL ────────────────────────
function Get-DownloadUrl([string]$asset) {
    $headers = @{ 'User-Agent' = 'bgscan-installer' }
    try {
        $release = Invoke-RestMethod -Uri $API -Headers $headers
    } catch {
        Write-Fatal "Failed to contact GitHub API: $_"
    }

    $url = $release.assets |
        Where-Object { $_.name -eq $asset } |
        Select-Object -ExpandProperty browser_download_url -First 1

    return $url
}

# ── prompt for a numbered choice ─────────────────────────────
function Read-Choice([string]$question, [string[]]$options) {
    Write-Host ""
    Write-Host "  $question" -ForegroundColor White
    Write-Host ""
    for ($i = 0; $i -lt $options.Length; $i++) {
        Write-Host "    [$($i+1)]  $($options[$i])" -ForegroundColor White
    }
    Write-Host ""
    Write-Host "  Your choice: " -ForegroundColor White -NoNewline
    return (Read-Host)
}

# ── simple spinner (dots) ─────────────────────────────────────
function Invoke-WithSpinner([scriptblock]$job, [string]$label = "Working") {
    $frames = @('⠋','⠙','⠹','⠸','⠼','⠴','⠦','⠧','⠇','⠏')
    $i = 0

    $jobObj = Start-Job -ScriptBlock $job
    while ($jobObj.State -eq 'Running') {
        Write-Host "`r  $($frames[$i])  $label   " -NoNewline -ForegroundColor Cyan
        $i = ($i + 1) % $frames.Length
        Start-Sleep -Milliseconds 80
    }
    Write-Host "`r" + (' ' * 60) + "`r" -NoNewline

    $result = Receive-Job $jobObj -Wait -AutoRemoveJob
    if ($jobObj.State -eq 'Failed') {
        Write-Fatal "Background job failed: $($jobObj.ChildJobs[0].JobStateInfo.Reason)"
    }
    return $result
}

# ════════════════════════════════════════════════════════════
#  MAIN
# ════════════════════════════════════════════════════════════
Write-Banner

# ── 1. Detect platform ───────────────────────────────────────
Write-Section "Detecting environment"

$OS   = "windows"
$ARCH = $env:PROCESSOR_ARCHITECTURE
Write-Info "Operating system : Windows"
Write-Info "Architecture     : $ARCH"

$ASSET = Get-Asset
Write-Info "Release asset    : $ASSET"
Write-Ok   "Platform supported"

# ── 2. Resolve install location ──────────────────────────────
Write-Section "Install location"
$INSTALL_DIR = Join-Path $PWD "bgscan"
Write-Info "Target : $INSTALL_DIR"

# ── 3. Handle existing install ───────────────────────────────
if (Test-Path $INSTALL_DIR) {
    Write-Host ""
    Write-Warn "An existing installation was found at $INSTALL_DIR"

    $choice = Read-Choice `
        "How would you like to proceed?" `
        @(
            "Remove the old installation and install fresh",
            "Back up the old installation (rename to bgscan_old)",
            "Cancel — exit the installer"
        )

    switch ($choice) {
        '1' {
            Write-Step "Removing old installation..."
            Remove-Item -Recurse -Force $INSTALL_DIR
            Write-Ok "Old installation removed"
        }
        '2' {
            $backup = "${INSTALL_DIR}_old"
            Write-Step "Backing up to $backup..."
            if (Test-Path $backup) { Remove-Item -Recurse -Force $backup }
            Rename-Item $INSTALL_DIR $backup
            Write-Ok "Backup saved to $backup"
        }
        default {
            Write-Host ""
            Write-Info "Installation cancelled. No changes were made."
            Write-Host ""
            exit 0
        }
    }
}

New-Item -ItemType Directory -Force -Path $INSTALL_DIR | Out-Null

# ── 4. Dependencies — nothing extra needed on Windows ────────
Write-Section "Checking dependencies"
Write-Ok "No additional dependencies required (Expand-Archive built in)"

# ── 5. Fetch download URL ────────────────────────────────────
Write-Section "Fetching release information"
Write-Step "Querying GitHub API..."

$URL = Get-DownloadUrl $ASSET
if (-not $URL) {
    Write-Fatal "Could not find a download URL for `"$ASSET`".`n`n  Check releases: https://github.com/$REPO/releases"
}

Write-Ok "Release URL resolved"
Write-Info $URL

# ── 6. Download ──────────────────────────────────────────────
Write-Section "Downloading"
$TMP     = [System.IO.Path]::GetTempPath()
$TMPDIR  = Join-Path $TMP ([System.Guid]::NewGuid().ToString())
New-Item -ItemType Directory -Force -Path $TMPDIR | Out-Null
$ZIPFILE = Join-Path $TMPDIR "bgscan.zip"

Write-Step "Downloading $ASSET..."
try {
    # Use BITS for progress if available, fall back to Invoke-WebRequest
    if (Get-Command Start-BitsTransfer -ErrorAction SilentlyContinue) {
        Start-BitsTransfer -Source $URL -Destination $ZIPFILE -DisplayName "bgscan" -Description "Downloading $ASSET"
    } else {
        $ProgressPreference = 'SilentlyContinue'   # faster on older PS
        Invoke-WebRequest -Uri $URL -OutFile $ZIPFILE -UseBasicParsing
        $ProgressPreference = 'Continue'
    }
} catch {
    Write-Fatal "Download failed: $_"
} finally {
    # Clean up temp on exit
    Register-EngineEvent -SourceIdentifier PowerShell.Exiting -Action {
        if (Test-Path $TMPDIR) { Remove-Item -Recurse -Force $TMPDIR }
    } | Out-Null
}
Write-Ok "Download complete"

# ── 7. Extract & install ─────────────────────────────────────
Write-Section "Installing"

Write-Step "Extracting archive..."
$EXTRACT_DIR = Join-Path $TMPDIR "extracted"
Expand-Archive -Path $ZIPFILE -DestinationPath $EXTRACT_DIR -Force

# Find the single top-level folder inside the archive (mirrors bash behaviour)
$SRC_DIR = Get-ChildItem -Path $EXTRACT_DIR -Directory | Select-Object -First 1
if (-not $SRC_DIR) {
    # Archive may have dumped files directly — use the extract root
    $SRC_DIR = Get-Item $EXTRACT_DIR
}

Write-Step "Installing files..."
Remove-Item -Recurse -Force $INSTALL_DIR -ErrorAction SilentlyContinue
New-Item -ItemType Directory -Force -Path $INSTALL_DIR | Out-Null
Copy-Item -Path "$($SRC_DIR.FullName)\*" -Destination $INSTALL_DIR -Recurse -Force

# Mark all bgscan executables as executable (not strictly required on Windows,
# but unblocks files downloaded from the internet)
Get-ChildItem -Path $INSTALL_DIR -Recurse -File |
    Where-Object { $_.Name -like "bgscan*" } |
    ForEach-Object { Unblock-File -Path $_.FullName }

Write-Ok "Installed at $INSTALL_DIR"

# ── 8. Done ──────────────────────────────────────────────────
Write-Host ""
Write-Host "  ────────────────────────────────────────────────────" -ForegroundColor DarkGray
Write-Host ""
Write-Host "  ✔  bgscan installed successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "  Get started:" -ForegroundColor White
Write-Host "  cd bgscan"    -ForegroundColor Cyan
Write-Host "  .\bgscan.exe" -ForegroundColor Cyan
Write-Host ""
Write-Host "  Docs & source: https://github.com/$REPO" -ForegroundColor DarkGray
Write-Host "  ────────────────────────────────────────────────────" -ForegroundColor DarkGray
Write-Host ""
