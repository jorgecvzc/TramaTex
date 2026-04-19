# TramaTex - Script para Detener el Ambiente de Desarrollo
# Ejecutar desde la raíz del proyecto: .\stop-dev.ps1

param(
    [switch]$FullCleanup,  # Incluye -v y --remove-orphans
    [switch]$RemoveImages,  # Incluye --rmi local
    [switch]$RemoveDbImage,
    [switch]$Help
)

if ($Help) {
    Write-Host "Uso: .\stop-dev.ps1 [-FullCleanup] [-RemoveImages] [-RemoveDbImage]" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Opciones:" -ForegroundColor Yellow
    Write-Host "  -FullCleanup   Ejecuta down -v --remove-orphans" -ForegroundColor White
    Write-Host "  -RemoveImages  Elimina imagenes locales del proyecto con --rmi local" -ForegroundColor White
    Write-Host "  -RemoveDbImage Elimina tambien la imagen de PostgreSQL definida en docker/.env" -ForegroundColor White
    Write-Host ""
    Write-Host "Ejemplos:" -ForegroundColor Yellow
    Write-Host "  .\stop-dev.ps1" -ForegroundColor Gray
    Write-Host "  .\stop-dev.ps1 -FullCleanup" -ForegroundColor Gray
    Write-Host "  .\stop-dev.ps1 -FullCleanup -RemoveImages" -ForegroundColor Gray
    Write-Host "  .\stop-dev.ps1 -FullCleanup -RemoveImages -RemoveDbImage" -ForegroundColor Gray
    exit 0
}

Write-Host "Deteniendo TramaTex Development Environment..." -ForegroundColor Cyan
Write-Host ""

# Navegar al directorio raíz si no estamos ahí
$rootDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $rootDir

# Detener servicios (incluyendo profile full si estuviera activo)
Write-Host "Deteniendo contenedores..." -ForegroundColor Yellow
$composeArgs = @(
    "compose",
    "-f", "docker/docker-compose.local.yml",
    "--env-file", "docker/.env",
    "--profile", "full",
    "down"
)

if ($FullCleanup) {
    Write-Host "Aplicando limpieza completa (-v --remove-orphans)..." -ForegroundColor Yellow
    $composeArgs += @("-v", "--remove-orphans")
}

if ($RemoveImages) {
    Write-Host "Eliminando imagenes locales del proyecto (--rmi local)..." -ForegroundColor Yellow
    $composeArgs += @("--rmi", "local")
}

docker @composeArgs

if ($RemoveDbImage) {
    try {
        $dbImage = (Get-Content "docker/.env" | Select-String '^DB_IMAGE=' | Select-Object -First 1).ToString().Split('=')[1].Trim()
        if ($dbImage) {
            Write-Host "Eliminando imagen de base de datos: $dbImage" -ForegroundColor Yellow
            docker image rm -f $dbImage *> $null
        }
    }
    catch {
        Write-Host "No se pudo eliminar la imagen de PostgreSQL automaticamente: $_" -ForegroundColor Yellow
    }
}

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "Servicios detenidos correctamente" -ForegroundColor Green
} else {
    Write-Host ""
    Write-Host "Algunos contenedores pueden no haberse detenido correctamente" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Para eliminar tambien los volumenes (datos), usa:" -ForegroundColor Cyan
Write-Host "   .\stop-dev.ps1 -FullCleanup" -ForegroundColor Gray
Write-Host "   .\stop-dev.ps1 -FullCleanup -RemoveImages" -ForegroundColor Gray
Write-Host "   .\stop-dev.ps1 -FullCleanup -RemoveImages -RemoveDbImage" -ForegroundColor Gray
Write-Host "   # equivalente manual:" -ForegroundColor DarkGray
Write-Host "   docker compose -f docker/docker-compose.local.yml --env-file docker/.env --profile full down -v --remove-orphans" -ForegroundColor DarkGray
Write-Host "   docker compose -f docker/docker-compose.local.yml --env-file docker/.env --profile full down -v --remove-orphans --rmi local" -ForegroundColor DarkGray
Write-Host ""
