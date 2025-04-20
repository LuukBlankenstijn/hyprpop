package state

import (
	"hyprpop/src/dto/state"
	"sync"
)

type State struct {
	mu      sync.RWMutex
	windows map[string]*state.Window
}

func newState() *State {
	return &State{
		windows: make(map[string]*state.Window),
	}
}

func (s *State) UpdateWindow(name string, window *state.Window) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.windows[name] = window
}

func (s *State) RemoveWindowByName(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.windows[name]; !ok {
		return
	}
	delete(s.windows, name)
}

func (s *State) RemoveWindow(window *state.Window) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for n, w := range s.windows {
		if w.Address == window.Address {
			delete(s.windows, n)
		}
	}
}

func (s *State) GetWindow(name string) *state.Window {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.windows[name]; !ok {
		return nil
	}
	return s.windows[name]
}

func (s *State) GetAllWindows() []*state.Window {
	s.mu.RLock()
	defer s.mu.RUnlock()
	windows := make([]*state.Window, 0, len(s.windows))
	for _, window := range s.windows {
		windows = append(windows, window)
	}
	return windows
}
