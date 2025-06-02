package utils

import (
	"encoding/json"
	"fmt"
	"hyprpop/src/dto/pubsub"
	"hyprpop/src/dto/state"
	"hyprpop/src/hypr/api"
	"os/exec"
	"strings"
)

func RegisterKeybind(e pubsub.Event, keybind state.Keybind) error {
	if isbound, _ := isBound(e, keybind); isbound {
		return nil
	}

	err := api.RegisterKeybind(e.ToString(), keybind)

	return err
}

func DeregisterAllKeybinds() error {
	binds, err := api.GetAllKeybinds()
	if err != nil {
		return fmt.Errorf("failed to get all keybinds: %w", err)
	}
	for _, b := range binds {
		if b.Dispatcher == "event" && strings.HasPrefix(b.Arg, "hyprpop:") {
			keybind := state.Keybind{
				Mod: state.ModToString(b.Mod),
				Key: b.Key,
			}
			_ = api.DeregisterKeybind(keybind)
		}
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
