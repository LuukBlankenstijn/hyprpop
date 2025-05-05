package hypr

import (
	"encoding/json"
	"fmt"
	"hyprpop/src/dto/state"
	"os/exec"
)

func getAllWindows() ([]state.Window, error) {
	cmd := exec.Command("hyprctl", "clients", "-j")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute hyprctl: %w", err)
	}

	// Parse JSON output
	var windows []state.Window
	if err := json.Unmarshal(output, &windows); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return windows, nil
}

func GetActiveWindow() (*state.Window, error) {
	cmd := exec.Command("hyprctl", "activewindow", "-j")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute hyprctl: %w", err)
	}

	var window state.Window
	if err := json.Unmarshal(output, &window); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &window, nil
}

func GetWindowByAddress(address string) (*state.Window, error) {
	windows, err := getAllWindows()
	if err != nil {
		return nil, err
	}
	for _, window := range windows {
		if window.Address == address {
			return &window, nil
		}
	}
	return nil, fmt.Errorf("window not found")
}

func GetWindowByName(name string) (*state.Window, error) {
	windows, err := getAllWindows()
	if err != nil {
		return nil, err
	}
	for _, window := range windows {
		if window.Class == name {
			return &window, nil
		}
	}
	return nil, fmt.Errorf("window not found")
}

func GetWindowByPid(pid int) (*state.Window, error) {
	windows, err := getAllWindows()
	if err != nil {
		return nil, err
	}
	for _, window := range windows {
		if window.Pid == pid {
			return &window, nil
		}
	}
	return nil, fmt.Errorf("window not found")
}

func SetSize(window state.Window, size state.Vec2) error {
	cmd := exec.Command("hyprctl",
		"dispatch",
		"resizewindowpixel",
		"exact",
		size.X.GetAsString(),
		size.Y.GetAsString()+",",
		"address:"+window.Address,
	)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func SetPosition(window state.Window, position state.Vec2) error {
	fmt.Printf("%+v\n", position)
	x, y, err := getExactPosition(position)
	if err != nil {
		return err
	}
	cmd := exec.Command("hyprctl",
		"dispatch",
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

func GetRelativeSize(window state.Window) (*state.Vec2, error) {
	hyprlandWindow, err := GetWindowByAddress(window.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to get window by address: %w", err)
	}

	monitor, err := getMonitorByWorkspace(&window.Workspace)
	if err != nil {
		return nil, fmt.Errorf("failed to get active monitor: %w", err)
	}

	relativeSize := state.Vec2{
		X: state.VectorValue{
			Value:        hyprlandWindow.Size.X.Value / float64(monitor.GetWidth()),
			IsPercentage: hyprlandWindow.Size.X.Value < float64(monitor.GetWidth()),
		},
		Y: state.VectorValue{
			Value:        hyprlandWindow.Size.Y.Value / float64(monitor.GetHeight()),
			IsPercentage: hyprlandWindow.Size.Y.Value < float64(monitor.GetHeight()),
		},
	}

	if !relativeSize.X.IsPercentage {
		relativeSize.X.Value = hyprlandWindow.Size.X.Value
	}
	if !relativeSize.Y.IsPercentage {
		relativeSize.Y.Value = hyprlandWindow.Size.Y.Value
	}

	return &relativeSize, nil
}

func GetRelativePosition(window state.Window) (*state.Vec2, error) {
	hyprlandWindow, err := GetWindowByAddress(window.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to get window by address: %w", err)
	}

	monitor, err := getMonitorByWorkspace(&window.Workspace)
	if err != nil {
		return nil, fmt.Errorf("failed to get active monitor: %w", err)
	}

	relativePosition := state.Vec2{
		X: state.VectorValue{
			Value:        hyprlandWindow.Position.X.Value / float64(monitor.GetWidth()),
			IsPercentage: hyprlandWindow.Position.X.Value < float64(monitor.GetWidth()),
		},
		Y: state.VectorValue{
			Value:        hyprlandWindow.Position.Y.Value / float64(monitor.GetHeight()),
			IsPercentage: hyprlandWindow.Position.Y.Value < float64(monitor.GetHeight()),
		},
	}

	if !relativePosition.X.IsPercentage {
		relativePosition.X.Value = hyprlandWindow.Size.X.Value
	}
	if !relativePosition.Y.IsPercentage {
		relativePosition.Y.Value = hyprlandWindow.Size.Y.Value
	}

	return &relativePosition, nil
}

func FocusWindow(window state.Window) error {
	cmd := exec.Command(
		"hyprctl",
		"dispatch",
		"focuswindow",
		"address:"+window.Address,
	)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func SyncInSizeAndPos(window *state.Window) error {
	newSize, err := GetRelativeSize(*window)
	if err != nil {
		return fmt.Errorf("error getting size: %v", err)
	}
	newPosition, err := GetRelativePosition(*window)
	if err != nil {
		fmt.Printf("error getting position: %v", err)
	}
	window.Size = *newSize
	window.Position = *newPosition
	return nil
}

func SyncOutSizeAndPos(window *state.Window) error {
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

func SetFloating(window *state.Window, floating bool) error {
	cmd := exec.Command(
		"hyprctl",
		"dispatch",
		"setfloating",
		"address:"+window.Address,
	)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
