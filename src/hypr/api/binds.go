package api

import (
	"fmt"
	"hyprpop/src/dto/state"
	"os/exec"
)

func RegisterKeybind(event string, keybind state.Keybind) error {
	cmd := exec.Command(
		"hyprctl",
		"keyword",
		"bind",
		keybind.ToString(),
		"event,",
		event,
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error creating keybind: %w", err)
	}
	return nil
}
