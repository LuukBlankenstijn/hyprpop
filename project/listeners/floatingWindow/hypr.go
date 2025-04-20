package floatingwindow

import (
	"hyprpop/project/dto/state"
	"os"
	"os/exec"
	"path/filepath"
)

const ChromiumProfileDir = ".config/hypr/hyprpop/chromium/floatingChromium"

func createChromiumWindow(window *state.WindowConfig) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	profilePath := filepath.Join(home, ChromiumProfileDir)

	cmd := exec.Command("chromium",
		"--app="+window.URL,
		"--class="+window.Name,
		"--user-data-dir="+profilePath,
	)

	err = cmd.Start()
	if err != nil {
		return err
	}

	return nil
}
