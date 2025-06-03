package floatingwindow

import (
	stateDto "hyprpop/src/dto/state"
	hyprapi "hyprpop/src/hypr/api"
	hyprutils "hyprpop/src/hypr/utils"
	"hyprpop/src/logging"
	"hyprpop/src/state"
	"time"
)

func setup(state *state.GlobalConfig) {
	createWindows(state.GetAppState(), state.GetConfigState().GetAllWindows()...)
}

func createSingleWindow(window stateDto.WindowConfig, state *state.State) {
	createWindows(state, &window)
}

func createWindows(state *state.State, windows ...*stateDto.WindowConfig) {
	workspace, err := hyprapi.GetActiveWorkSpace()
	if err != nil {
		logging.Warn("Failed to get active workspace: %+v", err)
	}
	defer func() {
		_ = hyprapi.FocusWorkspace(*workspace)
	}()

	if err := validateWindows(windows); err != nil {
		logging.Fatal("could not validate config windows: %w", err)
	}

	pids := make(map[int]stateDto.WindowConfig)
	// create windows
	for _, window := range windows {
		if window.Type != eventType {
			continue
		}
		pid, err := createWindow(window)
		if err != nil {
			logging.Error("Failed to created window: %w", err)
			continue
		}
		pids[pid] = *window
	}

	// sleep to allow hyprland to create the windows
	time.Sleep(1000 * time.Millisecond)

	// store the windows in the state and set default position and size
	for pid, window := range pids {
		createdWindow, err := hyprapi.GetWindowByPid(pid)
		if err != nil {
			time.Sleep(4 * time.Second)
			createdWindow, err = hyprapi.GetWindowByPid(pid)
			if err != nil {
				logging.Warn("failed to get window by PID: %+v", err)
				continue
			}
		}
		err = hyprapi.SetFloating(createdWindow, true)
		if err != nil {
			logging.Warn("failed to set window floating: %+v", err)
			continue
		}
		err = hyprutils.SetSize(*createdWindow, window.Size)
		if err != nil {
			logging.Warn("failed to set window size: %+v", err)
			continue
		}
		err = hyprutils.SetPosition(*createdWindow, window.Position)
		if err != nil {
			logging.Warn("failed to set window position: %+v", err)
			continue
		}
		err = hyprutils.SyncInSizeAndPos(createdWindow, &window.Size, &window.Position)
		if err != nil {
			logging.Warn("failed to save window state to memory")
			return
		}
		state.UpdateWindow(window.Name, createdWindow)
		err = hyprutils.MoveWindowToWorkspace(createdWindow, specialWorkspaceName, true)
		if err != nil {
			logging.Warn("failed to move window to hidden workspace")
			return
		}
	}
}

func validateWindows(windows []*stateDto.WindowConfig) error {
	for _, w := range windows {
		if w.Size.X.IsNegative || w.Size.Y.IsNegative {
			logging.Fatal("window %s cannot have size less then zero: %+v", nil, w.Name, w.Size)
		}
	}
	return nil
}
