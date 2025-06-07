package utils

import (
	"fmt"
	stateDto "hyprpop/src/dto/state"
	"hyprpop/src/hypr/api"
	"hyprpop/src/state"
	"os"
	"os/exec"
	"syscall"
)

func SetSize(window stateDto.Window, size stateDto.Vec2) error {
	monitor, err := api.GetMonitorById(window.MonitorId)
	if err != nil {
		return err
	}
	x, y, err := size.GetExactSize(*monitor)
	if err != nil {
		return err
	}
	cmd := exec.Command("hyprctl",
		"dispatch",
		"resizewindowpixel",
		"exact",
		x,
		y+",",
		"address:"+window.Address,
	)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func SetPosition(window stateDto.Window, position stateDto.Vec2) error {
	monitor, err := api.GetMonitorById(window.MonitorId)
	if err != nil {
		return err
	}
	x, y, err := position.GetExactPosition(*monitor)
	if err != nil {
		return err
	}
	cmd := exec.Command("hyprctl",
		"dispatch",
		"--",
		"movewindowpixel",
		"exact",
		x,
		y+",",
		"address:"+window.Address,
	)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func GetSize(window stateDto.Window, memorySize stateDto.Vec2) (*stateDto.Vec2, error) {
	// get hyprland objects
	hyprlandWindow, err := api.GetWindowByAddress(window.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to get window by address: %w", err)
	}

	monitor, err := api.GetMonitorById(window.MonitorId)
	if err != nil {
		return nil, fmt.Errorf("failed to get monitor: %w", err)
	}

	// get X
	var x stateDto.VectorValue
	if !memorySize.X.IsPercentage {
		x = memorySize.X
	} else {
		x = stateDto.VectorValue{
			Value:        hyprlandWindow.Size.X.Value / float64(monitor.GetWidth()),
			IsPercentage: true,
			IsNegative:   false,
		}
	}

	// get y
	var y stateDto.VectorValue
	if !memorySize.Y.IsPercentage {
		y = memorySize.Y
	} else {
		y = stateDto.VectorValue{
			Value:        hyprlandWindow.Size.Y.Value / float64(monitor.GetHeight()),
			IsPercentage: true,
			IsNegative:   false,
		}
	}

	return &stateDto.Vec2{
		X: x,
		Y: y,
	}, nil
}

func GetPosition(window stateDto.Window, memoryPosition stateDto.Vec2) (*stateDto.Vec2, error) {
	// get hyprland objects
	hyprlandWindow, err := api.GetWindowByAddress(window.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to get window by address: %w", err)
	}

	monitor, err := api.GetMonitorById(window.MonitorId)
	if err != nil {
		return nil, fmt.Errorf("failed to get monitor: %w", err)
	}

	// normalize position
	if err = makePositionLocal(hyprlandWindow); err != nil {
		return nil, fmt.Errorf("failed to localize position: %w", err)
	}

	// get X
	var x stateDto.VectorValue
	if !memoryPosition.X.IsPercentage {
		x = memoryPosition.X
	} else {
		x = stateDto.VectorValue{
			Value:        hyprlandWindow.Position.X.Value / float64(monitor.GetWidth()),
			IsPercentage: true,
			IsNegative:   false,
		}
	}

	// get y
	var y stateDto.VectorValue
	if !memoryPosition.Y.IsPercentage {
		y = memoryPosition.Y
	} else {
		y = stateDto.VectorValue{
			Value:        hyprlandWindow.Position.Y.Value / float64(monitor.GetHeight()),
			IsPercentage: true,
			IsNegative:   false,
		}
	}

	return &stateDto.Vec2{
		X: x,
		Y: y,
	}, nil
	// relativePosition := stateDto.Vec2{
	// 	X: stateDto.VectorValue{
	// 		Value:        hyprlandWindow.Position.X.Value / float64(monitor.GetWidth()),
	// 		IsPercentage: hyprlandWindow.Position.X.Value < float64(monitor.GetWidth()),
	// 	},
	// 	Y: stateDto.VectorValue{
	// 		Value:        hyprlandWindow.Position.Y.Value / float64(monitor.GetHeight()),
	// 		IsPercentage: hyprlandWindow.Position.Y.Value < float64(monitor.GetHeight()),
	// 	},
	// }
	//
	// if !relativePosition.X.IsPercentage {
	// 	relativePosition.X.Value = hyprlandWindow.Size.X.Value
	// }
	// if !relativePosition.Y.IsPercentage {
	// 	relativePosition.Y.Value = hyprlandWindow.Size.Y.Value
	// }
	//
	// return &relativePosition, nil
}

/**
* Loads the current size and position from the hyprctl api into the window object
 */
func SyncInSizeAndPos(window *stateDto.Window, memorySize *stateDto.Vec2, memoryPosition *stateDto.Vec2) error {
	newSize, err := GetSize(*window, *memorySize)
	if err != nil {
		return fmt.Errorf("error getting size: %v", err)
	}
	window.Size = *newSize
	newPosition, err := GetPosition(*window, *memoryPosition)
	if err != nil {
		fmt.Printf("error getting position: %v", err)
	}
	window.Position = *newPosition
	return nil
}

/**
* Sets the size and position of the window from the object with the hyprctl api
 */
func SyncOutSizeAndPos(window *stateDto.Window) error {
	err := SetSize(*window, window.Size)
	if err != nil {
		return fmt.Errorf("error setting size: %v", err)
	}
	err = SetPosition(*window, window.Position)
	if err != nil {
		return fmt.Errorf("error setting position: %v", err)
	}
	return nil
}

/**
* Gets the global position of a window, and makes it local to the monitor it is on.
 */
func makePositionLocal(window *stateDto.Window) error {
	monitor, err := api.GetMonitorById(window.MonitorId)
	if err != nil {
		return fmt.Errorf("failed to get monitor: %w", err)
	}

	position := stateDto.Vec2{
		X: stateDto.VectorValue{
			Value:        float64(int(window.Position.X.Value) - monitor.X),
			IsPercentage: false,
		},
		Y: stateDto.VectorValue{
			Value:        float64(int(window.Position.Y.Value) - monitor.Y),
			IsPercentage: false,
		},
	}

	window.Position = position
	return nil
}

func KillAllWindows(state state.GlobalConfig) {
	windows := state.GetAppState().GetAllWindows()
	for _, w := range windows {
		process, err := os.FindProcess(w.Pid)
		if err != nil {
			continue
		}

		err = process.Signal(syscall.SIGTERM)
		if err != nil {
			continue
		}
	}
}
