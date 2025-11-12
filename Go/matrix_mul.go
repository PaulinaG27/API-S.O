// matrix_mul.go
// Autores: Daniel Agudelo - 1001015358, Paulina Garcia - 1000414258
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func readMatrix(filename string) [][]int {
	file, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintf("No se pudo abrir %s", filename))
	}
	defer file.Close()

	var matrix [][]int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var row []int
		for _, s := range strings.Fields(line) {
			num, _ := strconv.Atoi(s)
			row = append(row, num)
		}
		matrix = append(matrix, row)
	}
	return matrix
}

func writeMatrix(filename string, matrix [][]int) {
	file, _ := os.Create(filename)
	defer file.Close()
	writer := bufio.NewWriter(file)
	for _, row := range matrix {
		for j, val := range row {
			if j > 0 {
				writer.WriteString(" ")
			}
			writer.WriteString(strconv.Itoa(val))
		}
		writer.WriteString("\n")
	}
	writer.Flush()
}

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

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Uso: %s <num_goroutines> <A.txt> <B.txt>\n", os.Args[0])
		os.Exit(1)
	}

	numGoroutines, _ := strconv.Atoi(os.Args[1])
	fileA := os.Args[2]
	fileB := os.Args[3]

	A := readMatrix(fileA)
	B := readMatrix(fileB)

	if len(A) == 0 || len(B) == 0 || len(A[0]) != len(B) {
		fmt.Fprintln(os.Stderr, "Error: Dimensiones de matrices incompatibles")
		os.Exit(1)
	}

	// --- Secuencial ---
	startSeq := time.Now()
	C_seq := multiplySequential(A, B)
	elapsedSeq := time.Since(startSeq)

	// --- Paralelo ---
	startPar := time.Now()
	C_par := multiplyParallel(A, B, numGoroutines)
	elapsedPar := time.Since(startPar)

	// --- Validación: asegurar que ambos resultados sean iguales ---
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
		fmt.Fprintln(os.Stderr, "Error: los resultados secuencial y paralelo no coinciden")
		os.Exit(1)
	}

	// Guardar resultado final
	writeMatrix("C.txt", C_par)

	// Mostrar tiempos
	fmt.Printf("Sequential time: %.3f seconds\n", elapsedSeq.Seconds())
	fmt.Printf("Parallel time (%d goroutines): %.3f seconds\n", numGoroutines, elapsedPar.Seconds())
	if elapsedPar.Seconds() > 0 {
		speedup := elapsedSeq.Seconds() / elapsedPar.Seconds()
		fmt.Printf("Speedup: %.2fx\n", speedup)
	}
}