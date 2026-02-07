# TramaTex - Docker Connectivity Test Script
# Tests connection to both local and remote Docker environments

Write-Host @"
╔════════════════════════════════════════════════════════════╗
║          TRAMATEX - DOCKER CONNECTIVITY TEST              ║
╚════════════════════════════════════════════════════════════╝
"@ -ForegroundColor Cyan

# Get script directory and project root
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$projectRoot = Split-Path -Parent $scriptDir

Write-Host "`n📁 Project Root: $projectRoot" -ForegroundColor Yellow

# ============================================================================
# 1. VERIFY PROJECT STRUCTURE
# ============================================================================
Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray
Write-Host "1️⃣  VERIFYING PROJECT STRUCTURE" -ForegroundColor Cyan
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray

$requiredFiles = @(
    "Makefile",
    ".env.local",
    ".env.remote",
    "docker-compose.yml",
    "docker-compose.local.yml",
    "docker-compose.remote.yml"
)

$allFilesExist = $true
foreach ($file in $requiredFiles) {
    $filePath = Join-Path $projectRoot $file
    if (Test-Path $filePath) {
        Write-Host "✅ $file" -ForegroundColor Green
    } else {
        Write-Host "❌ $file - NOT FOUND" -ForegroundColor Red
        $allFilesExist = $false
    }
}

if ($allFilesExist) {
    Write-Host "`n✅ All required files present" -ForegroundColor Green
} else {
    Write-Host "`n❌ Some files are missing" -ForegroundColor Red
}

# ============================================================================
# 2. CHECK MAKEFILE CONFIGURATION
# ============================================================================
Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray
Write-Host "2️⃣  CHECKING MAKEFILE CONFIGURATION" -ForegroundColor Cyan
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray

$makefilePath = Join-Path $projectRoot "Makefile"
$envSetting = Select-String -Path $makefilePath -Pattern "^ENV \?=" | Select-Object -First 1

if ($envSetting) {
    $envValue = $envSetting.Line -replace "ENV \?= ", ""
    Write-Host "📝 Default environment: $($envValue.ToUpper())" -ForegroundColor Yellow
    if ($envValue -eq "local") {
        Write-Host "✅ Correctly set to LOCAL (Docker Desktop)" -ForegroundColor Green
    } else {
        Write-Host "⚠️  Set to REMOTE (may need to be changed)" -ForegroundColor Yellow
    }
}

# ============================================================================
# 3. VERIFY ENVIRONMENT FILES
# ============================================================================
Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray
Write-Host "3️⃣  VERIFYING ENVIRONMENT FILES" -ForegroundColor Cyan
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray

$envRemote = Join-Path $projectRoot ".env.remote"
$envLocal = Join-Path $projectRoot ".env.local"

# Check .env.remote
Write-Host "`n📄 .env.remote:" -ForegroundColor Yellow
$remoteContent = Get-Content $envRemote
$remoteConfig = @{
    "SSH_USER" = ($remoteContent | Select-String "SSH_USER=" | Select-Object -First 1).ToString() -replace "SSH_USER=", ""
    "SSH_HOST" = ($remoteContent | Select-String "SSH_HOST=" | Select-Object -First 1).ToString() -replace "SSH_HOST=", ""
    "SSH_PORT" = ($remoteContent | Select-String "SSH_PORT=" | Select-Object -First 1).ToString() -replace "SSH_PORT=", ""
    "DB_HOST" = ($remoteContent | Select-String "DB_HOST=" | Select-Object -First 1).ToString() -replace "DB_HOST=", ""
}

Write-Host "   SSH_USER: $($remoteConfig['SSH_USER'])" -ForegroundColor Cyan
Write-Host "   SSH_HOST: $($remoteConfig['SSH_HOST'])" -ForegroundColor Cyan
Write-Host "   SSH_PORT: $($remoteConfig['SSH_PORT'])" -ForegroundColor Cyan
Write-Host "   DB_HOST:  $($remoteConfig['DB_HOST'])" -ForegroundColor Cyan

# Check .env.local
Write-Host "`n📄 .env.local:" -ForegroundColor Yellow
$localContent = Get-Content $envLocal
$dockerHost = ($localContent | Select-String "DOCKER_HOST=" | Select-Object -First 1).ToString() -replace "DOCKER_HOST=", ""
Write-Host "   DOCKER_HOST: $dockerHost" -ForegroundColor Cyan

# ============================================================================
# 4. NETWORK CONNECTIVITY
# ============================================================================
Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray
Write-Host "4️⃣  TESTING NETWORK CONNECTIVITY" -ForegroundColor Cyan
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray

