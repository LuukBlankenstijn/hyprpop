package core

import (
	"hyprpop/src/core/pubsub"
	"hyprpop/src/logging"
	"hyprpop/src/state"
)

type App struct {
	State  *state.GlobalConfig
	PubSub *PubSub
}

func (a *App) RegisterListener(listener func(*state.GlobalConfig)) {
	go listener(a.State)
}

func Initialize() (*App, error) {
	logging.SetupLogger()

	// state
	appState, err := state.InitState()
	if err != nil {
		return nil, err
	}

	// pubsub
	pubsub.Initialize()

	// hpyrsocket listener
	go Listen()

	return &App{
		State: appState,
	}, nil
}
