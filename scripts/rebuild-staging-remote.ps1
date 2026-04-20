# TramaTex - Rebuild total de staging remoto desde Windows (via SSH)
# Ejecutar desde la raíz del repo en Windows.

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
    Write-Host "  -PreserveDatabase       No elimina volúmenes de base de datos" -ForegroundColor White
    Write-Host "  -SkipImageRemove        No elimina imágenes antes de pull" -ForegroundColor White
    Write-Host ""
    Write-Host "Ejemplos:" -ForegroundColor Yellow
    Write-Host "  .\scripts\rebuild-staging-remote.ps1" -ForegroundColor Gray
    Write-Host "  .\scripts\rebuild-staging-remote.ps1 -CheckoutRef origin/staging" -ForegroundColor Gray
    exit 0
}

if (-not (Get-Command ssh -ErrorAction SilentlyContinue)) {
    Write-Host "ERR: No se encontró el comando 'ssh' en este equipo." -ForegroundColor Red
    exit 1
}

$checkoutValue = if ($NoCheckout) { "" } else { $CheckoutRef }
$preserveValue = if ($PreserveDatabase) { "true" } else { "false" }
$removeImagesValue = if ($SkipImageRemove) { "false" } else { "true" }

# Bloque de comandos remotos (Bash) - Construcción limpia
$remoteScript = "set -euo pipefail`n"
$remoteScript += "cd '$ProjectDir' || { echo 'ERR: No se pudo acceder a $ProjectDir'; exit 1; }`n"
$remoteScript += "echo 'INFO: Directorio actual: '`$(pwd)`n"
$remoteScript += "echo '[0/5] Verificando entorno remoto...'`n"
$remoteScript += "if [ ! -f 'docker/.env' ]; then`n"
$remoteScript += "    if [ -f 'docker/.env.staging.example' ]; then`n"
$remoteScript += "        cp docker/.env.staging.example docker/.env`n"
$remoteScript += "        echo 'WARN: Creado docker/.env desde ejemplo.'`n"
$remoteScript += "    fi`n"
$remoteScript += "fi`n"
$remoteScript += "echo '[1/5] Actualizando repositorio...'`n"
$remoteScript += "git fetch --prune origin`n"
$remoteScript += "if [ -n '$checkoutValue' ]; then`n"
$remoteScript += "    git reset --hard '$checkoutValue'`n"
$remoteScript += "fi`n"
$remoteScript += "echo 'INFO: Limpiando CRLF en scripts...'`n"
$remoteScript += "for f in scripts/*.sh; do sed -i 's/\r//g' `"`$f`"; done`n"
$remoteScript += "chmod +x ./scripts/*.sh`n"
$remoteScript += "echo '[2/5] Lanzando script de reconstruccion...'`n"
$remoteScript += "export PROJECT_DIR='$ProjectDir'`n"
$remoteScript += "export COMPOSE_FILE='docker/docker-compose.remote.yml'`n"
$remoteScript += "export ENV_FILE='docker/.env'`n"
$remoteScript += "export CHECKOUT_REF='$checkoutValue'`n"
$remoteScript += "export PRESERVE_DATABASE='$preserveValue'`n"
$remoteScript += "export REMOVE_IMAGES='$removeImagesValue'`n"
$remoteScript += "tr -d '\r' < ./scripts/rebuild-staging-remote.sh | bash -s`n"

Write-Host "Lanzando rebuild remoto en $User@$RemoteHost..." -ForegroundColor Cyan

# El pipe de PowerShell convierte \n en \r\n, lo que rompe bash en remoto.
# Escribimos los bytes UTF-8 puros directamente al stdin del proceso SSH.
$bytes = [System.Text.Encoding]::UTF8.GetBytes($remoteScript)
$proc = New-Object System.Diagnostics.Process
$proc.StartInfo.FileName = "ssh"
$proc.StartInfo.Arguments = "-o ConnectTimeout=10 $User@$RemoteHost bash -s"
$proc.StartInfo.UseShellExecute = $false
$proc.StartInfo.RedirectStandardInput = $true
$proc.Start() | Out-Null
$proc.StandardInput.BaseStream.Write($bytes, 0, $bytes.Length)
$proc.StandardInput.Close()
$proc.WaitForExit()
$sshExitCode = $proc.ExitCode

if ($sshExitCode -ne 0) {
    Write-Host "`n❌ El rebuild remoto falló con código $sshExitCode." -ForegroundColor Red
    exit $sshExitCode
}

Write-Host "`n✅ Rebuild remoto completado con éxito." -ForegroundColor Green
