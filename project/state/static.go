package state

import (
	"hyprwindow/project/dto/state"
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

func (s *Config) getWindowConfig(name string) *state.WindowConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.windowConfigs[name]
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
