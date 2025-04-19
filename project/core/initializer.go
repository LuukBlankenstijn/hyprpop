package core

import (
	"hyprwindow/project/dto/pubsub"
	"hyprwindow/project/state"
	"hyprwindow/project/utils"
)

type App struct {
	State  *state.GlobalConfig
	PubSub *utils.PubSub
}

func (a *App) RegisterListener(listener func(*state.GlobalConfig, chan pubsub.Event)) {
	go listener(a.State, a.PubSub.Subscribe(10))
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
