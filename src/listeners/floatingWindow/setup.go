package floatingwindow

import (
	"fmt"
	stateDto "hyprpop/src/dto/state"
	"hyprpop/src/state"
	"hyprpop/src/utils/hypr"
	"time"
)

func setup(state *state.GlobalConfig) {
	createWindows(state.GetConfigState().GetAllWindows(), state.GetAppState())
}

func createSingleWindow(window stateDto.WindowConfig, state *state.State) {
	createWindows([]stateDto.WindowConfig{window}, state)
}

func createWindows(windows []stateDto.WindowConfig, state *state.State) {
	pids := make(map[int]stateDto.WindowConfig)
	// create windows
	for _, window := range windows {
		if window.Type != eventType {
			continue
		}
		pid, err := createChromiumWindow(&window)
		if err != nil {
			fmt.Println(err)
			// TODO: log error
			continue
		}
		pids[pid] = window
	}

	// sleep to allow hyprland to create the windows
	time.Sleep(1000 * time.Millisecond)

	// store the windows in the state and set default position and size
	for pid, window := range pids {
		createdWindow, err := hypr.GetWindowByPid(pid)
		if err != nil {
			fmt.Println(err)
			time.Sleep(1 * time.Second)
			createdWindow, err = hypr.GetWindowByPid(pid)
			if err != nil {
				fmt.Println(err)
				// TODO: log error
				continue
			}
		}
		err = hypr.SetFloating(createdWindow, true)
		if err != nil {
			fmt.Println(err)
			// TODO: log error
			continue
		}
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
		err = hypr.SyncInSizeAndPos(createdWindow)
		if err != nil {
			fmt.Printf("Error syncing size and position: %v\n", err)
			return
		}
		state.UpdateWindow(window.Name, createdWindow)
	}
}
