package state

import (
	"net/url"
	"strings"
)

type WindowConfig struct {
	Keybind      Keybind `yaml:"keybind"`
	Position     Vec2    `yaml:"position"`
	Name         string  `yaml:"name"`
	Size         Vec2    `yaml:"size"`
	StartCommand string  `yaml:"startCommand"`
	Type         string  `yaml:"type"`
}

func (w WindowConfig) IsNative() bool {
	u, err := url.ParseRequestURI(w.StartCommand)
	if err != nil {
		return true
	}
	scheme := strings.ToLower(u.Scheme)
	return (scheme != "http" && scheme != "https") || u.Host == ""
}
