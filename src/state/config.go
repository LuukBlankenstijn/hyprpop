package state

import (
	"fmt"
	"hyprpop/src/dto/state"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type rawConfig struct {
	Windows []state.WindowConfig `yaml:"windows"`
}

var (
	configPath = ".config/hypr/hyprpop.yaml"
)

func loadConfig() error {
	// get path
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// read
	var data []byte
	configPath = filepath.Join(home, configPath)
	data, err = os.ReadFile(configPath)
	if err != nil {
		return err
	}

	// Unmarshal
	var newConfig rawConfig
	if err := yaml.Unmarshal(data, &newConfig); err != nil {
		return err
	}

	// validiate
	if err := validateState(newConfig); err != nil {
		return err
	}

	// save
	for _, config := range newConfig.Windows {
		c.configState.updateWindowConfig(config)
	}
	return nil
}

func validateState(config rawConfig) error {
	// check duplicates
	set := make(map[string]struct{})
	for _, w := range config.Windows {
		if _, exists := set[w.Name]; exists {
			return fmt.Errorf("duplicate window name: %s", w.Name)
		}
		set[w.Name] = struct{}{}
	}
	return nil
}
