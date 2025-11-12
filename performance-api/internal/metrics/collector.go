package metrics

import (
	"context"
	"math"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// SystemMetrics representa las métricas del sistema
type SystemMetrics struct {
	Timestamp    time.Time `json:"timestamp"`
	CPU          CPUInfo   `json:"cpu"`
	Memory       MemoryInfo `json:"memory"`
	Goroutines   int       `json:"goroutines"`
	NumCPU       int       `json:"num_cpu"`
}

// CPUInfo contiene información sobre el uso de CPU
type CPUInfo struct {
	Percent     float64   `json:"percent"`
	PerCPU      []float64 `json:"per_cpu,omitempty"`
	Count       int       `json:"count"`
}

// MemoryInfo contiene información sobre el uso de memoria
type MemoryInfo struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
	Free        uint64  `json:"free"`
}

// Collector gestiona la recolección de métricas del sistema
type Collector struct {
	mu              sync.RWMutex
	currentMetrics  *SystemMetrics
	metricsHistory  []SystemMetrics
	maxHistory      int
	collectionInterval time.Duration
	ctx             context.Context
	cancel          context.CancelFunc
}

// NewCollector crea una nueva instancia del recolector
func NewCollector() *Collector {
	ctx, cancel := context.WithCancel(context.Background())
	return &Collector{
		currentMetrics:    &SystemMetrics{},
		metricsHistory:    make([]SystemMetrics, 0),
		maxHistory:        100, // Mantener últimas 100 métricas
		collectionInterval: 15 * time.Second,
		ctx:               ctx,
		cancel:            cancel,
	}
}

// StartCollection inicia la recolección periódica de métricas
func (c *Collector) StartCollection(interval time.Duration) {
	c.collectionInterval = interval
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Recolectar métricas inmediatamente
	c.collectMetrics()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.collectMetrics()
		}
	}
}

// collectMetrics recolecta las métricas actuales del sistema
func (c *Collector) collectMetrics() {
	metrics := &SystemMetrics{
		Timestamp: time.Now(),
	}

	// Obtener información de CPU
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err == nil && len(cpuPercent) > 0 {
		metrics.CPU.Percent = cpuPercent[0]
	}

	// Obtener porcentaje por CPU
	cpuPercentAll, err := cpu.Percent(time.Second, true)
	if err == nil {
		metrics.CPU.PerCPU = cpuPercentAll
	}

	// Obtener número de CPUs
	cpuCount, err := cpu.Counts(true)
	if err == nil {
		metrics.CPU.Count = cpuCount
	}

	// Obtener información de memoria
	memInfo, err := mem.VirtualMemory()
	if err == nil {
		metrics.Memory.Total = memInfo.Total
		metrics.Memory.Available = memInfo.Available
		metrics.Memory.Used = memInfo.Used
		metrics.Memory.UsedPercent = memInfo.UsedPercent
		metrics.Memory.Free = memInfo.Free
	}

	// Información de goroutines
	metrics.Goroutines = runtime.NumGoroutine()
	metrics.NumCPU = runtime.NumCPU()

	c.mu.Lock()
	c.currentMetrics = metrics
	// Agregar al historial
	c.metricsHistory = append(c.metricsHistory, *metrics)
	if len(c.metricsHistory) > c.maxHistory {
		c.metricsHistory = c.metricsHistory[1:]
	}
	c.mu.Unlock()
}

// GetCurrentMetrics retorna las métricas actuales
func (c *Collector) GetCurrentMetrics() *SystemMetrics {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentMetrics
}

// GetMetricsHistory retorna el historial de métricas
func (c *Collector) GetMetricsHistory() []SystemMetrics {
	c.mu.RLock()
	defer c.mu.RUnlock()
	history := make([]SystemMetrics, len(c.metricsHistory))
	copy(history, c.metricsHistory)
	return history
}

// GetMetricsStats calcula estadísticas del historial de métricas
func (c *Collector) GetMetricsStats() *MetricsStatistics {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.metricsHistory) == 0 {
		return nil
	}

	stats := &MetricsStatistics{
		SampleCount: len(c.metricsHistory),
		TimeRange: TimeRange{
			Start: c.metricsHistory[0].Timestamp,
			End:   c.metricsHistory[len(c.metricsHistory)-1].Timestamp,
		},
	}

	// Calcular estadísticas de CPU
	cpuValues := make([]float64, 0, len(c.metricsHistory))
	for _, m := range c.metricsHistory {
		cpuValues = append(cpuValues, m.CPU.Percent)
	}
	stats.CPU = calculateStats(cpuValues)

	// Calcular estadísticas de memoria
	memUsedValues := make([]float64, 0, len(c.metricsHistory))
	for _, m := range c.metricsHistory {
		memUsedValues = append(memUsedValues, float64(m.Memory.Used))
	}
	stats.Memory = calculateStats(memUsedValues)

	// Calcular estadísticas de goroutines
	goroutineValues := make([]float64, 0, len(c.metricsHistory))
	for _, m := range c.metricsHistory {
		goroutineValues = append(goroutineValues, float64(m.Goroutines))
	}
	stats.Goroutines = calculateStats(goroutineValues)

	return stats
}

// Stop detiene la recolección de métricas
func (c *Collector) Stop() {
	c.cancel()
}

// calculateStats calcula estadísticas básicas de un conjunto de valores
func calculateStats(values []float64) StatInfo {
	if len(values) == 0 {
		return StatInfo{}
	}

	var sum, min, max float64
	min = values[0]
	max = values[0]

	for _, v := range values {
		sum += v
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	mean := sum / float64(len(values))

	// Calcular desviación estándar
	var variance float64
	for _, v := range values {
		diff := v - mean
		variance += diff * diff
	}
	variance = variance / float64(len(values))
	
	// Calcular raíz cuadrada (desviación estándar)
	stdDev := math.Sqrt(variance)

	return StatInfo{
		Min:    min,
		Max:    max,
		Mean:   mean,
		StdDev: stdDev,
	}
}

