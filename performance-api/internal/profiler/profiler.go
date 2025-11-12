package profiler

import (
	"bytes"
	"fmt"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

// Profiler gestiona el perfilamiento de funciones
type Profiler struct {
	mu sync.RWMutex
	profiles map[string]*ProfileData
}

// ProfileData contiene información de un perfil
type ProfileData struct {
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
	Data      string    `json:"data"`
}

// NewProfiler crea una nueva instancia del perfilador
func NewProfiler() *Profiler {
	return &Profiler{
		profiles: make(map[string]*ProfileData),
	}
}

// GetCPUProfile obtiene el perfil de CPU actual
func (p *Profiler) GetCPUProfile(seconds int) (*ProfileData, error) {
	var buf bytes.Buffer
	
	// Iniciar perfil de CPU
	if err := pprof.StartCPUProfile(&buf); err != nil {
		return nil, fmt.Errorf("error al iniciar CPU profile: %w", err)
	}
	
	// Esperar el tiempo especificado
	time.Sleep(time.Duration(seconds) * time.Second)
	
	// Detener perfil de CPU
	pprof.StopCPUProfile()
	
	profileData := &ProfileData{
		Name:      "cpu",
		Timestamp: time.Now(),
		Data:      buf.String(),
	}
	
	p.mu.Lock()
	p.profiles["cpu"] = profileData
	p.mu.Unlock()
	
	return profileData, nil
}

// GetHeapProfile obtiene el perfil de memoria heap
func (p *Profiler) GetHeapProfile() (*ProfileData, error) {
	var buf bytes.Buffer
	
	// Forzar garbage collection antes de obtener el perfil
	runtime.GC()
	
	// Obtener perfil de heap
	if err := pprof.WriteHeapProfile(&buf); err != nil {
		return nil, fmt.Errorf("error al obtener heap profile: %w", err)
	}
	
	profileData := &ProfileData{
		Name:      "heap",
		Timestamp: time.Now(),
		Data:      buf.String(),
	}
	
	p.mu.Lock()
	p.profiles["heap"] = profileData
	p.mu.Unlock()
	
	return profileData, nil
}

// GetGoroutineProfile obtiene el perfil de goroutines
func (p *Profiler) GetGoroutineProfile() (*ProfileData, error) {
	var buf bytes.Buffer
	
	profile := pprof.Lookup("goroutine")
	if profile == nil {
		return nil, fmt.Errorf("no se pudo obtener el perfil de goroutines")
	}
	
	if err := profile.WriteTo(&buf, 0); err != nil {
		return nil, fmt.Errorf("error al escribir goroutine profile: %w", err)
	}
	
	profileData := &ProfileData{
		Name:      "goroutine",
		Timestamp: time.Now(),
		Data:      buf.String(),
	}
	
	p.mu.Lock()
	p.profiles["goroutine"] = profileData
	p.mu.Unlock()
	
	return profileData, nil
}

// GetBlockProfile obtiene el perfil de bloqueos
func (p *Profiler) GetBlockProfile() (*ProfileData, error) {
	var buf bytes.Buffer
	
	profile := pprof.Lookup("block")
	if profile == nil {
		return nil, fmt.Errorf("no se pudo obtener el perfil de bloqueos")
	}
	
	if err := profile.WriteTo(&buf, 0); err != nil {
		return nil, fmt.Errorf("error al escribir block profile: %w", err)
	}
	
	profileData := &ProfileData{
		Name:      "block",
		Timestamp: time.Now(),
		Data:      buf.String(),
	}
	
	p.mu.Lock()
	p.profiles["block"] = profileData
	p.mu.Unlock()
	
	return profileData, nil
}

// GetProfile obtiene un perfil guardado por nombre
func (p *Profiler) GetProfile(name string) (*ProfileData, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	profile, exists := p.profiles[name]
	return profile, exists
}

// ListProfiles retorna la lista de perfiles disponibles
func (p *Profiler) ListProfiles() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	profiles := make([]string, 0, len(p.profiles))
	for name := range p.profiles {
		profiles = append(profiles, name)
	}
	return profiles
}

