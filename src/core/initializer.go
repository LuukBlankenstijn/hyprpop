package core

import (
	"hyprpop/src/dto/pubsub"
	"hyprpop/src/logging"
	"hyprpop/src/state"
)

type App struct {
	State  *state.GlobalConfig
	PubSub *PubSub
}

func (a *App) RegisterListener(listener func(*state.GlobalConfig, chan pubsub.Event)) {
	go listener(a.State, a.PubSub.Subscribe(10))
}

func Initialize() (*App, error) {
	logging.SetupLogger()

	// state
	appState, err := state.InitState()
	if err != nil {
		return nil, err
	}

	// pubsub
	pubSub := NewPubSub()

	// hpyrsocket listener
	go Listen(pubSub)

	return &App{
		State:  appState,
		PubSub: pubSub,
	}, nil
}
