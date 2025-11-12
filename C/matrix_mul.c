// matrix_mul.c
// Autores: Daniel Agudelo - 1001015358, Paulina Garcia - 1000414258
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>
#include <sys/shm.h>
#include <time.h>

int main(int argc, char *argv[]) {
    if (argc != 4) {
        fprintf(stderr, "Uso: %s <num_procesos> <A.txt> <B.txt>\n", argv[0]);
        exit(1);
    }

    int num_processes = atoi(argv[1]);
    char *fileA = argv[2];
    char *fileB = argv[3];

    FILE *fa = fopen(fileA, "r");
    FILE *fb = fopen(fileB, "r");
    if (!fa || !fb) {
        perror("Error abriendo archivos");
        exit(1);
    }

    // Contar filas y columnas de A
    int N = 0, M = 0;
    char line[4096];
    while (fgets(line, sizeof(line), fa)) {
        if (line[0] == '\n') continue;
        N++;
        if (N == 1) {
            // Contar columnas en primera línea
            char *token = strtok(line, " \t\r\n");
            while (token) {
                M++;
                token = strtok(NULL, " \t\r\n");
            }
        }
    }

    // Contar filas y columnas de B
    int B_rows = 0, P = 0;
    rewind(fb);
    while (fgets(line, sizeof(line), fb)) {
        if (line[0] == '\n') continue;
        B_rows++;
        if (B_rows == 1) {
            char *token = strtok(line, " \t\r\n");
            while (token) {
                P++;
                token = strtok(NULL, " \t\r\n");
            }
        }
    }

    // Validar compatibilidad
    if (M != B_rows) {
        fprintf(stderr, "Error: Dimensiones incompatibles. A: %dx%d, B: %dx%d\n", N, M, B_rows, P);
        exit(1);
    }

    // Leer matriz A
    int *A = malloc(N * M * sizeof(int));
    rewind(fa);
    for (int i = 0; i < N; i++) {
        for (int j = 0; j < M; j++) {
            fscanf(fa, "%d", &A[i * M + j]);
        }
    }
    fclose(fa);

    // Leer matriz B
    int *B = malloc(M * P * sizeof(int));
    rewind(fb);
    for (int i = 0; i < M; i++) {
        for (int j = 0; j < P; j++) {
            fscanf(fb, "%d", &B[i * P + j]);
        }
    }
    fclose(fb);

    // =============================
    // SECUENCIAL
    // =============================
    clock_t start_seq = clock();
    int *C_seq = malloc(N * P * sizeof(int));
    for (int i = 0; i < N; i++) {
        for (int j = 0; j < P; j++) {
            C_seq[i * P + j] = 0;
            for (int k = 0; k < M; k++) {
                C_seq[i * P + j] += A[i * M + k] * B[k * P + j];
            }
        }
    }
    clock_t end_seq = clock();
    double time_seq = ((double)(end_seq - start_seq)) / CLOCKS_PER_SEC;

    // =============================
    // PARALELO
    // =============================
    int shmid = shmget(IPC_PRIVATE, N * P * sizeof(int), IPC_CREAT | 0666);
    if (shmid == -1) {
        perror("shmget falló");
        exit(1);
    }
    int *C_par = (int *)shmat(shmid, NULL, 0);
    if (C_par == (void *)-1) {
        perror("shmat falló");
        exit(1);
    }

    for (int i = 0; i < N * P; i++) C_par[i] = 0;

    clock_t start_par = clock();

    if (num_processes > N) num_processes = N;
    int rows_per_proc = N / num_processes;
    int extra_rows = N % num_processes;

    for (int p = 0; p < num_processes; p++) {
        pid_t pid = fork();
        if (pid == 0) {
            int start_row = p * rows_per_proc + (p < extra_rows ? p : extra_rows);
            int end_row = start_row + rows_per_proc + (p < extra_rows ? 1 : 0);
            if (end_row > N) end_row = N;

            for (int i = start_row; i < end_row; i++) {
                for (int j = 0; j < P; j++) {
                    C_par[i * P + j] = 0;
                    for (int k = 0; k < M; k++) {
                        C_par[i * P + j] += A[i * M + k] * B[k * P + j];
                    }
                }
            }
            exit(0);
        }
    }

    for (int p = 0; p < num_processes; p++) {
        wait(NULL);
    }

    clock_t end_par = clock();
    double time_par = ((double)(end_par - start_par)) / CLOCKS_PER_SEC;

    // =============================
    // GUARDAR RESULTADO
    // =============================
    FILE *fc = fopen("C.txt", "w");
    for (int i = 0; i < N; i++) {
        for (int j = 0; j < P; j++) {
            fprintf(fc, "%d ", C_par[i * P + j]);
        }
        fprintf(fc, "\n");
    }
    fclose(fc);

    // =============================
    // MOSTRAR TIEMPOS
    // =============================
    printf("Sequential time: %.3f seconds\n", time_seq);
    printf("Parallel time (%d processes): %.3f seconds\n", num_processes, time_par);
    if (time_par > 0) {
        printf("Speedup: %.2fx\n", time_seq / time_par);
    }

    // =============================
    // LIMPIEZA
    // =============================
    shmdt(C_par);
    shmctl(shmid, IPC_RMID, NULL);
    free(A);
    free(B);
    free(C_seq);

    return 0;
}