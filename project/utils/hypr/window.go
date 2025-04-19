package hypr

import (
	"encoding/json"
	"fmt"
	"hyprwindow/project/dto/state"
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

func GetWindowByPid(pid int) (*state.Window, error) {
	windows, err := getAllWindows()
	if err != nil {
		return nil, err
	}
	fmt.Println(pid)
	for _, window := range windows {
		fmt.Println(window.Pid)
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
	cmd := exec.Command("hyprctl",
		"dispatch",
		"movewindowpixel",
		"exact",
		position.X.GetAsString(),
		position.Y.GetAsString()+",",
		"address:"+window.Address,
	)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
