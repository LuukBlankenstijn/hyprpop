package state

import (
	"hyprwindow/project/dto/state"
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
