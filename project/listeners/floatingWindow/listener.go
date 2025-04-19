package floatingwindow

import (
	"fmt"
	"hyprwindow/project/dto/pubsub"
	stateDto "hyprwindow/project/dto/state"
	"hyprwindow/project/state"
	"hyprwindow/project/utils"
	"hyprwindow/project/utils/hypr"
	"time"
)

const eventType = "floating"

func StartListening(state *state.GlobalConfig, channel chan pubsub.Event) {
	registerKeybinds(state.GetConfigState())
	setup(state)
	listen(channel)
}

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

func setup(state *state.GlobalConfig) {
	// create windows
	windows := state.GetConfigState().GetAllWindows()
	pids := make(map[int]stateDto.WindowConfig)
	for _, window := range windows {
		if window.Type != eventType {
			continue
		}
		pid, err := createChromiumWindow(&window)
		if err != nil {
			fmt.Println(err)
			// TODO: log error
			continue
		}
		pids[pid] = window
	}

	// sleep to allow hyprland to create the windows
	time.Sleep(400 * time.Millisecond)

	// store the windows in the state and set default position and size
	for pid, window := range pids {
		createdWindow, err := hypr.GetWindowByPid(pid)
		if err != nil {
			fmt.Println(err)
			// TODO: log error
			continue
		}
		state.GetAppState().AddWindow(window.Name, createdWindow)
		err = hypr.SetSize(*createdWindow, window.Size)
		if err != nil {
			fmt.Println(err)
			// TODO: log error
			continue
		}
		err = hypr.SetPosition(*createdWindow, window.Size)
		if err != nil {
			fmt.Println(err)
			// TODO: log error
			continue
		}
	}
}

func listen(chan pubsub.Event) {

}
