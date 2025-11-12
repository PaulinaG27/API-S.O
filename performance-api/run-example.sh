#!/bin/bash
# Script de ejemplo para ejecutar la API de Análisis de Rendimiento
# Linux/Mac Script

echo "🚀 Iniciando API de Análisis de Rendimiento..."
echo ""

# Verificar que Go esté instalado
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go no está instalado o no está en el PATH"
    exit 1
fi

echo "✅ Go encontrado: $(go version)"
echo ""

# Verificar dependencias
echo "📦 Verificando dependencias..."
go mod tidy
if [ $? -ne 0 ]; then
    echo "❌ Error al verificar dependencias"
    exit 1
fi

echo "✅ Dependencias verificadas"
echo ""

# Iniciar la API
echo "🌐 Iniciando servidor en http://localhost:8080"
echo ""
echo "Endpoints disponibles:"
echo "  - GET http://localhost:8080/api/metrics"
echo "  - GET http://localhost:8080/api/metrics/history"
echo "  - GET http://localhost:8080/api/metrics/stats"
echo "  - GET http://localhost:8080/api/profile/cpu?seconds=30"
echo "  - GET http://localhost:8080/api/profile/heap"
echo "  - GET http://localhost:8080/debug/pprof/"
echo ""
echo "Presiona Ctrl+C para detener el servidor"
echo ""

go run main.go

