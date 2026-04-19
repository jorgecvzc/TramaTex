# TramaTex - Rebuild total de staging remoto desde Windows (via SSH)
# Ejecutar desde la raiz del repo en Windows.

param(
    [string]$RemoteHost = "pcele",
    [string]$User = "ele",
    [string]$ProjectDir = "/opt/tramatex",
    [string]$CheckoutRef = "origin/staging",
    [switch]$NoCheckout,
    [switch]$PreserveDatabase,
    [switch]$SkipImageRemove,
    [switch]$Help
)

if ($Help) {
    Write-Host "Uso: .\scripts\rebuild-staging-remote.ps1 [opciones]" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Opciones:" -ForegroundColor Yellow
    Write-Host "  -RemoteHost <host>      Host remoto (default: pcele)" -ForegroundColor White
    Write-Host "  -User <user>            Usuario SSH (default: ele)" -ForegroundColor White
    Write-Host "  -ProjectDir <path>      Ruta del repo en remoto (default: /opt/tramatex)" -ForegroundColor White
    Write-Host "  -CheckoutRef <ref>      Ref para alinear staging (default: origin/staging)" -ForegroundColor White
    Write-Host "  -NoCheckout             Omite git fetch/checkout/reset en remoto" -ForegroundColor White
    Write-Host "  -PreserveDatabase       No elimina volumenes de base de datos" -ForegroundColor White
    Write-Host "  -SkipImageRemove        No elimina imagenes antes de pull" -ForegroundColor White
    Write-Host ""
    Write-Host "Ejemplos:" -ForegroundColor Yellow
    Write-Host "  .\scripts\rebuild-staging-remote.ps1" -ForegroundColor Gray
    Write-Host "  .\scripts\rebuild-staging-remote.ps1 -CheckoutRef origin/chore/staging-deploy-scripts" -ForegroundColor Gray
    Write-Host "  .\scripts\rebuild-staging-remote.ps1 -NoCheckout -PreserveDatabase" -ForegroundColor Gray
    exit 0
}

if (-not (Get-Command ssh -ErrorAction SilentlyContinue)) {
    Write-Host "No se encontro el comando 'ssh' en este equipo." -ForegroundColor Red
    exit 1
}

$checkoutValue = if ($NoCheckout) { "" } else { $CheckoutRef }
$preserveValue = if ($PreserveDatabase) { "true" } else { "false" }
$removeImagesValue = if ($SkipImageRemove) { "false" } else { "true" }

$remoteScript = @"
set -euo pipefail
cd '$ProjectDir'
if [ ! -f './scripts/rebuild-staging-remote.sh' ]; then
  echo 'Missing scripts/rebuild-staging-remote.sh in remote repo.' >&2
  exit 1
fi
chmod +x ./scripts/rebuild-staging-remote.sh
CHECKOUT_REF='$checkoutValue' PRESERVE_DATABASE='$preserveValue' REMOVE_IMAGES='$removeImagesValue' PROJECT_DIR='$ProjectDir' ./scripts/rebuild-staging-remote.sh
"@

Write-Host "Lanzando rebuild remoto en $User@$RemoteHost ..." -ForegroundColor Cyan
$remoteScript | ssh "$User@$RemoteHost" "bash -s"

if ($LASTEXITCODE -ne 0) {
    Write-Host "El rebuild remoto fallo." -ForegroundColor Red
    exit $LASTEXITCODE
}

Write-Host "Rebuild remoto completado." -ForegroundColor Green
