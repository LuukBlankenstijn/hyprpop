package floatingwindow

import (
	"fmt"
	"hyprpop/src/core"
	"hyprpop/src/dto/pubsub"
	stateDto "hyprpop/src/dto/state"
	hyprapi "hyprpop/src/hypr/api"
	hyprutils "hyprpop/src/hypr/utils"
	"hyprpop/src/logging"
	"hyprpop/src/state"
)

const (
	eventType            = "floating"
	specialWorkspaceName = "special:hyprpop:" + eventType
)

var store *state.GlobalConfig

func StartListening(state *state.GlobalConfig, channel chan pubsub.Event) {
	store = state
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

		core.CleanupState(state)
	}
}

func getMemoryWindow(store *state.GlobalConfig, name string) (*stateDto.Window, error) {
	state := store.GetAppState()
	config := store.GetConfigState()
	window := state.GetWindow(name)
	if window == nil {
		windowConfig, ok := config.GetWindowConfig(name)
		if !ok {
			return nil, fmt.Errorf("window %s not found in config", name)
		}
		createSingleWindow(*windowConfig, state)
		window = state.GetWindow(name)
		if window == nil {
			return nil, fmt.Errorf("could not find or create window")
		}
	}
	return window, nil
}

func handleEvent(store *state.GlobalConfig, event pubsub.Event) {
	// get the window from the state, or create it if it doesn't exist
	memoryWindow, err := getMemoryWindow(store, event.Name)
	if err != nil {
		logging.Error("failed getting window from memory: ", err)
		return
	}

	// get the current window from hyprland
	currentWindow, err := hyprapi.GetWindowByAddress(memoryWindow.Address)
	if err != nil {
		logging.Warn("Window %s not found in Hyprland\n", event.Name)
		return
	}

	currentWorkspace, err := hyprapi.GetActiveWorkSpace()
	if err != nil {
		logging.Error("Active workspace not found", err)
		return
	}

	defer func() {
		if err != nil {
			logging.Error("failed to handle event %+v", err, event)
		}
	}()
	if currentWindow.Workspace == *currentWorkspace {
		activeWindow, err := hyprapi.GetActiveWindow()
		if err != nil {
			return
		}
		if activeWindow.Address == currentWindow.Address {
			err = toHiddenWorkspace(currentWindow, event.Name)
		} else {
			_ = hyprapi.FocusWindow(*currentWindow)
			err = hyprapi.MoveWindowToTop(*currentWindow)
		}
	} else {
		if currentWindow.Workspace.Name == specialWorkspaceName {
			err = fromHiddenWorkspace(currentWindow, event.Name)
		} else {
			err = fromOtherWorkspace(currentWindow, event.Name)
		}
	}
}

func toHiddenWorkspace(currentWindow *stateDto.Window, eventName string) error {
	err := hyprutils.SyncInSizeAndPos(currentWindow)
	if err != nil {
		fmt.Printf("Error syncing size and position: %v\n", err)
		return fmt.Errorf("error syncing size and poisition %w", err)
	}
	store.GetAppState().UpdateWindow(eventName, currentWindow)

	// move to special workspace
	return hyprutils.MoveWindowToWorkspace(currentWindow, specialWorkspaceName, true)
}

func fromHiddenWorkspace(currentWindow *stateDto.Window, eventName string) error {
	activeWorkspace, err := hyprapi.GetActiveWorkSpace()
	if err != nil {
		return err
	}
	err = hyprutils.MoveWindowToWorkspace(currentWindow, activeWorkspace.Name, false)
	if err != nil {
		return err
	}

	memoryWindow, err := getMemoryWindow(store, eventName)
	if err != nil {
		return fmt.Errorf("failed to get window from memory: %w", err)
	}

	// set new monitor and workspace manually
	monitor, _ := hyprapi.GetMonitorByWorkspace(activeWorkspace)
	currentWindow.Workspace = *activeWorkspace
	currentWindow.MonitorId = monitor.Id

	// set size and position
	err = hyprutils.SetSize(*currentWindow, memoryWindow.Size)
	if err != nil {
		return fmt.Errorf("error setting size: %w", err)
	}
	err = hyprutils.SetPosition(*currentWindow, memoryWindow.Position)
	if err != nil {
		return fmt.Errorf("error setting position: %w", err)
	}
	err = hyprapi.FocusWindow(*currentWindow)
	if err == nil {
		err = hyprapi.MoveWindowToTop(*currentWindow)
	}
	return err
}

func fromOtherWorkspace(currentWindow *stateDto.Window, eventName string) error {
	// save size and position
	err := hyprutils.SyncInSizeAndPos(currentWindow)
	if err != nil {
		return fmt.Errorf("error syncing size and position: %w", err)
	}
	store.GetAppState().UpdateWindow(eventName, currentWindow)

	activeWorkspace, err := hyprapi.GetActiveWorkSpace()
	if err != nil {
		return fmt.Errorf("error sync size and position: %w", err)
	}

	// move to current workspace
	_ = hyprutils.MoveWindowToWorkspace(currentWindow, activeWorkspace.Name, false)

	// set new monitor and workspace manually
	monitor, _ := hyprapi.GetMonitorByWorkspace(activeWorkspace)
	currentWindow.Workspace = *activeWorkspace
	currentWindow.MonitorId = monitor.Id

	// set size and position
	err = hyprutils.SyncOutSizeAndPos(currentWindow)
	if err != nil {
		return fmt.Errorf("error syncing size and position: %w", err)
	}
	err = hyprapi.FocusWindow(*currentWindow)
	if err == nil {
		err = hyprapi.MoveWindowToTop(*currentWindow)
	}
	return err
}
