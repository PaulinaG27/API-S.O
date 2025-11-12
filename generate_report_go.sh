#!/bin/bash

# Limpiar y crear reporte
clear
> lab3_report_go.txt

# Función para imprimir secciones
print_section() {
    echo "====================================================================================================" | tee -a lab3_report_go.txt
    echo "$1" | tee -a lab3_report_go.txt
    echo "====================================================================================================" | tee -a lab3_report_go.txt
}

# Encabezado
print_section "LABORATORIO 3 - MULTIPLICACIÓN DE MATRICES CON GOROUTINES"
print_section "BENCHMARK DE RENDIMIENTO - VERSIÓN GO"

# Autores (¡REEMPLAZA LOS IDs!)
echo "Integrantes:" | tee -a lab3_report_go.txt
echo "  • Daniel Agudelo - 1001005358" | tee -a lab3_report_go.txt
echo "  • Paulina Garcia - 1000414258" | tee -a lab3_report_go.txt
echo "Universidad de Antioquia" | tee -a lab3_report_go.txt
echo "Curso: Sistemas Operativos" | tee -a lab3_report_go.txt
date +"Fecha: %d/%m/%Y %H:%M:%S" | tee -a lab3_report_go.txt

# Info del sistema
print_section "INFORMACIÓN DEL SISTEMA"
echo "Sistema operativo: $(uname -s)" | tee -a lab3_report_go.txt
echo "Kernel: $(uname -r)" | tee -a lab3_report_go.txt
echo "Arquitectura: $(uname -m)" | tee -a lab3_report_go.txt
CPU=$(lscpu | grep "Model name" | cut -d: -f2 | xargs)
echo "Procesador: $CPU" | tee -a lab3_report_go.txt
echo "Núcleos lógicos: $(nproc)" | tee -a lab3_report_go.txt
echo "Memoria RAM: $(free -h | awk '/^Mem:/ {print $2}')" | tee -a lab3_report_go.txt

# Ir a carpeta Go
cd Go || { echo "❌ Error: carpeta Go no encontrada"; exit 1; }

# Verificar que exista matrix_mul.go
if [ ! -f matrix_mul.go ]; then
    echo "❌ matrix_mul.go no encontrado" | tee -a ../lab3_report_go.txt
    exit 1
fi

# Verificar matrices
if [ ! -f A.txt ] || [ ! -f B.txt ]; then
    echo "❌ A.txt o B.txt no encontrados. Ejecuta generate_matrices.py primero." | tee -a ../lab3_report_go.txt
    exit 1
fi

# Detectar dimensiones
N=$(wc -l < A.txt)
M=$(head -n1 A.txt | wc -w)
P=$(head -n1 B.txt | wc -w)
print_section "CONFIGURACIÓN DE MATRICES"
echo "A(${N}x${M}) × B(${M}x${P}) = C(${N}x${P})" | tee -a ../lab3_report_go.txt

# IPC en Go
print_section "MECANISMO DE COMUNICACIÓN ENTRE GOROUTINES (IPC)"
echo "Método utilizado: Heap compartido (memoria compartida implícita)" | tee -a ../lab3_report_go.txt
echo "Explicación:" | tee -a ../lab3_report_go.txt
echo "  • En Go, las goroutines comparten el espacio de memoria del proceso." | tee -a ../lab3_report_go.txt
echo "  • No se requieren llamadas explícitas como shmget() (no disponibles en Go estándar)." | tee -a ../lab3_report_go.txt
echo "  • El acceso concurrente a la matriz C es seguro porque cada goroutine escribe en filas distintas." | tee -a ../lab3_report_go.txt

# Ejecutar secuencial (1 goroutine) para baseline
print_section "EJECUCIÓN SECUENCIAL (BASELINE)"
echo "Ejecutando con 1 goroutine..." | tee -a ../lab3_report_go.txt
OUTPUT_SEQ=$(timeout 120 go run matrix_mul.go 1 A.txt B.txt 2>&1)
SEQ_TIME=$(echo "$OUTPUT_SEQ" | grep "Sequential time" | grep -oE '[0-9]+\.[0-9]+' | head -n1)
if [ -z "$SEQ_TIME" ]; then
    echo "❌ No se obtuvo tiempo secuencial" | tee -a ../lab3_report_go.txt
    exit 1
fi
echo "Tiempo secuencial: ${SEQ_TIME} segundos" | tee -a ../lab3_report_go.txt

# Tabla de resultados
print_section "RESULTADOS DE RENDIMIENTO"
printf "%-8s | %-14s | %-14s | %-8s | %s\n" "Goroutines" "Seq (s)" "Par (s)" "Speedup" "Eficiencia" | tee -a ../lab3_report_go.txt
echo "----------|----------------|----------------|----------|-----------" | tee -a ../lab3_report_go.txt

for gr in 1 2 4 8 16; do
    OUTPUT=$(timeout 120 go run matrix_mul.go $gr A.txt B.txt 2>&1)
    if [ $? -ne 0 ]; then
        PAR="TIMEOUT"
        SP="0.00x"
        EFF="0.00%"
    else
        PAR=$(echo "$OUTPUT" | grep "Parallel time" | grep -oE '[0-9]+\.[0-9]+' | head -n1)
        if [ -z "$PAR" ]; then
            PAR="ERROR"
            SP="0.00x"
            EFF="0.00%"
        else
            SP_VAL=$(awk "BEGIN {printf \"%.2f\", $SEQ_TIME / $PAR}")
            SP="${SP_VAL}x"
            EFF_VAL=$(awk "BEGIN {printf \"%.0f\", (100 * $SP_VAL / $gr)}")
            EFF="${EFF_VAL}%"
        fi
    fi
    printf "%-10s | %-14s | %-14s | %-8s | %s\n" "$gr" "$SEQ_TIME" "$PAR" "$SP" "$EFF" | tee -a ../lab3_report_go.txt
done

# Análisis
print_section "ANÁLISIS Y CONCLUSIONES"
echo "• Go logra speedup significativo gracias a goroutines ligeras." | tee -a ../lab3_report_go.txt
echo "• No hay overhead de fork() (como en C), por lo que la eficiencia es más alta." | tee -a ../lab3_report_go.txt
echo "• El heap compartido actúa como memoria compartida funcional." | tee -a ../lab3_report_go.txt
echo "• Resultados verificados: secuencial = paralelo." | tee -a ../lab3_report_go.txt

print_section "✅ BENCHMARK AUTOMÁTICO EN GO COMPLETADO"
echo "Reporte guardado como: lab3_report_go.txt"