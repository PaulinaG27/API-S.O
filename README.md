<<<<<<< HEAD
# API-S.O
=======
# Laboratorio 3 - Multiplicación de Matrices con Procesos

**Curso:** Sistemas Operativos  
**Universidad de Antioquia**  
**Autores:** Daniel Agudelo, Paulina Garcia

Este laboratorio implementa la multiplicación de matrices de forma secuencial y paralela en **C** (usando `fork()` y memoria compartida) y en **Go** (usando goroutines y heap compartido). Incluye generación automática de matrices de prueba y scripts para reportes de rendimiento.

---

## 📁 Estructura del Proyecto
Lab03/
├── C/ # Código y datos para versión en C
│ ├── matrix_mul.c
│ ├── A.txt
│ └── B.txt
├── Go/ # Código y datos para versión en Go
│ ├── matrix_mul.go
│ ├── A.txt
│ └── B.txt
├── generate_matrices.py # Generador de matrices aleatorias
├── generate_report.sh # Script automático para reporte en C
└── generate_report_go.sh # Script automático para reporte en Go


---

## ⚙️ Cómo usar (recomendado en Ubuntu/WSL2)

1. Generar matrices de prueba

```bash
cd /ruta/a/Lab03
python3 generate_matrices.py

El script te pedirá:

Filas de A (N)
Columnas de A / Filas de B (M)
Columnas de B (P)
Ejemplo: 500 500 500 → matrices 500×500 y 500×500.


2. Ejecución manual (opcional)

En C:
cd C
gcc -o matrix_mul matrix_mul.c
./matrix_mul 1 A.txt B.txt   # Secuencial (baseline)
./matrix_mul 2 A.txt B.txt   # Paralelo con 2 procesos
./matrix_mul 4 A.txt B.txt   # Paralelo con 4 procesos

En Go:
cd ../Go
go run matrix_mul.go 1 A.txt B.txt
go run matrix_mul.go 2 A.txt B.txt
go run matrix_mul.go 4 A.txt B.txt

Ambos programas generan C.txt con el resultado y muestran:

Tiempo secuencial
Tiempo paralelo
Speedup


3. Generar reporte automático (recomendado)

Los scripts generan un reporte completo con:

Información del sistema
Tabla de tiempos, speedup y eficiencia
Análisis y conclusiones

Para C:
chmod +x generate_report.sh
./generate_report.sh
# Salida: lab3_report_c.txt

Para Go:
chmod +x generate_report_go.sh
./generate_report_go.sh
# Salida: lab3_report_go.txt
>>>>>>> 0e07404 (Primer avance de la API para el proyecto final de S.O)
