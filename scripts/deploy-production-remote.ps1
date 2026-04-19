# TramaTex - Deploy to production remote desde Windows (via SSH)
# Ejecutar desde la raiz del repo en Windows.
# Nota: las imagenes deben estar ya publicadas en GHCR (via GitHub Actions o manualmente).
# El servidor NO construye imagenes — solo las descarga de ghcr.io.

param(
    [string]$ProdHost = "",
    [string]$User = "root",
    [string]$ProjectDir = "/opt/tramatex",
    [string]$CheckoutRef = "origin/master",
    [switch]$NoCheckout,
    [switch]$WipeDatabase,
    [string]$GhcrUser = "",
    [string]$GhcrToken = "",
    [switch]$Help
)

if ($Help) {
    Write-Host "Uso: .\scripts\deploy-production-remote.ps1 [opciones]" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Opciones:" -ForegroundColor Yellow
    Write-Host "  -ProdHost <host>         IP o hostname del servidor de produccion (requerido)" -ForegroundColor White
    Write-Host "  -User <user>             Usuario SSH (default: root)" -ForegroundColor White
    Write-Host "  -ProjectDir <path>       Ruta del repo en remoto (default: /opt/tramatex)" -ForegroundColor White
    Write-Host "  -CheckoutRef <ref>       Ref para alinear el servidor (default: origin/master)" -ForegroundColor White
    Write-Host "  -NoCheckout              Omite git fetch/checkout/reset en remoto" -ForegroundColor White
    Write-Host "  -WipeDatabase            Elimina volumen de base de datos (DESTRUCTIVO)" -ForegroundColor White
    Write-Host "  -GhcrUser <user>         Usuario GHCR para docker login en remoto" -ForegroundColor White
    Write-Host "  -GhcrToken <token>       Token GHCR para docker login en remoto" -ForegroundColor White
    Write-Host ""
    Write-Host "Ejemplos:" -ForegroundColor Yellow
    Write-Host "  .\scripts\deploy-production-remote.ps1 -ProdHost 1.2.3.4" -ForegroundColor Gray
    Write-Host "  .\scripts\deploy-production-remote.ps1 -ProdHost 1.2.3.4 -NoCheckout" -ForegroundColor Gray
    Write-Host "  .\scripts\deploy-production-remote.ps1 -ProdHost 1.2.3.4 -WipeDatabase" -ForegroundColor Gray
    Write-Host ""
    Write-Host "Nota: para disparar el deploy completo via GitHub Actions (recomendado):" -ForegroundColor Yellow
    Write-Host "  gh workflow run deploy-production.yml --ref master" -ForegroundColor Gray
    exit 0
}

if (-not $ProdHost) {
    Write-Host "ERROR: -ProdHost es requerido (IP o hostname del servidor de produccion)" -ForegroundColor Red
    Write-Host "Ejemplo: .\scripts\deploy-production-remote.ps1 -ProdHost 1.2.3.4" -ForegroundColor Yellow
    exit 1
}

if (-not (Get-Command ssh -ErrorAction SilentlyContinue)) {
    Write-Host "No se encontro el comando 'ssh' en este equipo." -ForegroundColor Red
    exit 1
}

$checkoutValue   = if ($NoCheckout)    { "" }      else { $CheckoutRef }
$wipeValue       = if ($WipeDatabase)  { "true" }  else { "false" }
$preserveValue   = if ($WipeDatabase)  { "false" } else { "true" }
$ghcrUserValue   = $GhcrUser
$ghcrTokenValue  = $GhcrToken

$remoteScript = @"
set -euo pipefail
cd '$ProjectDir'
if [ ! -f './scripts/deploy-production-remote.sh' ]; then
  echo 'Missing scripts/deploy-production-remote.sh in remote repo.' >&2
  exit 1
fi
chmod +x ./scripts/deploy-production-remote.sh
CHECKOUT_REF='$checkoutValue' \
  PRESERVE_DATABASE='$preserveValue' \
  GHCR_USER='$ghcrUserValue' \
  GHCR_TOKEN='$ghcrTokenValue' \
  PROJECT_DIR='$ProjectDir' \
  ./scripts/deploy-production-remote.sh $(if ($WipeDatabase) { '--wipe-database' } else { '' })
"@

Write-Host "Lanzando deploy a produccion en $User@$ProdHost ..." -ForegroundColor Cyan
Write-Host "  Ref: $(if ($NoCheckout) { '(sin checkout)' } else { $CheckoutRef })" -ForegroundColor Gray
Write-Host "  Base de datos: $(if ($WipeDatabase) { 'WIPE (datos eliminados)' } else { 'preservada' })" -ForegroundColor $(if ($WipeDatabase) { 'Red' } else { 'Gray' })
Write-Host ""

$remoteScript | ssh "$User@$ProdHost" "bash -s"

if ($LASTEXITCODE -ne 0) {
    Write-Host "El deploy a produccion fallo." -ForegroundColor Red
    exit $LASTEXITCODE
}

Write-Host "Deploy a produccion completado." -ForegroundColor Green
