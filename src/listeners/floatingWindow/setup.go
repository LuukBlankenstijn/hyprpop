package floatingwindow

import (
	"hyprpop/src/core/pubsub"
	stateDto "hyprpop/src/dto/state"
	hyprapi "hyprpop/src/hypr/api"
	hyprutils "hyprpop/src/hypr/utils"
	"hyprpop/src/logging"
	"hyprpop/src/state"
	"sync"
)

func setup(state *state.GlobalConfig) {
	createWindows(state.GetAppState(), state.GetConfigState().GetAllWindows()...)
}

func createWindow(window *stateDto.WindowConfig, state *state.State, wg *sync.WaitGroup) {
	workspace, err := hyprapi.GetActiveWorkSpace()
	if err != nil {
		logging.Warn("Failed to get active workspace: %+v", err)
	}
	channel := pubsub.Subscribe(10)
	defer func() {
		_ = hyprapi.FocusWorkspace(*workspace)
		pubsub.Unsubscribe(channel)
		wg.Done()
	}()

	pid, err := createHyprWindow(window)
	if err != nil {
		logging.Error("Failed to create window: %w", err)
		return
	}

channelLoop:
	for range channel {
		createdWindow, err := hyprapi.GetWindowByPid(pid)
		if err != nil {
			continue
		}

		err = hyprapi.SetFloating(createdWindow, true)
		if err != nil {
			logging.Warn("failed to set window floating: %+v", err)
			return
		}
		err = hyprutils.SetSize(*createdWindow, window.Size)
		if err != nil {
			logging.Warn("failed to set window size: %+v", err)
			return
		}
		err = hyprutils.SetPosition(*createdWindow, window.Position)
		if err != nil {
			logging.Warn("failed to set window position: %+v", err)
			return
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
		break channelLoop
	}
	registerKeybind(*window)
}

func createSingleWindow(window stateDto.WindowConfig, state *state.State) {
	createWindows(state, &window)
}

func createWindows(state *state.State, windows ...*stateDto.WindowConfig) {
	if err := validateWindows(windows); err != nil {
		logging.Fatal("could not validate config windows: %w", err)
	}
	var wg sync.WaitGroup

	for _, window := range windows {
		wg.Add(1)
		go createWindow(window, state, &wg)
	}
	wg.Wait()
}

func validateWindows(windows []*stateDto.WindowConfig) error {
	for _, w := range windows {
		if w.Size.X.IsNegative || w.Size.Y.IsNegative {
			logging.Fatal("window %s cannot have size less then zero: %+v", nil, w.Name, w.Size)
		}
	}
	return nil
}
