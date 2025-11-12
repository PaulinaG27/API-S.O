# Script de ejemplo para ejecutar la API de Análisis de Rendimiento
# PowerShell Script

Write-Host "🚀 Iniciando API de Análisis de Rendimiento..." -ForegroundColor Green
Write-Host ""

# Verificar que Go esté instalado
$goVersion = go version 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Error: Go no está instalado o no está en el PATH" -ForegroundColor Red
    exit 1
}

Write-Host "✅ Go encontrado: $goVersion" -ForegroundColor Green
Write-Host ""

# Verificar dependencias
Write-Host "📦 Verificando dependencias..." -ForegroundColor Yellow
go mod tidy
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Error al verificar dependencias" -ForegroundColor Red
    exit 1
}

Write-Host "✅ Dependencias verificadas" -ForegroundColor Green
Write-Host ""

# Iniciar la API
Write-Host "🌐 Iniciando servidor en http://localhost:8080" -ForegroundColor Cyan
Write-Host ""
Write-Host "Endpoints disponibles:" -ForegroundColor Yellow
Write-Host "  - GET http://localhost:8080/api/metrics" -ForegroundColor White
Write-Host "  - GET http://localhost:8080/api/metrics/history" -ForegroundColor White
Write-Host "  - GET http://localhost:8080/api/metrics/stats" -ForegroundColor White
Write-Host "  - GET http://localhost:8080/api/profile/cpu?seconds=30" -ForegroundColor White
Write-Host "  - GET http://localhost:8080/api/profile/heap" -ForegroundColor White
Write-Host "  - GET http://localhost:8080/debug/pprof/" -ForegroundColor White
Write-Host ""
Write-Host "Presiona Ctrl+C para detener el servidor" -ForegroundColor Yellow
Write-Host ""

go run main.go

