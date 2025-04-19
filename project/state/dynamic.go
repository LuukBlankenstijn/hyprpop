package state

import (
	"hyprwindow/project/dto/state"
	"sync"
)

type State struct {
	mu      sync.RWMutex
	windows map[string]*state.Window
}

func (s *State) AddWindow(name string, createdWindow *state.Window) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.windows[name] = createdWindow
}

func newState() *State {
	return &State{
		windows: make(map[string]*state.Window),
	}
}
