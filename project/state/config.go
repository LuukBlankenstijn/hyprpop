package state

import (
	"fmt"
	"hyprwindow/project/dto/state"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

type rawConfig struct {
	Windows []state.WindowConfig `yaml:"windows"`
}

var (
	configPath = ".config/hypr/hyprwindow.yaml"
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

func watchConfig() {
	// create watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("Error creating watcher:", err)
		return
	}

	// defer close watcher
	defer func() {
		err := watcher.Close()
		if err != nil {
			log.Println("Error closing watcher:", err)
		}
	}()

	// watch path to watch
	err = watcher.Add(configPath)
	if err != nil {
		log.Println("Error watching file:", err)
		return
	}

	// watch events
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Rename == fsnotify.Rename {
				// sleep to avoid file not found errors
				time.Sleep(10 * time.Millisecond)
				if err := loadConfig(); err != nil {
					log.Println("Error reloading config:", err)
				}
				err = watcher.Add(configPath)
				if err != nil {
					log.Println("Error adding file to watcher:", err)
					return
				}
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("Config file changed, reloading...")
				if err := loadConfig(); err != nil {
					log.Println("Error reloading config:", err)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Watcher error:", err)
		}
	}
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
