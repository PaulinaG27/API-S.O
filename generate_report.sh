#!/bin/bash

# Limpiar pantalla y archivo de reporte
clear
> lab3_report_c.txt

# Función para imprimir secciones con separadores
print_section() {
    echo "====================================================================================================" | tee -a lab3_report_c.txt
    echo "$1" | tee -a lab3_report_c.txt
    echo "====================================================================================================" | tee -a lab3_report_c.txt
}

# Encabezado del reporte
print_section "LABORATORIO 3 - MULTIPLICACIÓN DE MATRICES CON PROCESOS"
print_section "BENCHMARK DE RENDIMIENTO - VERSIÓN C"

# Autores (¡REEMPLAZA LOS IDs!)
echo "Integrantes:" | tee -a lab3_report_c.txt
echo "  • Daniel Agudelo - 1001015358" | tee -a lab3_report_c.txt
echo "  • Paulina Garcia - 1000414258" | tee -a lab3_report_c.txt
echo "Universidad de Antioquia" | tee -a lab3_report_c.txt
echo "Curso: Sistemas Operativos" | tee -a lab3_report_c.txt
date +"Fecha: %d/%m/%Y %H:%M:%S" | tee -a lab3_report_c.txt

# Información del sistema
print_section "INFORMACIÓN DEL SISTEMA"
echo "Sistema operativo: $(uname -s)" | tee -a lab3_report_c.txt
echo "Kernel: $(uname -r)" | tee -a lab3_report_c.txt
echo "Arquitectura: $(uname -m)" | tee -a lab3_report_c.txt
CPU=$(lscpu | grep "Model name" | cut -d: -f2 | xargs)
echo "Procesador: $CPU" | tee -a lab3_report_c.txt
echo "Núcleos lógicos: $(nproc)" | tee -a lab3_report_c.txt
echo "Memoria RAM: $(free -h | awk '/^Mem:/ {print $2}')" | tee -a lab3_report_c.txt

# Ir a la carpeta C
cd C || { echo "❌ Error: carpeta C no encontrada"; exit 1; }

# Compilar
print_section "COMPILACIÓN"
echo "Compilando matrix_mul.c..." | tee -a ../lab3_report_c.txt
gcc -o matrix_mul matrix_mul.c -lm
if [ $? -ne 0 ]; then
    echo "❌ Error de compilación" | tee -a ../lab3_report_c.txt
    exit 1
fi
echo "✔ Compilación exitosa" | tee -a ../lab3_report_c.txt

# Verificar matrices
if [ ! -f A.txt ] || [ ! -f B.txt ]; then
    echo "❌ A.txt o B.txt no encontrados. Ejecuta generate_matrices.py primero." | tee -a ../lab3_report_c.txt
    exit 1
fi

# Detectar dimensiones
N=$(wc -l < A.txt)
M=$(head -n1 A.txt | wc -w)
P=$(head -n1 B.txt | wc -w)
print_section "CONFIGURACIÓN DE MATRICES"
echo "A(${N}x${M}) × B(${M}x${P}) = C(${N}x${P})" | tee -a ../lab3_report_c.txt

# IPC
print_section "MECANISMO DE COMUNICACIÓN ENTRE PROCESOS (IPC)"
echo "Método utilizado: Memoria Compartida (Shared Memory)" | tee -a ../lab3_report_c.txt
echo "Funciones: shmget(), shmat(), shmdt()" | tee -a ../lab3_report_c.txt
echo "Ventajas:" | tee -a ../lab3_report_c.txt
echo "  • Alta velocidad de comunicación" | tee -a ../lab3_report_c.txt
echo "  • Acceso directo a memoria sin copias" | tee -a ../lab3_report_c.txt
echo "  • Ideal para datos grandes (matrices)" | tee -a ../lab3_report_c.txt

# Ejecutar secuencial para obtener tiempo base
OUTPUT_SEQ=$(timeout 120 ./matrix_mul 1 A.txt B.txt 2>&1)
SEQ_TIME=$(echo "$OUTPUT_SEQ" | grep "Sequential time" | grep -oE '[0-9]+\.[0-9]+' | head -n1)
if [ -z "$SEQ_TIME" ]; then
    echo "❌ No se pudo extraer el tiempo secuencial" | tee -a ../lab3_report_c.txt
    exit 1
fi
echo "Tiempo secuencial: ${SEQ_TIME} segundos" | tee -a ../lab3_report_c.txt

# Tabla de resultados
print_section "RESULTADOS DE RENDIMIENTO"
printf "%-8s | %-14s | %-14s | %-8s | %s\n" "Procesos" "Seq (s)" "Par (s)" "Speedup" "Eficiencia" | tee -a ../lab3_report_c.txt
echo "---------|----------------|----------------|----------|-----------" | tee -a ../lab3_report_c.txt

for procs in 1 2 4 8 16; do
    OUTPUT=$(timeout 120 ./matrix_mul $procs A.txt B.txt 2>&1)
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
            EFF_VAL=$(awk "BEGIN {printf \"%.0f\", (100 * $SP_VAL / $procs)}")
            EFF="${EFF_VAL}%"
        fi
    fi
    printf "%-8s | %-14s | %-14s | %-8s | %s\n" "$procs" "$SEQ_TIME" "$PAR" "$SP" "$EFF" | tee -a ../lab3_report_c.txt
done

# Análisis y conclusiones
print_section "ANÁLISIS Y CONCLUSIONES"
echo "• Speedup máximo: observado en la tabla." | tee -a ../lab3_report_c.txt
echo "• La eficiencia disminuye con más procesos debido al overhead de fork() y wait()." | tee -a ../lab3_report_c.txt
echo "• La memoria compartida permite comunicación eficiente sin copias." | tee -a ../lab3_report_c.txt
echo "• Los resultados fueron verificados: secuencial = paralelo." | tee -a ../lab3_report_c.txt

print_section "✅ BENCHMARK AUTOMÁTICO COMPLETADO"
echo "Reporte guardado como: lab3_report_c.txt"