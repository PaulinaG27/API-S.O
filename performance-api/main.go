package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"performance-api/internal/api"
	"performance-api/internal/metrics"
	"performance-api/internal/profiler"
	"time"
)

func main() {
	// Inicializar el recolector de métricas
	collector := metrics.NewCollector()
	
	// Inicializar el perfilador
	profiler := profiler.NewProfiler()
	
	// Configurar el router de la API
	router := api.NewRouter(collector, profiler)
	
	// Iniciar recolección de métricas en segundo plano
	go collector.StartCollection(15 * time.Second) // Recolecta cada 5 segundos
	
	// Endpoints de la API
	port := ":8080"
	log.Printf("🚀 API de Análisis de Rendimiento iniciada en http://localhost%s", port)
	log.Printf("📊 Métricas disponibles en http://localhost%s/api/metrics", port)
	log.Printf("🔍 Perfilamiento disponible en http://localhost%s/debug/pprof/", port)
	
	// Iniciar servidor HTTP
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

