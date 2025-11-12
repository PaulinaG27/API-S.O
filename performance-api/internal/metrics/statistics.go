package metrics

import "time"

// MetricsStatistics contiene estadísticas calculadas del historial de métricas
type MetricsStatistics struct {
	SampleCount int       `json:"sample_count"`
	TimeRange    TimeRange `json:"time_range"`
	CPU          StatInfo  `json:"cpu"`
	Memory       StatInfo  `json:"memory"`
	Goroutines   StatInfo  `json:"goroutines"`
}

// TimeRange representa un rango de tiempo
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// StatInfo contiene estadísticas básicas (min, max, mean, std dev)
type StatInfo struct {
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
	Mean   float64 `json:"mean"`
	StdDev float64 `json:"std_dev"`
}

