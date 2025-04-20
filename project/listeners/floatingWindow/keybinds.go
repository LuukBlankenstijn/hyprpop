package floatingwindow

import (
	"fmt"
	"hyprpop/project/dto/pubsub"
	"hyprpop/project/state"
	"hyprpop/project/utils"
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
			fmt.Println(err)
			// TODO: log error
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
