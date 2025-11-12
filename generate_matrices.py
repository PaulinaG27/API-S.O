import os
import random

def generate_matrix(rows, cols, min_val=1, max_val=10):
    return [[random.randint(min_val, max_val) for _ in range(cols)] for _ in range(rows)]

def save_matrix_to_file(matrix, filepath):
    with open(filepath, 'w') as f:
        for row in matrix:
            f.write(' '.join(map(str, row)) + '\n')

def main():
    print("=== Generador de matrices aleatorias para Lab03 ===")
    
    # Pedir dimensiones
    try:
        N = int(input("Número de filas de A (N): "))
        M = int(input("Número de columnas de A / filas de B (M): "))
        P = int(input("Número de columnas de B (P): "))
    except ValueError:
        print("Error: Ingresa números enteros.")
        return

    if N <= 0 or M <= 0 or P <= 0:
        print("Error: Las dimensiones deben ser positivas.")
        return

    # Generar matrices
    A = generate_matrix(N, M)
    B = generate_matrix(M, P)

    # Rutas de las carpetas
    base_dir = os.path.dirname(os.path.abspath(__file__))
    c_dir = os.path.join(base_dir, "C")
    go_dir = os.path.join(base_dir, "Go")

    # Crear carpetas si no existen
    os.makedirs(c_dir, exist_ok=True)
    os.makedirs(go_dir, exist_ok=True)

    # Guardar en ambas carpetas
    save_matrix_to_file(A, os.path.join(c_dir, "A.txt"))
    save_matrix_to_file(B, os.path.join(c_dir, "B.txt"))

    save_matrix_to_file(A, os.path.join(go_dir, "A.txt"))
    save_matrix_to_file(B, os.path.join(go_dir, "B.txt"))

    print(f"\n✅ Matrices generadas y guardadas en:")
    print(f"  - {c_dir}")
    print(f"  - {go_dir}")
    print(f"\nDimensiones: A = {N}x{M}, B = {M}x{P}, C = {N}x{P}")

if __name__ == "__main__":
    main()