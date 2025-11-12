package api

import (
	"encoding/json"
	"net/http"
	"performance-api/internal/metrics"
	"performance-api/internal/profiler"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Router gestiona las rutas de la API
type Router struct {
	collector *metrics.Collector
	profiler  *profiler.Profiler
	mux       *mux.Router
}

// NewRouter crea un nuevo router con los handlers configurados
func NewRouter(collector *metrics.Collector, profiler *profiler.Profiler) *Router {
	r := &Router{
		collector: collector,
		profiler:  profiler,
		mux:       mux.NewRouter(),
	}
	
	r.setupRoutes()
	return r
}

// setupRoutes configura todas las rutas de la API
func (r *Router) setupRoutes() {
	// Endpoints de métricas
	r.mux.HandleFunc("/api/metrics", r.handleGetMetrics).Methods("GET")
	r.mux.HandleFunc("/api/metrics/history", r.handleGetMetricsHistory).Methods("GET")
	r.mux.HandleFunc("/api/metrics/stats", r.handleGetMetricsStats).Methods("GET")
	
	// Endpoints de perfilamiento
	r.mux.HandleFunc("/api/profile/cpu", r.handleCPUProfile).Methods("GET")
	r.mux.HandleFunc("/api/profile/heap", r.handleHeapProfile).Methods("GET")
	r.mux.HandleFunc("/api/profile/goroutine", r.handleGoroutineProfile).Methods("GET")
	r.mux.HandleFunc("/api/profile/block", r.handleBlockProfile).Methods("GET")
	r.mux.HandleFunc("/api/profile/list", r.handleListProfiles).Methods("GET")
	
	// Endpoint de salud
	r.mux.HandleFunc("/api/health", r.handleHealth).Methods("GET")
	
	// Endpoint raíz
	r.mux.HandleFunc("/", r.handleRoot).Methods("GET")
}

// ServeHTTP implementa http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// handleGetMetrics retorna las métricas actuales del sistema
func (r *Router) handleGetMetrics(w http.ResponseWriter, req *http.Request) {
	metrics := r.collector.GetCurrentMetrics()
	r.respondJSON(w, http.StatusOK, metrics)
}

// handleGetMetricsHistory retorna el historial de métricas
func (r *Router) handleGetMetricsHistory(w http.ResponseWriter, req *http.Request) {
	history := r.collector.GetMetricsHistory()
	r.respondJSON(w, http.StatusOK, map[string]interface{}{
		"count":   len(history),
		"history": history,
	})
}

// handleGetMetricsStats retorna estadísticas del historial de métricas
func (r *Router) handleGetMetricsStats(w http.ResponseWriter, req *http.Request) {
	stats := r.collector.GetMetricsStats()
	if stats == nil {
		r.respondError(w, http.StatusNotFound, "No hay métricas disponibles aún")
		return
	}
	r.respondJSON(w, http.StatusOK, stats)
}

// handleCPUProfile genera un perfil de CPU
func (r *Router) handleCPUProfile(w http.ResponseWriter, req *http.Request) {
	seconds := 30 // Por defecto 30 segundos
	if s := req.URL.Query().Get("seconds"); s != "" {
		if parsed, err := strconv.Atoi(s); err == nil && parsed > 0 && parsed <= 300 {
			seconds = parsed
		}
	}
	
	profile, err := r.profiler.GetCPUProfile(seconds)
	if err != nil {
		r.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(profile.Data))
}

// handleHeapProfile genera un perfil de memoria heap
func (r *Router) handleHeapProfile(w http.ResponseWriter, req *http.Request) {
	profile, err := r.profiler.GetHeapProfile()
	if err != nil {
		r.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(profile.Data))
}

// handleGoroutineProfile genera un perfil de goroutines
func (r *Router) handleGoroutineProfile(w http.ResponseWriter, req *http.Request) {
	profile, err := r.profiler.GetGoroutineProfile()
	if err != nil {
		r.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(profile.Data))
}

// handleBlockProfile genera un perfil de bloqueos
func (r *Router) handleBlockProfile(w http.ResponseWriter, req *http.Request) {
	profile, err := r.profiler.GetBlockProfile()
	if err != nil {
		r.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(profile.Data))
}

// handleListProfiles lista los perfiles disponibles
func (r *Router) handleListProfiles(w http.ResponseWriter, req *http.Request) {
	profiles := r.profiler.ListProfiles()
	r.respondJSON(w, http.StatusOK, map[string]interface{}{
		"profiles": profiles,
	})
}

// handleHealth retorna el estado de salud de la API
func (r *Router) handleHealth(w http.ResponseWriter, req *http.Request) {
	r.respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"uptime":    "running",
	})
}

// handleRoot retorna información sobre la API
func (r *Router) handleRoot(w http.ResponseWriter, req *http.Request) {
	info := map[string]interface{}{
		"name":        "API de Análisis de Rendimiento",
		"version":     "1.0.0",
		"description": "API para recolectar y analizar métricas de rendimiento de aplicaciones",
		"endpoints": map[string]string{
			"metrics":        "/api/metrics",
			"metrics_history": "/api/metrics/history",
			"metrics_stats":  "/api/metrics/stats",
			"cpu_profile":    "/api/profile/cpu?seconds=30",
			"heap_profile":   "/api/profile/heap",
			"goroutine_profile": "/api/profile/goroutine",
			"block_profile":  "/api/profile/block",
			"health":         "/api/health",
		},
	}
	r.respondJSON(w, http.StatusOK, info)
}

// respondJSON envía una respuesta JSON
func (r *Router) respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// respondError envía una respuesta de error
func (r *Router) respondError(w http.ResponseWriter, status int, message string) {
	r.respondJSON(w, status, map[string]string{
		"error": message,
	})
}

