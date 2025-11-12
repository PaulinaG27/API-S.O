# Script de Validación Rápida de la API
# Ejecuta pruebas básicas para verificar que todo funciona

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Validación de la API de Rendimiento" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$API_URL = "http://localhost:8080"
$errors = 0

# Función para probar un endpoint
function Test-Endpoint {
    param(
        [string]$Name,
        [string]$Url,
        [bool]$ShouldReturnJson = $true
    )
    
    Write-Host "Probando: $Name" -ForegroundColor Yellow -NoNewline
    Write-Host " -> $Url" -ForegroundColor Gray
    
    try {
        $response = Invoke-RestMethod -Uri $Url -ErrorAction Stop
        if ($ShouldReturnJson) {
            if ($response) {
                Write-Host "  ✅ OK" -ForegroundColor Green
                return $true
            } else {
                Write-Host "  ❌ Respuesta vacía" -ForegroundColor Red
                return $false
            }
        } else {
            Write-Host "  ✅ OK" -ForegroundColor Green
            return $true
        }
    } catch {
        Write-Host "  ❌ Error: $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

# Verificar que la API esté corriendo
Write-Host "1. Verificando que la API esté corriendo..." -ForegroundColor Cyan
try {
    $health = Invoke-RestMethod -Uri "$API_URL/api/health" -ErrorAction Stop
    if ($health.status -eq "healthy") {
        Write-Host "  ✅ API está corriendo y saludable" -ForegroundColor Green
    } else {
        Write-Host "  ⚠️  API responde pero estado: $($health.status)" -ForegroundColor Yellow
    }
} catch {
    Write-Host "  ❌ ERROR: La API no está corriendo en $API_URL" -ForegroundColor Red
    Write-Host "     Por favor, inicia la API primero con: go run main.go" -ForegroundColor Yellow
    exit 1
}

Write-Host ""

# Probar endpoints
Write-Host "2. Probando endpoints de la API..." -ForegroundColor Cyan
Write-Host ""

if (-not (Test-Endpoint "Información de la API" "$API_URL/")) { $errors++ }
Start-Sleep -Seconds 1

if (-not (Test-Endpoint "Estado de salud" "$API_URL/api/health")) { $errors++ }
Start-Sleep -Seconds 1

if (-not (Test-Endpoint "Métricas actuales" "$API_URL/api/metrics")) { $errors++ }
Start-Sleep -Seconds 2

if (-not (Test-Endpoint "Historial de métricas" "$API_URL/api/metrics/history")) { $errors++ }
Start-Sleep -Seconds 1

if (-not (Test-Endpoint "Lista de perfiles" "$API_URL/api/profile/list")) { $errors++ }
Start-Sleep -Seconds 1

# Verificar métricas específicas
Write-Host ""
Write-Host "3. Validando estructura de métricas..." -ForegroundColor Cyan
try {
    $metrics = Invoke-RestMethod -Uri "$API_URL/api/metrics"
    $valid = $true
    
    $checks = @(
        @{Name="timestamp"; Value=$metrics.timestamp},
        @{Name="cpu.percent"; Value=$metrics.cpu.percent},
        @{Name="cpu.count"; Value=$metrics.cpu.count},
        @{Name="memory.total"; Value=$metrics.memory.total},
        @{Name="memory.used"; Value=$metrics.memory.used},
        @{Name="memory.used_percent"; Value=$metrics.memory.used_percent},
        @{Name="goroutines"; Value=$metrics.goroutines},
        @{Name="num_cpu"; Value=$metrics.num_cpu}
    )
    
    foreach ($check in $checks) {
        if ($null -eq $check.Value) {
            Write-Host "  ❌ Falta: $($check.Name)" -ForegroundColor Red
            $valid = $false
            $errors++
        } else {
            Write-Host "  ✅ $($check.Name): $($check.Value)" -ForegroundColor Green
        }
    }
    
    if ($valid) {
        Write-Host ""
        Write-Host "  ✅ Todas las métricas están presentes" -ForegroundColor Green
    }
} catch {
    Write-Host "  ❌ Error al validar métricas: $($_.Exception.Message)" -ForegroundColor Red
    $errors++
}

# Verificar estadísticas (puede que no haya suficientes muestras aún)
Write-Host ""
Write-Host "4. Verificando estadísticas..." -ForegroundColor Cyan
try {
    $stats = Invoke-RestMethod -Uri "$API_URL/api/metrics/stats"
    if ($stats.sample_count -gt 0) {
        Write-Host "  ✅ Estadísticas disponibles ($($stats.sample_count) muestras)" -ForegroundColor Green
        Write-Host "     CPU - Min: $([math]::Round($stats.cpu.min, 2))%, Max: $([math]::Round($stats.cpu.max, 2))%, Media: $([math]::Round($stats.cpu.mean, 2))%" -ForegroundColor Gray
    } else {
        Write-Host "  ⚠️  No hay suficientes muestras aún (espera unos minutos)" -ForegroundColor Yellow
    }
} catch {
    Write-Host "  ⚠️  Estadísticas no disponibles aún (normal si la API acaba de iniciar)" -ForegroundColor Yellow
}

# Resumen
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
if ($errors -eq 0) {
    Write-Host "  ✅ Validación completada sin errores" -ForegroundColor Green
} else {
    Write-Host "  ⚠️  Se encontraron $errors error(es)" -ForegroundColor Yellow
}
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Mostrar información útil
Write-Host "Próximos pasos:" -ForegroundColor Cyan
Write-Host "  1. Ejecuta la aplicación de prueba: cd test-app; go run main.go" -ForegroundColor White
Write-Host "  2. Monitorea métricas mientras la app corre" -ForegroundColor White
Write-Host "  3. Genera perfiles: Invoke-WebRequest -Uri '$API_URL/api/profile/cpu?seconds=30' -OutFile 'cpu.prof'" -ForegroundColor White
Write-Host "  4. Analiza perfiles: go tool pprof cpu.prof" -ForegroundColor White
Write-Host ""

