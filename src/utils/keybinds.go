package utils

import (
	"encoding/json"
	"fmt"
	"hyprpop/src/dto/pubsub"
	"hyprpop/src/dto/state"
	"os/exec"
	"strings"
)

func RegisterKeybind(e pubsub.Event, keybind state.Keybind) error {
	if isbound, _ := isBound(e, keybind); isbound {
		return nil
	}
	k := fmt.Sprintf("%s, %s,", keybind.Mod, keybind.Key)
	event := fmt.Sprintf("hyprpop:%s:%s", e.Type, e.Name)
	cmd := exec.Command(
		"hyprctl",
		"keyword",
		"bind",
		k,
		"event,",
		event,
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error creating keybind: %w", err)
	}
	return nil
}

func isBound(event pubsub.Event, keybind state.Keybind) (bool, error) {
	cmd := exec.Command("hyprctl", "binds", "-j")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to execute hyprctl: %w", err)
	}

	var binds []state.Keybind
	if err := json.Unmarshal(output, &binds); err != nil {
		return false, fmt.Errorf("failed to parse JSON: %w", err)
	}

	var events []struct {
		Dispatcher string `json:"dispatcher"`
		Arg        string `json:"arg"`
	}
	if err := json.Unmarshal(output, &events); err != nil {
		return false, fmt.Errorf("failed to parse JSON: %w", err)
	}

	for index := range binds {
		e := strings.ReplaceAll(events[index].Arg, "hyprpop:", "")
		parts := strings.Split(e, ":")
		if binds[index] == keybind &&
			parts[0] == event.Type &&
			parts[1] == event.Name &&
			events[index].Dispatcher == "event" {
			return true, nil
		}
	}

	return false, nil
}
