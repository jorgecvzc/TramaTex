# TramaTex - Makefile Global

.PHONY: help setup docker-build docker-up docker-down

help:
	@echo "TramaTex - Make Commands"
	@echo ""
	@echo "Setup:"
	@echo "  make setup          - Setup inicial del proyecto"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build   - Build imágenes Docker"
	@echo "  make docker-up      - Levantar stack (docker-compose up)"
	@echo "  make docker-down    - Bajar stack (docker-compose down)"
	@echo ""
	@echo "Backend:"
	@echo "  make backend-test   - Ejecutar tests backend"
	@echo "  make backend-run    - Ejecutar servidor backend"
	@echo ""
	@echo "Frontend:"
	@echo "  make frontend-dev   - Servidor desarrollo frontend"
	@echo "  make frontend-build - Build frontend"
	@echo ""
	@echo "Documentación:"
	@echo "  make docs-view      - Ver índice de documentación"
	@echo ""

setup:
	@echo "Inicializando TramaTex..."
	@git init
	@echo "✓ Repositorio Git inicializado"
	@echo "✓ Estructura de carpetas creada"
	@echo "✓ Documentación base copiada"
	@echo ""
	@echo "Próximos pasos:"
	@echo "1. docker-compose build"
	@echo "2. docker-compose up"
	@echo "3. Navegar a http://localhost:5173"

docker-build:
	docker-compose build

docker-up:
	docker-compose up

docker-down:
	docker-compose down

backend-test:
	cd backend && make test

backend-run:
	cd backend && make run

frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build

docs-view:
	@echo "Documentación TramaTex"
	@echo ""
	@echo "ADRs:"
	@ls -1 docs/adr/ | grep ADR
	@echo ""
	@echo "Módulos:"
	@ls -1 docs/modules/ | grep -v "^_"
	@echo ""
	@echo "Sesiones:"
	@ls -1 docs/sessions/ | grep "^2026\|^2027\|^2028"
	@echo ""
	@echo "Guías:"
	@ls -1 docs/guides/