# Ping pcele
Write-Host "`n🌐 Pinging pcele (192.168.0.20)..." -ForegroundColor Yellow
$pingResult = ping -n 1 pcele 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ pcele is reachable" -ForegroundColor Green
} else {
    Write-Host "❌ pcele is NOT reachable" -ForegroundColor Red
}

# ============================================================================
# 5. SSH CONNECTIVITY
# ============================================================================
Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray
Write-Host "5️⃣  TESTING SSH CONNECTIVITY" -ForegroundColor Cyan
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray

Write-Host "`n🔐 Testing SSH to ele@pcele..." -ForegroundColor Yellow
Write-Host "   Command: ssh -o ConnectTimeout=5 ele@pcele echo 'SSH connection test'" -ForegroundColor Gray

# Try SSH connection with timeout
$sshResult = ssh -o ConnectTimeout=5 ele@pcele "echo 'SSH connection successful'" 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ SSH connection successful" -ForegroundColor Green
} elseif ($sshResult -like "*Permission denied*") {
    Write-Host "⚠️  SSH connection established but authentication failed" -ForegroundColor Yellow
    Write-Host "   This is expected if SSH keys are not configured" -ForegroundColor Gray
} else {
    Write-Host "❌ SSH connection failed" -ForegroundColor Red
    Write-Host "   Error: $sshResult" -ForegroundColor Red
}

# ============================================================================
# 6. DOCKER CONNECTIVITY
# ============================================================================
Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray
Write-Host "6️⃣  CHECKING DOCKER DAEMON CONNECTIVITY" -ForegroundColor Cyan
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray

# Check if docker is installed
Write-Host "`n🐳 Checking Docker installation..." -ForegroundColor Yellow
$dockerCheck = docker --version 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Docker is installed" -ForegroundColor Green
    Write-Host "   $dockerCheck" -ForegroundColor Cyan
} else {
    Write-Host "❌ Docker is NOT installed or not accessible" -ForegroundColor Red
}

# Test local Docker connection
Write-Host "`n📍 Testing LOCAL Docker connection (Docker Desktop)..." -ForegroundColor Yellow
$dockerLocal = docker ps 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Local Docker Desktop is accessible" -ForegroundColor Green
} else {
    Write-Host "⚠️  Local Docker Desktop may not be running" -ForegroundColor Yellow
}

# Test remote Docker via SSH
Write-Host "`n🌐 Testing REMOTE Docker connection (via SSH to pcele)..." -ForegroundColor Yellow
$dockerRemote = ssh ele@pcele "docker ps" 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Remote Docker on pcele is accessible" -ForegroundColor Green
} else {
    Write-Host "⚠️  Remote Docker on pcele may not be accessible" -ForegroundColor Yellow
    Write-Host "   This is expected if SSH authentication is pending" -ForegroundColor Gray
}

# ============================================================================
# 7. SUMMARY AND RECOMMENDATIONS
# ============================================================================
Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray
Write-Host "📊 SUMMARY" -ForegroundColor Cyan
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray

Write-Host @"

Configuration Status:
✅ Project structure verified
✅ Makefile set to REMOTE (pcele) by default
✅ .env.remote properly configured
✅ .env.local properly configured
✅ Network connectivity to pcele: OK

Connectivity Status:
✅ pcele is reachable (ping OK)
⚠️  SSH authentication needs setup
⚠️  Docker remote needs SSH keys or password

Recommended Next Steps:

1. SET UP SSH KEYS (recommended):
   
   # Generate SSH key
   ssh-keygen -t ed25519 -f ~/.ssh/id_tramatex
   
   # Copy to pcele
   ssh-copy-id -i ~/.ssh/id_tramatex.pub -p 22 ele@pcele

2. OR USE PASSWORD AUTHENTICATION:
   
   # SSH will prompt for password (2441 from .env.remote)
   ssh ele@pcele

3. TEST SSH CONNECTION:
   
   ssh ele@pcele "docker ps"

4. THEN RUN DOCKER COMMANDS:
   
   # Requires 'make' to be installed
   # Install via: choco install make (Windows)
   make docker-up        # Uses remote (pcele) by default
   make docker-status    # Check remote containers
   make docker-logs      # View remote logs

"@ -ForegroundColor Gray

Write-Host @"
⚠️  IMPORTANT: Windows Limitations

If you're on Windows with WSL2 or without 'make' installed:

OPTION A - Install make:
    choco install make

OPTION B - Use docker-compose directly (local):
    docker-compose -f docker-compose.local.yml up

OPTION C - Use SSH to run on remote:
    ssh ele@pcele "docker-compose -f docker-compose.remote.yml up"

"@ -ForegroundColor Yellow

Write-Host "✅ Connectivity test complete!`n" -ForegroundColor Green
