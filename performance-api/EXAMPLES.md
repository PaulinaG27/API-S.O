# Ejemplos de Uso de la API

Este documento contiene ejemplos prácticos de cómo usar la API de Análisis de Rendimiento.

## Requisitos Previos

1. La API debe estar ejecutándose en `http://localhost:8080`
2. Opcionalmente, ejecuta la aplicación de prueba (`test-app/main.go`) para generar carga

## Ejemplos con cURL

### 1. Obtener Métricas Actuales

```bash
curl http://localhost:8080/api/metrics
```

Respuesta esperada:
```json
{
  "timestamp": "2024-11-09T18:30:00Z",
  "cpu": {
    "percent": 25.5,
    "per_cpu": [23.1, 27.9, 24.2, 26.8],
    "count": 4
  },
  "memory": {
    "total": 17179869184,
    "available": 10737418240,
    "used": 6442450944,
    "used_percent": 37.5,
    "free": 10737418240
  },
  "goroutines": 8,
  "num_cpu": 4
}
```

### 2. Obtener Historial de Métricas

```bash
curl http://localhost:8080/api/metrics/history
```

### 3. Obtener Estadísticas del Historial

```bash
curl http://localhost:8080/api/metrics/stats
```

Respuesta esperada:
```json
{
  "sample_count": 50,
  "time_range": {
    "start": "2024-11-09T18:25:00Z",
    "end": "2024-11-09T18:30:00Z"
  },
  "cpu": {
    "min": 15.2,
    "max": 45.8,
    "mean": 28.5,
    "std_dev": 8.3
  },
  "memory": {
    "min": 5000000000,
    "max": 7000000000,
    "mean": 6000000000,
    "std_dev": 500000000
  },
  "goroutines": {
    "min": 5,
    "max": 15,
    "mean": 9.2,
    "std_dev": 2.8
  }
}
```

### 4. Generar Perfil de CPU (10 segundos)

```bash
curl "http://localhost:8080/api/profile/cpu?seconds=10" -o cpu.prof
```

Luego analizar con:
```bash
go tool pprof cpu.prof
```

En el prompt de pprof:
```
(pprof) top
(pprof) top10
(pprof) list nombreFuncion
(pprof) web
(pprof) quit
```

### 5. Generar Perfil de Memoria Heap

```bash
curl http://localhost:8080/api/profile/heap -o heap.prof
go tool pprof heap.prof
```

### 6. Generar Perfil de Goroutines

```bash
curl http://localhost:8080/api/profile/goroutine -o goroutine.prof
go tool pprof goroutine.prof
```

### 7. Generar Perfil de Bloqueos

```bash
curl http://localhost:8080/api/profile/block -o block.prof
go tool pprof block.prof
```

### 8. Verificar Estado de Salud

```bash
curl http://localhost:8080/api/health
```

### 9. Obtener Información de la API

```bash
curl http://localhost:8080/
```

## Ejemplos con PowerShell

### Obtener Métricas y Formatear JSON

```powershell
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/metrics"
$response | ConvertTo-Json -Depth 10
```

### Guardar Perfil de CPU

```powershell
Invoke-WebRequest -Uri "http://localhost:8080/api/profile/cpu?seconds=30" -OutFile "cpu.prof"
```

## Ejemplos con Postman

### Configuración de Colección

1. Crear una nueva colección llamada "Performance API"
2. Agregar las siguientes requests:

#### GET - Métricas Actuales
- **Method**: GET
- **URL**: `http://localhost:8080/api/metrics`
- **Headers**: `Content-Type: application/json`

#### GET - Historial de Métricas
- **Method**: GET
- **URL**: `http://localhost:8080/api/metrics/history`

#### GET - Estadísticas
- **Method**: GET
- **URL**: `http://localhost:8080/api/metrics/stats`

#### GET - Perfil de CPU
- **Method**: GET
- **URL**: `http://localhost:8080/api/profile/cpu?seconds=30`
- **Save Response**: Guardar como archivo `.prof`

## Análisis Experimental

### Escenario 1: Análisis de Carga Normal

1. Iniciar la API
2. Consultar métricas cada 5 segundos durante 2 minutos
3. Calcular estadísticas
4. Identificar patrones de uso

### Escenario 2: Análisis con Carga Alta

1. Iniciar la API
2. Ejecutar la aplicación de prueba (`test-app/main.go`)
3. Generar perfiles de CPU y memoria durante la carga
4. Analizar los perfiles para identificar cuellos de botella

### Escenario 3: Análisis Comparativo

1. Ejecutar diferentes versiones de una aplicación
2. Recolectar métricas para cada versión
3. Comparar estadísticas (media, desviación estándar)
4. Identificar mejoras o regresiones

## Script de Monitoreo Continuo

Crea un script para monitorear métricas continuamente:

```bash
#!/bin/bash
# monitor.sh

API_URL="http://localhost:8080/api/metrics"
OUTPUT_FILE="metrics_log.json"

echo "Iniciando monitoreo cada 5 segundos..."
echo "Presiona Ctrl+C para detener"

while true; do
    timestamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    metrics=$(curl -s "$API_URL")
    echo "{\"timestamp\": \"$timestamp\", \"metrics\": $metrics}" >> "$OUTPUT_FILE"
    sleep 5
done
```

## Análisis de Perfiles con go tool pprof

### Comandos Útiles de pprof

```bash
# Ver las funciones que más CPU consumen
top

# Ver top 10
top10

# Ver código fuente de una función específica
list nombreFuncion

# Generar gráfico (requiere Graphviz)
web

# Ver árbol de llamadas
tree

# Exportar a formato PDF
pdf

# Ver estadísticas de memoria
top -cum
```

### Interpretación de Resultados

- **flat**: Tiempo dedicado directamente a la función
- **cum**: Tiempo acumulado (incluyendo funciones llamadas)
- **%**: Porcentaje del tiempo total

## Integración con Herramientas Externas

### Grafana (Futuro)

La API puede extenderse para exportar métricas en formato compatible con Prometheus/Grafana.

### Logging

Las métricas pueden ser enviadas a sistemas de logging como ELK Stack o Splunk.

