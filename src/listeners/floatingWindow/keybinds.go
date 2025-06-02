package floatingwindow

import (
	"hyprpop/src/dto/pubsub"
	"hyprpop/src/logging"
	"hyprpop/src/state"
	"hyprpop/src/utils"
)

func registerKeybinds(config *state.Config) {
	windows := config.GetAllWindows()
	for _, window := range windows {
		if window.Type != eventType {
			continue
		}
		event := getEvent(window.Name)
		err := utils.RegisterKeybind(event, window.Keybind)
		if err != nil {
			logging.Warn("failed to register keybind: %+v", err)
			continue
		}
	}
}

func getEvent(name string) pubsub.Event {
	return pubsub.Event{
		Type: eventType,
		Name: name,
	}
}
