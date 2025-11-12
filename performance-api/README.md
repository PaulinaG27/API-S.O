# API de Análisis de Rendimiento

**Proyecto Final de Sistemas Operativos**  
**Universidad de Antioquia - Facultad de Ingeniería**

## 📋 Descripción

API desarrollada en Golang para la recolección y análisis de métricas de rendimiento de aplicaciones. El proyecto incluye monitoreo de parámetros fundamentales como uso de CPU, uso de memoria, y mecanismos de perfilamiento para funciones específicas de aplicaciones.

## 🎯 Características

- ✅ Recolección de métricas de CPU (porcentaje de uso, por núcleo)
- ✅ Recolección de métricas de memoria (total, disponible, usado, porcentaje)
- ✅ Monitoreo de goroutines y número de CPUs
- ✅ Perfilamiento de CPU usando pprof
- ✅ Perfilamiento de memoria heap
- ✅ Perfilamiento de goroutines
- ✅ Perfilamiento de bloqueos
- ✅ Historial de métricas con estadísticas (min, max, media, desviación estándar)
- ✅ API REST con endpoints documentados

## 🏗️ Arquitectura

```
performance-api/
├── main.go                 # Punto de entrada de la aplicación
├── internal/
│   ├── api/               # Módulo de API REST
│   │   └── router.go      # Configuración de rutas y handlers
│   ├── metrics/           # Módulo de recolección de métricas
│   │   ├── collector.go   # Recolector de métricas del sistema
│   │   └── statistics.go  # Cálculo de estadísticas
│   └── profiler/          # Módulo de perfilamiento
│       └── profiler.go    # Gestión de perfiles pprof
├── test-app/              # Aplicación de prueba para análisis
│   └── main.go
├── Dockerfile
└── README.md
```

## 🚀 Instalación y Uso

### Requisitos

- Go 1.21 o superior
- Git

### Instalación Local

1. **Clonar o navegar al directorio del proyecto:**
```bash
cd performance-api
```

2. **Instalar dependencias:**
```bash
go mod download
```

3. **Ejecutar la API:**
```bash
go run main.go
```

La API estará disponible en `http://localhost:8080`

### Uso con Docker

1. **Construir la imagen:**
```bash
docker build -t performance-api .
```

2. **Ejecutar el contenedor:**
```bash
docker run -p 8080:8080 performance-api
```

## 📡 Endpoints de la API

### Métricas

- **GET `/api/metrics`** - Obtiene las métricas actuales del sistema
- **GET `/api/metrics/history`** - Obtiene el historial de métricas recolectadas
- **GET `/api/metrics/stats`** - Obtiene estadísticas del historial (min, max, media, desviación estándar)

### Perfilamiento

- **GET `/api/profile/cpu?seconds=30`** - Genera un perfil de CPU (por defecto 30 segundos, máximo 300)
- **GET `/api/profile/heap`** - Genera un perfil de memoria heap
- **GET `/api/profile/goroutine`** - Genera un perfil de goroutines
- **GET `/api/profile/block`** - Genera un perfil de bloqueos
- **GET `/api/profile/list`** - Lista los perfiles disponibles

### Utilidades

- **GET `/api/health`** - Estado de salud de la API
- **GET `/`** - Información sobre la API y endpoints disponibles

### Perfilamiento nativo de Go (pprof)

La API también expone los endpoints estándar de pprof en `/debug/pprof/`:
- `/debug/pprof/` - Índice de perfiles
- `/debug/pprof/heap` - Perfil de heap
- `/debug/pprof/profile?seconds=30` - Perfil de CPU
- `/debug/pprof/goroutine` - Perfil de goroutines
- `/debug/pprof/block` - Perfil de bloqueos

## 🧪 Aplicación de Prueba

Para probar la API con una aplicación que consume recursos, puedes usar la aplicación de prueba incluida basada en multiplicación de matrices:

```bash
cd test-app
go run main.go
```

Esta aplicación ejecuta multiplicación de matrices de diferentes tamaños:
- **Versión secuencial**: Multiplicación tradicional sin paralelismo
- **Versión paralela**: Multiplicación usando múltiples goroutines
- **Diferentes tamaños**: Desde 100x100 hasta 800x800 matrices
- **Ejecución continua**: Ejecuta multiplicaciones periódicamente cada 10 segundos
- **Validación**: Verifica que ambas versiones produzcan el mismo resultado
- **Métricas**: Muestra tiempos de ejecución y speedup

La aplicación genera matrices aleatorias automáticamente y ejecuta tanto la versión secuencial como paralela, comparando sus rendimientos. Mantén esta aplicación ejecutándose mientras consultas las métricas en la API para analizar el consumo de CPU, memoria y el comportamiento de las goroutines.

## 📊 Ejemplos de Uso

### Obtener métricas actuales

```bash
curl http://localhost:8080/api/metrics
```

Respuesta ejemplo:
```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "cpu": {
    "percent": 45.2,
    "per_cpu": [42.1, 48.3, 44.5, 46.0],
    "count": 4
  },
  "memory": {
    "total": 17179869184,
    "available": 8589934592,
    "used": 8589934592,
    "used_percent": 50.0,
    "free": 8589934592
  },
  "goroutines": 12,
  "num_cpu": 4
}
```

### Obtener estadísticas

```bash
curl http://localhost:8080/api/metrics/stats
```

### Generar perfil de CPU

```bash
curl http://localhost:8080/api/profile/cpu?seconds=10 > cpu.prof
```

Luego puedes analizarlo con:
```bash
go tool pprof cpu.prof
```

### Generar perfil de memoria

```bash
curl http://localhost:8080/api/profile/heap > heap.prof
go tool pprof heap.prof
```

## 🔬 Análisis Experimental

Para realizar análisis estadístico del rendimiento:

1. Ejecuta la API y la aplicación de prueba simultáneamente
2. Consulta `/api/metrics` periódicamente o usa `/api/metrics/history` para obtener el historial
3. Usa `/api/metrics/stats` para obtener estadísticas calculadas
4. Genera perfiles de CPU y memoria durante diferentes cargas de trabajo
5. Analiza los perfiles con `go tool pprof` para identificar funciones que necesitan optimización

## 🛠️ Tecnologías Utilizadas

- **Golang 1.21** - Lenguaje de programación
- **gorilla/mux** - Router HTTP
- **gopsutil** - Recolección de métricas del sistema
- **pprof** - Perfilamiento de aplicaciones Go
- **Docker** - Contenerización

## 📚 Conceptos de Sistemas Operativos Aplicados

- **Virtualización del CPU**: Medición del consumo de CPU y análisis de planificación
- **Virtualización de memoria**: Monitoreo del uso de memoria y detección de patrones
- **Concurrencia**: Análisis de goroutines y bloqueos
- **Perfilado de funciones**: Identificación de cuellos de botella en el código

## 👥 Autores

- Daniel Andrés Agudelo García
- Paulina García Aristizábal

## 📝 Licencia

Este proyecto es parte del curso de Sistemas Operativos de la Universidad de Antioquia.

