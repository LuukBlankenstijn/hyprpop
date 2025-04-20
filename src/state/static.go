package state

import (
	"hyprpop/src/dto/state"
	"sync"
)

type Config struct {
	mu            sync.RWMutex
	windowConfigs map[string]*state.WindowConfig
}

func newConfig() *Config {
	return &Config{
		windowConfigs: make(map[string]*state.WindowConfig),
	}
}

func (s *Config) updateWindowConfig(config state.WindowConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.windowConfigs[config.Name] = &config
}

func (s *Config) GetWindowConfig(name string) (*state.WindowConfig, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	windowConfig, ok := s.windowConfigs[name]
	return windowConfig, ok
}

func (s *Config) GetAllWindows() []state.WindowConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	windows := make([]state.WindowConfig, 0, len(s.windowConfigs))
	for _, v := range s.windowConfigs {
		windows = append(windows, *v)
	}
	return windows
}
