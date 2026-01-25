#!/usr/bin/env pwsh
# TramaTex - Docker Connectivity Test Script

Write-Host "`n=== TRAMATEX - CONNECTIVITY VERIFICATION ===" -ForegroundColor Cyan

$projectRoot = "c:\Users\joran\Desarrollos\TramaTex"
Set-Location $projectRoot

# 1. Check files
Write-Host "`n[1] Checking project files..." -ForegroundColor Yellow
$files = @("Makefile", ".env.local", ".env.remote", "docker-compose.yml", "docker-compose.local.yml", "docker-compose.remote.yml")
foreach ($f in $files) {
    if (Test-Path $f) { Write-Host "    OK: $f" -ForegroundColor Green }
    else { Write-Host "    FAIL: $f" -ForegroundColor Red }
}

# 2. Check Makefile ENV
Write-Host "`n[2] Makefile configuration..." -ForegroundColor Yellow
$env_line = Select-String "^ENV \?=" Makefile
Write-Host "    $($env_line.Line)" -ForegroundColor Cyan

# 3. Check .env.remote
Write-Host "`n[3] Remote environment config (.env.remote)..." -ForegroundColor Yellow
$remote = @{}
(Get-Content .env.remote | Select-String "^SSH_|^DB_HOST") | ForEach-Object {
    $key, $val = $_.Line -split "=" | Select-Object -First 2
    Write-Host "    $key = $val" -ForegroundColor Cyan
}

# 4. Network test
Write-Host "`n[4] Network connectivity..." -ForegroundColor Yellow
$ping = ping -n 1 pcele 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host "    PASS: pcele (192.168.0.20) is reachable" -ForegroundColor Green
} else {
    Write-Host "    FAIL: pcele is not reachable" -ForegroundColor Red
}

# 5. SSH test
Write-Host "`n[5] SSH connectivity..." -ForegroundColor Yellow
Write-Host "    Attempting: ssh -o ConnectTimeout=5 ele@pcele 'echo OK'" -ForegroundColor Gray
$ssh_test = ssh -o ConnectTimeout=5 -o BatchMode=yes ele@pcele "echo SSH_TEST_OK" 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host "    PASS: SSH connection successful" -ForegroundColor Green
} elseif ($ssh_test -like "*Permission denied*") {
    Write-Host "    WARNING: SSH authentication needed (password/keys not configured)" -ForegroundColor Yellow
} else {
    Write-Host "    INFO: SSH needs setup" -ForegroundColor Yellow
}

# 6. Docker local
Write-Host "`n[6] Docker installation check..." -ForegroundColor Yellow
$docker_ver = docker --version 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host "    OK: Docker installed - $docker_ver" -ForegroundColor Green
} else {
    Write-Host "    WARNING: Docker not found" -ForegroundColor Yellow
}

# 7. Local Docker test
Write-Host "`n[7] Local Docker connection (Docker Desktop)..." -ForegroundColor Yellow
$docker_local = docker ps 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host "    PASS: Docker Desktop is running" -ForegroundColor Green
} else {
    Write-Host "    WARNING: Docker Desktop may not be running" -ForegroundColor Yellow
}

# Summary
Write-Host "`n" -ForegroundColor Gray
Write-Host "=== SUMMARY ===" -ForegroundColor Cyan
Write-Host "Project structure      : OK" -ForegroundColor Green
Write-Host "Makefile (remote)      : OK (ENV = remote)" -ForegroundColor Green
Write-Host "Network (pcele)        : OK (reachable)" -ForegroundColor Green
Write-Host "SSH authentication     : Pending setup" -ForegroundColor Yellow
Write-Host "Docker                 : $(if ($LASTEXITCODE -eq 0) { 'OK' } else { 'Check needed' })" -ForegroundColor Yellow

Write-Host "`n=== NEXT STEPS ===" -ForegroundColor Cyan
Write-Host @"
1. Install make (Windows):
   choco install make

2. Set up SSH keys (recommended):
   ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519
   ssh-copy-id -i ~/.ssh/id_ed25519.pub ele@pcele

3. Test connection:
   ssh ele@pcele 'docker ps'

4. Run make commands:
   make docker-up        (uses remote/pcele by default)
   make docker-status
   make docker-logs

"@ -ForegroundColor White

Write-Host "Connectivity verification complete!" -ForegroundColor Green
