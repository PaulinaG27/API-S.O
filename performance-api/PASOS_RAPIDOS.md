# Pasos Rápidos de Ejecución

Guía rápida para ejecutar y validar la API en menos de 5 minutos.

## 🚀 Inicio Rápido

### Paso 1: Preparar el Entorno (30 segundos)

```powershell
# Navegar al directorio
cd C:\Users\Pauli\Desktop\Repositorios\Lab03\performance-api

# Verificar dependencias
go mod tidy
```

### Paso 2: Iniciar la API (10 segundos)

```powershell
# Opción A: Usar el script
.\run-example.ps1

# Opción B: Ejecución directa
go run main.go
```

**✅ Verificación**: Debe mostrar mensaje de inicio en `http://localhost:8080`

### Paso 3: Validar que Funciona (30 segundos)

En una **nueva terminal** (mantener la API corriendo):

```powershell
cd C:\Users\Pauli\Desktop\Repositorios\Lab03\performance-api

# Ejecutar script de validación
.\validar-api.ps1
```

**✅ Verificación**: Todos los checks deben pasar (✅)

### Paso 4: Probar Endpoints Manualmente (1 minuto)

```powershell
# Ver estado de salud
Invoke-RestMethod -Uri "http://localhost:8080/api/health"

# Ver métricas actuales
$m = Invoke-RestMethod -Uri "http://localhost:8080/api/metrics"
$m | ConvertTo-Json -Depth 10

# Ver información de la API
Invoke-RestMethod -Uri "http://localhost:8080/"
```

**✅ Verificación**: Debe retornar JSON válido con datos

### Paso 5: Ejecutar Aplicación de Prueba (2 minutos)

En una **nueva terminal**:

```powershell
cd C:\Users\Pauli\Desktop\Repositorios\Lab03\performance-api\test-app
go run main.go
```

**✅ Verificación**: Debe ejecutar multiplicaciones de matrices y mantenerse activa, mostrando:
- Pruebas iniciales con diferentes tamaños de matrices
- Tiempos de ejecución secuencial vs paralelo
- Speedup calculado
- Ejecución periódica cada 10 segundos

### Paso 6: Monitorear Métricas con Carga (1 minuto)

En otra terminal o navegador:

```powershell
# Monitoreo continuo
while ($true) {
    Clear-Host
    $m = Invoke-RestMethod -Uri "http://localhost:8080/api/metrics"
    Write-Host "CPU: $($m.cpu.percent)% | Memoria: $([math]::Round($m.memory.used_percent, 2))% | Goroutines: $($m.goroutines)"
    Start-Sleep -Seconds 3
}
```

**✅ Verificación**: Métricas deben cambiar cuando la app de prueba está activa

### Paso 7: Generar y Analizar Perfil (2 minutos)

```powershell
# Generar perfil de CPU (30 segundos)
Invoke-WebRequest -Uri "http://localhost:8080/api/profile/cpu?seconds=30" -OutFile "cpu.prof"

# Analizar perfil
go tool pprof cpu.prof
```

En pprof:
```
(pprof) top
(pprof) quit
```

**✅ Verificación**: Debe mostrar funciones ordenadas por consumo de CPU

---

## 📋 Checklist de Validación Completa

- [ ] ✅ API inicia sin errores
- [ ] ✅ Endpoint `/api/health` responde
- [ ] ✅ Endpoint `/api/metrics` retorna datos válidos
- [ ] ✅ Endpoint `/api/metrics/history` funciona
- [ ] ✅ Endpoint `/api/metrics/stats` calcula estadísticas
- [ ] ✅ Aplicación de prueba se ejecuta
- [ ] ✅ Métricas cambian con carga
- [ ] ✅ Perfiles se generan correctamente
- [ ] ✅ `go tool pprof` puede analizar perfiles

---

## 🐛 Solución Rápida de Problemas

| Problema | Solución |
|----------|----------|
| Puerto 8080 ocupado | `netstat -ano \| findstr :8080` luego `taskkill /PID [PID] /F` |
| "cannot find module" | `go mod tidy` |
| API no responde | Verificar que esté corriendo, revisar logs |
| Métricas vacías | Esperar 5-10 segundos después de iniciar |

---

## 📚 Documentación Completa

Para más detalles, consulta:
- `GUIA_EJECUCION.md` - Guía completa paso a paso
- `README.md` - Documentación general del proyecto
- `EXAMPLES.md` - Ejemplos de uso avanzados

---

**Tiempo total estimado**: ~5 minutos para validación completa

