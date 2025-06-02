package state

import (
	"errors"
	"sync"
)

type GlobalConfig struct {
	configState *Config
	appState    *State
}

var (
	c          GlobalConfig
	once       sync.Once
	intialized bool = false
)

func InitState() (*GlobalConfig, error) {
	once.Do(func() {
		c = GlobalConfig{
			newConfig(),
			newState(),
		}
	})

	if intialized {
		return nil, errors.New("state is already intialized")
	}

	if err := loadConfig(); err != nil {
		return nil, err
	}

	intialized = true
	return &c, nil
}

func (c *GlobalConfig) GetAppState() *State {
	return c.appState
}

func (c *GlobalConfig) GetConfigState() *Config {
	return c.configState
}
