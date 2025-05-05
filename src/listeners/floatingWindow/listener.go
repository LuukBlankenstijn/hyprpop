package floatingwindow

import (
	"fmt"
	"hyprpop/src/dto/pubsub"
	stateDto "hyprpop/src/dto/state"
	"hyprpop/src/state"
	"hyprpop/src/utils"
	"hyprpop/src/utils/hypr"
)

const (
	eventType            = "floating"
	specialWorkspaceName = "special:hyprpop:" + eventType
)

func StartListening(state *state.GlobalConfig, channel chan pubsub.Event) {
	registerKeybinds(state.GetConfigState())
	setup(state)
	listen(state, channel)
}

func listen(store *state.GlobalConfig, channel chan pubsub.Event) {

	state := store.GetAppState()
	for event := range channel {
		if event.Type != eventType {
			continue
		}

		// execute the required action
		handleEvent(store, event)

		utils.CleanupState(state)
	}
}

func getMemoryWindow(store *state.GlobalConfig, name string) *stateDto.Window {
	state := store.GetAppState()
	config := store.GetConfigState()
	window := state.GetWindow(name)
	if window == nil {
		windowConfig, ok := config.GetWindowConfig(name)
		if !ok {
			fmt.Printf("Window %s not found in config\n", name)
			return nil
		}
		createSingleWindow(*windowConfig, state)
		window = state.GetWindow(name)
	}
	return window
}

func handleEvent(store *state.GlobalConfig, event pubsub.Event) {
	// get the window from the state, or create it if it doesn't exist
	memoryWindow := getMemoryWindow(store, event.Name)

	// get the current window from hyprland
	currentWindow, err := hypr.GetWindowByAddress(memoryWindow.Address)
	if err != nil {
		fmt.Printf("Window %s not found in Hyprland\n", event.Name)
		return
	}

	currentWorkspace, _ := hypr.GetActiveWorkSpace()
	// TODO: handle hyprland errors
	if currentWindow.Workspace == *currentWorkspace {
		activeWindow, _ := hypr.GetActiveWindow()
		if activeWindow.Address == currentWindow.Address {
			// save size and position
			err := hypr.SyncInSizeAndPos(currentWindow)
			if err != nil {
				fmt.Printf("Error syncing size and position: %v\n", err)
				return
			}
			store.GetAppState().UpdateWindow(event.Name, currentWindow)

			// move to special workspace
			_ = hypr.MoveWindowToWorkspace(activeWindow, specialWorkspaceName, true)
		} else {
			_ = hypr.FocusWindow(*currentWindow)
		}
	} else {
		// different workspace
		activeWorkspace, _ := hypr.GetActiveWorkSpace()
		if currentWindow.Workspace.Name == specialWorkspaceName {
			_ = hypr.MoveWindowToWorkspace(currentWindow, activeWorkspace.Name, false)
			currentWindow.Workspace = *activeWorkspace

			// set size and position
			err = hypr.SetSize(*currentWindow, memoryWindow.Size)
			if err != nil {
				fmt.Printf("Error setting size: %v\n", err)
				return
			}
			err = hypr.SetPosition(*currentWindow, memoryWindow.Position)
			if err != nil {
				fmt.Printf("Error setting position: %v\n", err)
				return
			}
		} else {
			// save size and position
			err := hypr.SyncInSizeAndPos(currentWindow)
			if err != nil {
				fmt.Printf("Error syncing size and position: %v\n", err)
				return
			}
			store.GetAppState().UpdateWindow(event.Name, currentWindow)

			// move to current workspace
			_ = hypr.MoveWindowToWorkspace(currentWindow, activeWorkspace.Name, false)
			currentWindow.Workspace = *activeWorkspace

			// set size and position
			err = hypr.SyncOutSizeAndPos(currentWindow)
			if err != nil {
				fmt.Printf("Error syncing size and position: %v\n", err)
				return
			}
		}
	}
}
