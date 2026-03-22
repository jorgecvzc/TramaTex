# TramaTex - Script para levantar el ambiente de desarrollo
# Ejecutar desde la raiz del proyecto: .\start-dev.ps1
# Con frontend (Nginx): .\start-dev.ps1 -Full

param(
    [switch]$Full  # Incluir frontend/Nginx (profile: full)
)

Write-Host "Iniciando TramaTex Development Environment..." -ForegroundColor Cyan
if ($Full) { Write-Host "  Modo: FULL (DB + API + Frontend/Nginx)" -ForegroundColor Magenta }
else       { Write-Host "  Modo: DEV  (DB + API — frontend via npm run dev)" -ForegroundColor Magenta }
Write-Host ""

$previousErrorActionPreference = $ErrorActionPreference
$ErrorActionPreference = 'Continue'

$rootDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $rootDir

Write-Host "Verificando Docker..." -ForegroundColor Yellow
try {
    docker ps > $null 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Docker no esta corriendo. Inicia Docker Desktop." -ForegroundColor Red
        exit 1
    }
    Write-Host "Docker esta corriendo" -ForegroundColor Green
}
catch {
    Write-Host "Error al verificar Docker: $_" -ForegroundColor Red
    exit 1
}

$profileArgs = @()
if ($Full) { $profileArgs = @("--profile", "full") }

Write-Host ""
Write-Host "Limpiando contenedores previos..." -ForegroundColor Yellow
try {
    docker compose -f docker/docker-compose.local.yml --env-file docker/.env @profileArgs down -v *> $null
}
catch {
    # Ignorar errores en limpieza previa
}

Write-Host ""
Write-Host "Construyendo y levantando servicios..." -ForegroundColor Yellow
docker compose -f docker/docker-compose.local.yml --env-file docker/.env @profileArgs up -d --build

if ($LASTEXITCODE -ne 0) {
    Write-Host ""
    Write-Host "Error al iniciar los servicios" -ForegroundColor Red
    Write-Host "Verifica logs con: docker compose -f docker/docker-compose.local.yml --env-file docker/.env logs" -ForegroundColor Gray
    exit 1
}

Write-Host ""
Write-Host "Servicios iniciados correctamente" -ForegroundColor Green
Write-Host ""
Write-Host "Estado de contenedores:" -ForegroundColor Cyan
Start-Sleep -Seconds 2
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

Write-Host ""
Write-Host "URLs disponibles:" -ForegroundColor Cyan
Write-Host "  API Health: http://localhost:8080/api/health" -ForegroundColor White
Write-Host "  Database:   localhost:5432 (user: tramatex, db: tramatex)" -ForegroundColor White
if ($Full) {
    Write-Host "  Frontend:   http://localhost:3000 (Nginx)" -ForegroundColor White
} else {
    Write-Host "  Frontend:   Ejecuta 'cd apps/frontend && npm run dev' (puerto 5173)" -ForegroundColor White
}

Write-Host ""
Write-Host "Esperando que la API este lista..." -ForegroundColor Yellow
$maxAttempts = 30
$attempt = 0
$apiReady = $false

while ($attempt -lt $maxAttempts) {
    try {
        $response = Invoke-WebRequest -Uri "http://localhost:8080/api/health" -UseBasicParsing -ErrorAction SilentlyContinue
        if ($response.StatusCode -eq 200) {
            $apiReady = $true
            break
        }
    }
    catch {
    }

    $attempt++
    Start-Sleep -Seconds 2
    Write-Host "." -NoNewline -ForegroundColor Gray
}

Write-Host ""
if ($apiReady) {
    Write-Host "API esta lista y respondiendo" -ForegroundColor Green
}
else {
    Write-Host "La API tardo mas de lo esperado. Verifica logs con: docker logs tramatex_api" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Ambiente de desarrollo listo" -ForegroundColor Green
Write-Host ""

$ErrorActionPreference = $previousErrorActionPreference
