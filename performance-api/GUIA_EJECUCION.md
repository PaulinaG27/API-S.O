# Guía de Ejecución y Validación - API de Análisis de Rendimiento

Esta guía contiene los pasos detallados para ejecutar, probar y validar la API de Análisis de Rendimiento.

## 📋 Índice

1. [Preparación del Entorno](#preparación-del-entorno)
2. [Ejecución de la API](#ejecución-de-la-api)
3. [Validación de Endpoints](#validación-de-endpoints)
4. [Análisis con Aplicación de Prueba](#análisis-con-aplicación-de-prueba)
5. [Generación y Análisis de Perfiles](#generación-y-análisis-de-perfiles)
6. [Análisis Estadístico](#análisis-estadístico)
7. [Troubleshooting](#troubleshooting)

---

## 1. Preparación del Entorno

### 1.1 Verificar Instalación de Go

```powershell
# En PowerShell
go version
```

**Resultado esperado**: Debe mostrar la versión de Go (ej: `go version go1.21.x windows/amd64`)

Si no está instalado, descargar desde: https://golang.org/dl/

### 1.2 Navegar al Directorio del Proyecto

```powershell
cd C:\Users\Pauli\Desktop\Repositorios\Lab03\performance-api
```

### 1.3 Verificar y Descargar Dependencias

```powershell
go mod tidy
go mod download
```

**Resultado esperado**: Debe descargar las dependencias sin errores.

### 1.4 Compilar el Proyecto (Opcional)

```powershell
go build -o performance-api.exe .
```

**Resultado esperado**: Debe crear el archivo `performance-api.exe` sin errores.

---

## 2. Ejecución de la API

### 2.1 Iniciar la API

**Opción A: Usando el script PowerShell**
```powershell
.\run-example.ps1
```

**Opción B: Ejecución directa**
```powershell
go run main.go
```

**Resultado esperado**: Debe mostrar:
```
🚀 API de Análisis de Rendimiento iniciada en http://localhost:8080
📊 Métricas disponibles en http://localhost:8080/api/metrics
🔍 Perfilamiento disponible en http://localhost:8080/debug/pprof/
```

### 2.2 Verificar que la API está Corriendo

Abrir un navegador o usar PowerShell:

```powershell
# Verificar estado de salud
Invoke-RestMethod -Uri "http://localhost:8080/api/health"
```

**Resultado esperado**:
```json
{
  "status": "healthy",
  "timestamp": "2024-11-09T...",
  "uptime": "running"
}
```

---

## 3. Validación de Endpoints

### 3.1 Obtener Información de la API

```powershell
Invoke-RestMethod -Uri "http://localhost:8080/" | ConvertTo-Json -Depth 10
```

**Validación**: Debe mostrar información sobre la API y lista de endpoints.

### 3.2 Obtener Métricas Actuales

```powershell
$metrics = Invoke-RestMethod -Uri "http://localhost:8080/api/metrics"
$metrics | ConvertTo-Json -Depth 10
```

**Validación**: Debe mostrar:
- ✅ `timestamp`: Fecha y hora actual
- ✅ `cpu.percent`: Porcentaje de uso de CPU (0-100)
- ✅ `cpu.count`: Número de CPUs
- ✅ `memory.total`: Memoria total en bytes
- ✅ `memory.used`: Memoria usada en bytes
- ✅ `memory.used_percent`: Porcentaje de memoria usada
- ✅ `goroutines`: Número de goroutines activas
- ✅ `num_cpu`: Número de CPUs lógicos

### 3.3 Obtener Historial de Métricas

Esperar al menos 10-15 segundos después de iniciar la API (para que se recolecten varias métricas), luego:

```powershell
$history = Invoke-RestMethod -Uri "http://localhost:8080/api/metrics/history"
$history | ConvertTo-Json -Depth 10
```

**Validación**: 
- ✅ `count`: Número de métricas en el historial (debe ser > 0)
- ✅ `history`: Array con las métricas recolectadas

### 3.4 Obtener Estadísticas

```powershell
$stats = Invoke-RestMethod -Uri "http://localhost:8080/api/metrics/stats"
$stats | ConvertTo-Json -Depth 10
```

**Validación** (después de recolectar varias métricas):
- ✅ `sample_count`: Número de muestras
- ✅ `time_range`: Rango de tiempo (start y end)
- ✅ `cpu`: Estadísticas de CPU (min, max, mean, std_dev)
- ✅ `memory`: Estadísticas de memoria (min, max, mean, std_dev)
- ✅ `goroutines`: Estadísticas de goroutines (min, max, mean, std_dev)

### 3.5 Verificar Endpoints de Perfilamiento

```powershell
# Listar perfiles disponibles
Invoke-RestMethod -Uri "http://localhost:8080/api/profile/list"
```

**Validación**: Debe mostrar lista de perfiles disponibles.

---

## 4. Análisis con Aplicación de Prueba

### 4.1 Abrir Nueva Terminal

Mantener la API corriendo en una terminal y abrir una nueva terminal para la aplicación de prueba.

### 4.2 Ejecutar Aplicación de Prueba

```powershell
cd C:\Users\Pauli\Desktop\Repositorios\Lab03\performance-api\test-app
go run main.go
```

**Resultado esperado**: Debe mostrar:
```
🚀 Aplicación de Prueba - Multiplicación de Matrices
==================================================
Autores: Daniel Agudelo, Paulina Garcia

💻 CPUs disponibles: [número]

📊 Ejecutando pruebas iniciales...

1. Matrices pequeñas (100x100):
   ✅ Matriz 100x100 | Secuencial: X.XXXs | Paralelo (2 goroutines): X.XXXs | Speedup: X.XXx
...
⏳ Ejecutando multiplicaciones periódicas para análisis continuo...
   Presiona Ctrl+C para detener
```

### 4.3 Monitorear Métricas Durante la Carga

En otra terminal (o en el navegador), consultar métricas mientras la aplicación de prueba está corriendo:

```powershell
# Consultar métricas cada 5 segundos
while ($true) {
    Clear-Host
    Write-Host "=== Métricas del Sistema ===" -ForegroundColor Cyan
    $m = Invoke-RestMethod -Uri "http://localhost:8080/api/metrics"
    Write-Host "CPU: $($m.cpu.percent)%" -ForegroundColor Yellow
    Write-Host "Memoria: $([math]::Round($m.memory.used_percent, 2))% ($([math]::Round($m.memory.used/1GB, 2)) GB usados)" -ForegroundColor Yellow
    Write-Host "Goroutines: $($m.goroutines)" -ForegroundColor Yellow
    Write-Host "Timestamp: $($m.timestamp)" -ForegroundColor Gray
    Start-Sleep -Seconds 5
}
```

**Validación**: 
- ✅ El uso de CPU debe aumentar cuando la aplicación de prueba está activa
- ✅ El número de goroutines puede variar
- ✅ Las métricas deben actualizarse cada 5 segundos

### 4.4 Detener la Aplicación de Prueba

Presionar `Ctrl+C` en la terminal donde está corriendo `test-app`.

---

## 5. Generación y Análisis de Perfiles

### 5.1 Generar Perfil de CPU

**Importante**: Este proceso tomará el tiempo especificado (por defecto 30 segundos).

```powershell
# Generar perfil de CPU durante 30 segundos
Invoke-WebRequest -Uri "http://localhost:8080/api/profile/cpu?seconds=30" -OutFile "cpu.prof"
```

**Nota**: Mientras se genera el perfil, es recomendable tener la aplicación de prueba ejecutándose para obtener datos significativos.

### 5.2 Analizar Perfil de CPU con pprof

```powershell
go tool pprof cpu.prof
```

**Comandos útiles en pprof**:
```
(pprof) top          # Ver top 10 funciones que más CPU consumen
(pprof) top10        # Ver top 10
(pprof) list main    # Ver código de la función main
(pprof) web          # Generar gráfico (requiere Graphviz)
(pprof) quit         # Salir
```

**Validación**: 
- ✅ Debe mostrar funciones ordenadas por consumo de CPU
- ✅ Las funciones de `test-app` deben aparecer si está corriendo

### 5.3 Generar Perfil de Memoria Heap

```powershell
Invoke-WebRequest -Uri "http://localhost:8080/api/profile/heap" -OutFile "heap.prof"
go tool pprof heap.prof
```

**Comandos útiles**:
```
(pprof) top          # Ver funciones que más memoria usan
(pprof) top -cum     # Ver memoria acumulada
(pprof) list [func]  # Ver código de función específica
(pprof) quit
```

### 5.4 Generar Perfil de Goroutines

```powershell
Invoke-WebRequest -Uri "http://localhost:8080/api/profile/goroutine" -OutFile "goroutine.prof"
go tool pprof goroutine.prof
```

**Validación**: Debe mostrar el estado de todas las goroutines activas.

### 5.5 Acceder a pprof Web Interface

Abrir en el navegador:
```
http://localhost:8080/debug/pprof/
```

**Validación**: Debe mostrar la página de índice de pprof con enlaces a diferentes perfiles.

---

## 6. Análisis Estadístico

### 6.1 Recolectar Múltiples Muestras

Dejar la API corriendo durante al menos 2-3 minutos, consultando métricas periódicamente o dejando que se acumule el historial automáticamente.

### 6.2 Obtener Estadísticas Completas

```powershell
$stats = Invoke-RestMethod -Uri "http://localhost:8080/api/metrics/stats"
$stats | ConvertTo-Json -Depth 10
```

### 6.3 Analizar los Resultados

**Validar que las estadísticas sean coherentes**:
- ✅ `cpu.min` ≤ `cpu.mean` ≤ `cpu.max`
- ✅ `cpu.std_dev` ≥ 0
- ✅ `memory.min` ≤ `memory.mean` ≤ `memory.max`
- ✅ `memory.std_dev` ≥ 0
- ✅ `sample_count` debe ser > 0

### 6.4 Comparar Escenarios

**Escenario 1: Sistema en Reposo**
1. Iniciar solo la API
2. Esperar 2 minutos
3. Obtener estadísticas
4. Guardar resultados: `stats_reposo.json`

**Escenario 2: Sistema con Carga**
1. Iniciar la API
2. Iniciar la aplicación de prueba
3. Esperar 2 minutos
4. Obtener estadísticas
5. Guardar resultados: `stats_carga.json`

**Comparación**:
```powershell
# Cargar ambos archivos y comparar
$reposo = Get-Content stats_reposo.json | ConvertFrom-Json
$carga = Get-Content stats_carga.json | ConvertFrom-Json

Write-Host "=== Comparación de CPU ===" -ForegroundColor Cyan
Write-Host "Reposo - Media: $($reposo.cpu.mean)%, Max: $($reposo.cpu.max)%"
Write-Host "Carga  - Media: $($carga.cpu.mean)%, Max: $($carga.cpu.max)%"
Write-Host "Diferencia: $($carga.cpu.mean - $reposo.cpu.mean)%" -ForegroundColor Yellow
```

### 6.5 Identificar Funciones a Optimizar

Basado en los perfiles generados:
1. Identificar funciones con mayor consumo de CPU (`top` en pprof)
2. Identificar funciones con mayor uso de memoria (`top -cum` en pprof)
3. Analizar el código de estas funciones (`list [func]` en pprof)
4. Documentar hallazgos

---

## 7. Troubleshooting

### Problema: "Puerto 8080 ya está en uso"

**Solución**:
```powershell
# Encontrar proceso usando el puerto
netstat -ano | findstr :8080

# Matar el proceso (reemplazar PID con el número encontrado)
taskkill /PID [PID] /F
```

### Problema: "go: cannot find module"

**Solución**:
```powershell
cd performance-api
go mod tidy
go mod download
```

### Problema: "Error al iniciar CPU profile"

**Solución**: Asegurarse de que solo hay un perfil de CPU activo a la vez. Esperar a que termine antes de iniciar otro.

### Problema: "No hay métricas disponibles aún"

**Solución**: Esperar al menos 5-10 segundos después de iniciar la API para que se recolecten las primeras métricas.

### Problema: "go tool pprof no funciona"

**Solución**: Verificar que Go está correctamente instalado:
```powershell
go version
go env GOROOT
go env GOPATH
```

### Problema: Métricas muestran valores incorrectos

**Solución**: 
- Verificar permisos del sistema (en Windows puede requerir ejecutar como administrador)
- Verificar que gopsutil puede acceder a las métricas del sistema

---

## 8. Checklist de Validación Completa

Usa este checklist para asegurar que todo funciona correctamente:

- [ ] Go está instalado y funciona (`go version`)
- [ ] Dependencias se descargaron correctamente (`go mod tidy` sin errores)
- [ ] La API inicia sin errores
- [ ] Endpoint `/api/health` responde correctamente
- [ ] Endpoint `/api/metrics` retorna métricas válidas
- [ ] Endpoint `/api/metrics/history` retorna historial
- [ ] Endpoint `/api/metrics/stats` calcula estadísticas correctamente
- [ ] Endpoint `/api/profile/cpu` genera perfil de CPU
- [ ] Endpoint `/api/profile/heap` genera perfil de memoria
- [ ] `go tool pprof` puede analizar los perfiles generados
- [ ] La aplicación de prueba se ejecuta correctamente
- [ ] Las métricas cambian cuando hay carga en el sistema
- [ ] Las estadísticas son coherentes (min ≤ mean ≤ max)

---

## 9. Próximos Pasos para el Proyecto

1. **Análisis Experimental Completo**:
   - Recolectar datos durante diferentes escenarios
   - Comparar rendimiento de diferentes implementaciones
   - Documentar resultados en un informe

2. **Optimización**:
   - Identificar funciones problemáticas usando perfiles
   - Refactorizar código basado en métricas
   - Medir mejoras después de optimizaciones

3. **Extensión de Funcionalidades**:
   - Agregar métricas de disco I/O
   - Agregar métricas de red
   - Implementar exportación a Prometheus/Grafana
   - Agregar alertas basadas en umbrales

---

## 10. Recursos Adicionales

- **Documentación de pprof**: https://pkg.go.dev/net/http/pprof
- **Documentación de gopsutil**: https://github.com/shirou/gopsutil
- **Go Performance Best Practices**: https://go.dev/doc/effective_go#performance

---

**Última actualización**: 2024-11-09  
**Versión de la API**: 1.0.0

