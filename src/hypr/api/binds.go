package api

import (
	"encoding/json"
	"fmt"
	"hyprpop/src/dto/state"
	"os/exec"
)

func RegisterKeybind(event string, keybind state.Keybind) error {
	cmd := exec.Command(
		"hyprctl",
		"keyword",
		"bind",
		keybind.ToString()+",",
		"event,",
		event,
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error creating keybind: %w", err)
	}
	return nil
}

func DeregisterKeybind(keybind state.Keybind) error {
	cmd := exec.Command(
		"hyprctl",
		"keyword",
		"unbind",
		keybind.ToString(),
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error creating keybind: %w", err)
	}
	return nil
}

func GetAllKeybinds() ([]state.HyprKeybind, error) {
	cmd := exec.Command("hyprctl", "binds", "-j")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute hyprctl: %w", err)
	}

	var binds []state.Keybind
	if err := json.Unmarshal(output, &binds); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	var events []state.HyprKeybind
	if err := json.Unmarshal(output, &events); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return events, nil
}
