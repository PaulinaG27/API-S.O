# Dockerfile para la API de Análisis de Rendimiento
FROM golang:1.21-alpine AS builder

# Instalar dependencias del sistema necesarias para gopsutil
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum* ./
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o performance-api .

# Imagen final
FROM alpine:latest

# Instalar dependencias runtime necesarias
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar el binario compilado
COPY --from=builder /app/performance-api .

# Exponer el puerto
EXPOSE 8080

# Ejecutar la aplicación
CMD ["./performance-api"]

