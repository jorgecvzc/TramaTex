# TramaTex - Force production deployment via GitHub Actions
# Requires: GitHub CLI (gh) installed and authenticated
# Usage: .\scripts\deploy-production.ps1

param(
    [switch]$Watch,
    [switch]$FreshDB,
    [switch]$Help
)

if ($Help) {
    Write-Host "Uso: .\scripts\deploy-production.ps1 [opciones]" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Opciones:" -ForegroundColor Yellow
    Write-Host "  -FreshDB  Destruye y recrea la base de datos (PELIGROSO: borra todos los datos)" -ForegroundColor Red
    Write-Host "  -Watch    Espera y muestra el progreso del workflow en tiempo real" -ForegroundColor White
    Write-Host "  -Help     Muestra esta ayuda" -ForegroundColor White
    exit 0
}

if (-not (Get-Command gh -ErrorAction SilentlyContinue)) {
    Write-Host "ERR: No se encontro el comando 'gh' (GitHub CLI)." -ForegroundColor Red
    Write-Host "     Instalalo desde: https://cli.github.com/" -ForegroundColor Yellow
    exit 1
}

Write-Host "Lanzando deploy a produccion (master)..." -ForegroundColor Cyan

if ($FreshDB) {
    Write-Host ""
    Write-Host "ADVERTENCIA: -FreshDB activado. Se destruiran TODOS LOS DATOS de produccion." -ForegroundColor Red
    $confirm = Read-Host "Escribe 'SI' para confirmar"
    if ($confirm -ne "SI") {
        Write-Host "Operacion cancelada." -ForegroundColor Yellow
        exit 0
    }
    gh workflow run deploy-production.yml --ref master --field rebuild_images=true --field fresh_db=true
} else {
    gh workflow run deploy-production.yml --ref master
}

if ($LASTEXITCODE -ne 0) {
    Write-Host "ERR: No se pudo lanzar el workflow. Asegurate de estar autenticado con 'gh auth login'." -ForegroundColor Red
    exit $LASTEXITCODE
}

Write-Host "OK: Workflow lanzado." -ForegroundColor Green
Write-Host "    Puedes seguirlo en: https://github.com/jorgecvzc/TramaTex/actions/workflows/deploy-production.yml" -ForegroundColor Gray

if ($Watch) {
    Write-Host ""
    Write-Host "Esperando inicio del run..." -ForegroundColor Yellow
    Start-Sleep -Seconds 3
    gh run watch $(gh run list --workflow=deploy-production.yml --limit=1 --json databaseId --jq '.[0].databaseId')
}
