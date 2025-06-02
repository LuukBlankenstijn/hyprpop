package core

import (
	hyprapi "hyprpop/src/hypr/api"
	"hyprpop/src/state"
)

func CleanupState(state *state.State) {
	for _, window := range state.GetAllWindows() {
		hyprWindow, err := hyprapi.GetWindowByAddress(window.Address)
		if (err != nil) || (hyprWindow == nil) {
			state.RemoveWindow(window)
		}
	}
}
