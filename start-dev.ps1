# TramaTex - Script para Levantar el Ambiente de Desarrollo
# Ejecutar desde la raíz del proyecto: .\start-dev.ps1

Write-Host "🚀 Iniciando TramaTex Development Environment..." -ForegroundColor Cyan
Write-Host ""

# Navegar al directorio raíz si no estamos ahí
$rootDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $rootDir

# Verificar que Docker esté corriendo
Write-Host "🔍 Verificando Docker..." -ForegroundColor Yellow
try {
    docker ps > $null 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Docker no está corriendo. Por favor, inicia Docker Desktop." -ForegroundColor Red
        exit 1
    }
    Write-Host "✅ Docker está corriendo" -ForegroundColor Green
} catch {
    Write-Host "❌ Error al verificar Docker: $_" -ForegroundColor Red
    exit 1
}

# Limpiar contenedores previos (opcional)
Write-Host ""
Write-Host "🧹 Limpiando contenedores previos..." -ForegroundColor Yellow
docker-compose -f docker/docker-compose.local.yml --env-file docker/.env down -v 2>$null

# Levantar servicios
Write-Host ""
Write-Host "🏗️  Construyendo y levantando servicios..." -ForegroundColor Yellow
docker-compose -f docker/docker-compose.local.yml --env-file docker/.env up -d --build

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "✅ Servicios iniciados correctamente" -ForegroundColor Green
    Write-Host ""
    Write-Host "📊 Estado de los contenedores:" -ForegroundColor Cyan
    Start-Sleep -Seconds 3
    docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
    
    Write-Host ""
    Write-Host "🔗 URLs disponibles:" -ForegroundColor Cyan
    Write-Host "   API Health: http://localhost:8080/api/health" -ForegroundColor White
    Write-Host "   Database:   localhost:5432 (user: tramatex, db: tramatex)" -ForegroundColor White
    
    Write-Host ""
    Write-Host "📝 Comandos útiles:" -ForegroundColor Cyan
    Write-Host "   Ver logs API:  docker logs tramatex_api -f" -ForegroundColor Gray
    Write-Host "   Ver logs DB:   docker logs tramatex_db -f" -ForegroundColor Gray
    Write-Host "   Detener todo:  docker-compose -f docker/docker-compose.local.yml --env-file docker/.env down" -ForegroundColor Gray
    
    # Esperar a que la API esté lista
    Write-Host ""
    Write-Host "⏳ Esperando que la API esté lista..." -ForegroundColor Yellow
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
        } catch {
            # API no está lista aún
        }
        $attempt++
        Start-Sleep -Seconds 2
        Write-Host "." -NoNewline -ForegroundColor Gray
    }
    
    Write-Host ""
    if ($apiReady) {
        Write-Host "✅ API está lista y respondiendo!" -ForegroundColor Green
    } else {
        Write-Host "⚠️  La API tardó más de lo esperado. Verifica los logs:" -ForegroundColor Yellow
        Write-Host "   docker logs tramatex_api" -ForegroundColor Gray
    }
} else {
    Write-Host ""
    Write-Host "❌ Error al iniciar los servicios" -ForegroundColor Red
    Write-Host "   Verifica los logs con: docker-compose -f docker/docker-compose.local.yml --env-file docker/.env logs" -ForegroundColor Gray
    exit 1
}

Write-Host ""
Write-Host "🎉 Ambiente de desarrollo listo!" -ForegroundColor Green
Write-Host ""
