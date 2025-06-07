package floatingwindow

import (
	"fmt"
	"hyprpop/src/dto/pubsub"
	stateDto "hyprpop/src/dto/state"
	hyprapi "hyprpop/src/hypr/api"
	hyprutils "hyprpop/src/hypr/utils"
	"hyprpop/src/logging"
	"strings"
)

func registerKeybind(window stateDto.WindowConfig) {
	if window.Type != eventType {
		return
	}
	event := getEvent(window.Name)
	err := hyprutils.RegisterKeybind(event, window.Keybind)
	if err != nil {
		logging.Warn("failed to register keybind: %+v", err)
	}

}

func deregisterKeybinds() error {
	binds, err := hyprapi.GetAllKeybinds()
	if err != nil {
		return fmt.Errorf("failed to get all keybinds: %w", err)
	}
	for _, b := range binds {
		if b.Dispatcher == "event" && strings.HasPrefix(b.Arg, "hyprpop:"+eventType) {
			keybind := stateDto.Keybind{
				Mod: stateDto.ModToString(b.Mod),
				Key: b.Key,
			}
			_ = hyprapi.DeregisterKeybind(keybind)
		}
	}
	return nil
}

func getEvent(name string) pubsub.Event {
	return pubsub.Event{
		Type: eventType,
		Name: name,
	}
}
