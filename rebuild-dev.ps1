# TramaTex - Reconstruccion total local (DB + API + Front)
# Ejecutar desde la raiz del proyecto: .\rebuild-dev.ps1

param(
    [switch]$NoFrontend, # Si se indica, reconstruye solo DB + API
    [switch]$PreserveDatabase,
    [switch]$Help
)

if ($Help) {
    Write-Host "Uso: .\rebuild-dev.ps1 [-NoFrontend] [-PreserveDatabase]" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Este script realiza una reconstruccion total usando solo scripts del proyecto:" -ForegroundColor Yellow
    Write-Host "  1. .\stop-dev.ps1 con las opciones de limpieza correspondientes" -ForegroundColor White
    Write-Host "  2. .\start-dev.ps1 con -RebuildImages (y -ResetData si se resetea BBDD)" -ForegroundColor White
    Write-Host ""
    Write-Host "Opciones:" -ForegroundColor Yellow
    Write-Host "  -NoFrontend    Reconstruye solo DB + API" -ForegroundColor White
    Write-Host "  -PreserveDatabase  Reconstruye aplicacion sin tocar volumenes/datos de BBDD" -ForegroundColor White
    Write-Host ""
    Write-Host "Ejemplos:" -ForegroundColor Yellow
    Write-Host "  .\rebuild-dev.ps1" -ForegroundColor Gray
    Write-Host "  .\rebuild-dev.ps1 -NoFrontend" -ForegroundColor Gray
    Write-Host "  .\rebuild-dev.ps1 -PreserveDatabase" -ForegroundColor Gray
    exit 0
}

$rootDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $rootDir

Write-Host "Iniciando reconstruccion total de TramaTex..." -ForegroundColor Cyan
if ($NoFrontend) {
    Write-Host "  Modo: DB + API" -ForegroundColor Magenta
} else {
    Write-Host "  Modo: FULL (DB + API + Frontend/Nginx)" -ForegroundColor Magenta
}
Write-Host ""

# 1) Parar y limpiar completamente
if ($PreserveDatabase) {
    Write-Host "[1/2] Limpieza de aplicacion preservando BBDD..." -ForegroundColor Yellow
    powershell -ExecutionPolicy Bypass -File .\stop-dev.ps1 -RemoveImages
} else {
    Write-Host "[1/2] Limpieza completa de stack, datos e imagen de PostgreSQL..." -ForegroundColor Yellow
    powershell -ExecutionPolicy Bypass -File .\stop-dev.ps1 -FullCleanup -RemoveImages -RemoveDbImage
}
if ($LASTEXITCODE -ne 0) {
    Write-Host "Error en limpieza completa" -ForegroundColor Red
    exit 1
}

# 2) Reconstruir imagenes y levantar
Write-Host "[2/2] Reconstruyendo imagenes y levantando servicios..." -ForegroundColor Yellow
if ($NoFrontend) {
    if ($PreserveDatabase) {
        powershell -ExecutionPolicy Bypass -File .\start-dev.ps1 -RebuildImages
    } else {
        powershell -ExecutionPolicy Bypass -File .\start-dev.ps1 -ResetData -RebuildImages
    }
} else {
    if ($PreserveDatabase) {
        powershell -ExecutionPolicy Bypass -File .\start-dev.ps1 -Full -RebuildImages
    } else {
        powershell -ExecutionPolicy Bypass -File .\start-dev.ps1 -Full -ResetData -RebuildImages
    }
}

if ($LASTEXITCODE -ne 0) {
    Write-Host "Error al reconstruir/levantar servicios" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "Reconstruccion total completada" -ForegroundColor Green
