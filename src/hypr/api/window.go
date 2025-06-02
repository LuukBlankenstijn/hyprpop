package api

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

/**
* Give a window focus.
 */
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

func MoveWindowToTop(window state.Window) error {
	cmd := exec.Command(
		"hyprctl",
		"dispatch",
		"alterzorder",
		"top,",
		"address:"+window.Address,
	)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

/**
* Makes sets the window to float depending on floating
 */
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
