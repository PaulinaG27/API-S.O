// main.go
// Aplicación de prueba para análisis de rendimiento - Multiplicación de Matrices
// Autores: Daniel Agudelo - 1001015358, Paulina Garcia - 1000414258
package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

// generateMatrix genera una matriz aleatoria de dimensiones NxM
func generateMatrix(N, M int) [][]int {
	matrix := make([][]int, N)
	for i := range matrix {
		matrix[i] = make([]int, M)
		for j := range matrix[i] {
			matrix[i][j] = rand.Intn(100) // Valores entre 0 y 99
		}
	}
	return matrix
}

// multiplySequential multiplica dos matrices de forma secuencial
func multiplySequential(A, B [][]int) [][]int {
	N, M := len(A), len(A[0])
	P := len(B[0])
	C := make([][]int, N)
	for i := range C {
		C[i] = make([]int, P)
		for j := 0; j < P; j++ {
			sum := 0
			for k := 0; k < M; k++ {
				sum += A[i][k] * B[k][j]
			}
			C[i][j] = sum
		}
	}
	return C
}

// multiplyParallel multiplica dos matrices de forma paralela usando goroutines
func multiplyParallel(A, B [][]int, numGoroutines int) [][]int {
	N, M := len(A), len(A[0])
	P := len(B[0])
	C := make([][]int, N)
	for i := range C {
		C[i] = make([]int, P)
	}

	done := make(chan bool, numGoroutines)

	rowsPerGoroutine := N / numGoroutines
	extraRows := N % numGoroutines

	for g := 0; g < numGoroutines; g++ {
		go func(goroutineID int) {
			startRow := goroutineID*rowsPerGoroutine + min(goroutineID, extraRows)
			endRow := startRow + rowsPerGoroutine
			if goroutineID < extraRows {
				endRow++
			}
			if endRow > N {
				endRow = N
			}

			for i := startRow; i < endRow; i++ {
				for j := 0; j < P; j++ {
					sum := 0
					for k := 0; k < M; k++ {
						sum += A[i][k] * B[k][j]
					}
					C[i][j] = sum
				}
			}
			done <- true
		}(g)
	}

	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	return C
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// runMatrixMultiplication ejecuta una prueba de multiplicación de matrices
func runMatrixMultiplication(size int, numGoroutines int) {
	// Generar matrices aleatorias
	A := generateMatrix(size, size)
	B := generateMatrix(size, size)

	// Validar dimensiones
	if len(A) == 0 || len(B) == 0 || len(A[0]) != len(B) {
		fmt.Printf("   ⚠️  Error: Dimensiones incompatibles\n")
		return
	}

	// Ejecutar versión secuencial
	startSeq := time.Now()
	C_seq := multiplySequential(A, B)
	elapsedSeq := time.Since(startSeq)

	// Ejecutar versión paralela
	startPar := time.Now()
	C_par := multiplyParallel(A, B, numGoroutines)
	elapsedPar := time.Since(startPar)

	// Validar que ambos resultados sean iguales
	valid := true
	for i := range C_seq {
		for j := range C_seq[i] {
			if C_seq[i][j] != C_par[i][j] {
				valid = false
				break
			}
		}
		if !valid {
			break
		}
	}

	if !valid {
		fmt.Printf("   ❌ Error: Los resultados no coinciden\n")
		return
	}

	// Calcular speedup
	var speedup float64
	if elapsedPar.Seconds() > 0 {
		speedup = elapsedSeq.Seconds() / elapsedPar.Seconds()
	}

	fmt.Printf("   ✅ Matriz %dx%d | Secuencial: %.3fs | Paralelo (%d goroutines): %.3fs | Speedup: %.2fx\n",
		size, size, elapsedSeq.Seconds(), numGoroutines, elapsedPar.Seconds(), speedup)
}

func main() {
	fmt.Println("🚀 Aplicación de Prueba - Multiplicación de Matrices")
	fmt.Println("==================================================")
	fmt.Println("Autores: Daniel Agudelo, Paulina Garcia")
	fmt.Println("")

	// Configurar número de CPUs a usar
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf("💻 CPUs disponibles: %d\n", runtime.NumCPU())
	fmt.Println("")

	// Inicializar generador de números aleatorios
	rand.Seed(time.Now().UnixNano())

	// Ejecutar pruebas iniciales con diferentes tamaños
	fmt.Println("📊 Ejecutando pruebas iniciales...")
	fmt.Println("")

	// Prueba 1: Matrices pequeñas (100x100)
	fmt.Println("1. Matrices pequeñas (100x100):")
	runMatrixMultiplication(100, 2)

	// Prueba 2: Matrices medianas (300x300)
	fmt.Println("\n2. Matrices medianas (300x300):")
	runMatrixMultiplication(300, 4)

	// Prueba 3: Matrices grandes (500x500)
	fmt.Println("\n3. Matrices grandes (500x500):")
	runMatrixMultiplication(500, runtime.NumCPU())

	// Prueba 4: Matrices muy grandes (800x800) - si el sistema lo permite
	fmt.Println("\n4. Matrices muy grandes (800x800):")
	runMatrixMultiplication(800, runtime.NumCPU())

	fmt.Println("\n✅ Pruebas iniciales completadas")
	fmt.Println("\n💡 Mantén esta aplicación ejecutándose mientras consultas la API")
	fmt.Println("   en http://localhost:8080/api/metrics")
	fmt.Println("\n⏳ Ejecutando multiplicaciones periódicas para análisis continuo...")
	fmt.Println("   Presiona Ctrl+C para detener")
	fmt.Println("")

	// Configuración para ejecución continua
	ticker := time.NewTicker(10 * time.Second) // Ejecutar cada 10 segundos
	defer ticker.Stop()

	iteration := 0
	for {
		select {
		case <-ticker.C:
			iteration++
			
			// Variar el tamaño de las matrices para diferentes cargas
			var size int
			var numGoroutines int
			
			switch iteration % 4 {
			case 0:
				size = 200
				numGoroutines = 2
			case 1:
				size = 400
				numGoroutines = 4
			case 2:
				size = 600
				numGoroutines = runtime.NumCPU()
			case 3:
				size = 300
				numGoroutines = runtime.NumCPU()
			}

			fmt.Printf("[Iteración %d] ", iteration)
			runMatrixMultiplication(size, numGoroutines)
			
			// Forzar garbage collection periódicamente
			if iteration%10 == 0 {
				runtime.GC()
				fmt.Printf("   🗑️  Garbage collection ejecutado\n")
			}
		}
	}
}
