#!/bin/bash
# Script para setup y ejecutar tests en pcele
# Uso: bash ~/setup-tests.sh

set -e

echo "════════════════════════════════════════════════════════════════"
echo "  SETUP Y EJECUTAR TESTS - TRAMATEX EN PCELE"
echo "════════════════════════════════════════════════════════════════"
echo ""

# 1. Verificar Docker
echo "1️⃣  Verificando Docker..."
docker --version
docker-compose --version
echo "✅ Docker disponible"
echo ""

# 2. Navegar a directorio
echo "2️⃣  Navegando a ~/tramatex..."
cd ~/tramatex
pwd
echo "✅ En directorio correcto"
echo ""

# 3. Listar contenido
echo "3️⃣  Contenido del directorio:"
ls -la | grep -E "(docker|Dockerfile|apps/tramatex-api)"
echo ""

# 4. Levantar servicios con sudo
echo "4️⃣  Levantando servicios Docker..."
sudo docker-compose up -d
echo "✅ Servicios levantados"
echo ""

# 5. Esperar a que servicios estén listos
echo "5️⃣  Esperando 15 segundos a que los servicios estén listos..."
sleep 15
echo "✅ Servicios listos"
echo ""

# 6. Verificar servicios
echo "6️⃣  Verificando servicios corriendo:"
sudo docker-compose ps
echo ""

# 7. Verificar healthcheck
echo "7️⃣  Verificando healthcheck del API..."
sudo docker exec tramatex_api curl -s http://localhost:8080/api/health || echo "⚠️  API aún no listo, continuando..."
echo ""

# 8. Ejecutar tests
echo "8️⃣  EJECUTANDO TESTS..."
echo "════════════════════════════════════════════════════════════════"
sudo docker exec tramatex_api go test ./tests/... -v

# Capturar resultado
TEST_RESULT=$?

echo ""
echo "════════════════════════════════════════════════════════════════"

if [ $TEST_RESULT -eq 0 ]; then
    echo "✅ TODOS LOS TESTS PASARON"
else
    echo "❌ ALGUNOS TESTS FALLARON"
fi

echo "════════════════════════════════════════════════════════════════"
echo ""

# 9. Mostrar logs si hay errores
if [ $TEST_RESULT -ne 0 ]; then
    echo "9️⃣  Mostrando logs de error:"
    sudo docker-compose logs api | tail -30
fi

echo ""
echo "✅ Setup completado"
echo "Resultado: $TEST_RESULT"

exit $TEST_RESULT
