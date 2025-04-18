package core

import (
	"hyprwindow/project/state"
	"hyprwindow/project/utils"
)

type App struct {
	State  *state.GlobalConfig
	PubSub *utils.PubSub
}

func Initialize() (*App, error) {
	// state
	appState, err := state.InitState()
	if err != nil {
		return nil, err
	}

	// pubsub
	pubSub := utils.NewPubSub()

	// hpyrsocket listener
	go Listen(pubSub)

	return &App{
		State:  appState,
		PubSub: pubSub,
	}, nil
}
