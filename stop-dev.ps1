# TramaTex - Script para Detener el Ambiente de Desarrollo
# Ejecutar desde la raíz del proyecto: .\stop-dev.ps1

Write-Host "Deteniendo TramaTex Development Environment..." -ForegroundColor Cyan
Write-Host ""

# Navegar al directorio raíz si no estamos ahí
$rootDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $rootDir

# Detener servicios (incluyendo profile full si estuviera activo)
Write-Host "Deteniendo contenedores..." -ForegroundColor Yellow
docker compose -f docker/docker-compose.local.yml --env-file docker/.env --profile full down

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "Servicios detenidos correctamente" -ForegroundColor Green
} else {
    Write-Host ""
    Write-Host "Algunos contenedores pueden no haberse detenido correctamente" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Para eliminar tambien los volumenes (datos), usa:" -ForegroundColor Cyan
Write-Host "   docker compose -f docker/docker-compose.local.yml --env-file docker/.env --profile full down -v" -ForegroundColor Gray
Write-Host ""
