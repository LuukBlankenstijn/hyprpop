package utils

import (
	"hyprpop/project/state"
	"hyprpop/project/utils/hypr"
)

func CleanupState(state *state.State) {
	for _, window := range state.GetAllWindows() {
		hyprWindow, err := hypr.GetWindowByAddress(window.Address)
		if (err != nil) || (hyprWindow == nil) {
			state.RemoveWindow(window)
		}
	}
}
