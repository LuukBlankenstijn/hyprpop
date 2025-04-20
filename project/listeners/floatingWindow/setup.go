package floatingwindow

import (
	"fmt"
	stateDto "hyprwindow/project/dto/state"
	"hyprwindow/project/state"
	"hyprwindow/project/utils/hypr"
	"time"
)

func setup(state *state.GlobalConfig) {
	createWindows(state.GetConfigState().GetAllWindows(), state.GetAppState())
}

func createSingleWindow(window stateDto.WindowConfig, state *state.State) {
	createWindows([]stateDto.WindowConfig{window}, state)
}

func createWindows(windows []stateDto.WindowConfig, state *state.State) {
	// create windows
	for _, window := range windows {
		if window.Type != eventType {
			continue
		}
		err := createChromiumWindow(&window)
		if err != nil {
			fmt.Println(err)
			// TODO: log error
			continue
		}
	}

	// sleep to allow hyprland to create the windows
	time.Sleep(400 * time.Millisecond)

	// store the windows in the state and set default position and size
	for _, window := range windows {
		createdWindow, err := hypr.GetWindowByName(window.Name)
		if err != nil {
			fmt.Println(err)
			// TODO: log error
			continue
		}
		state.UpdateWindow(window.Name, createdWindow)
		err = hypr.SetSize(*createdWindow, window.Size)
		if err != nil {
			fmt.Println(err)
			// TODO: log error
			continue
		}
		err = hypr.SetPosition(*createdWindow, window.Position)
		if err != nil {
			fmt.Println(err)
			// TODO: log error
			continue
		}
	}
}
