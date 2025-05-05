package floatingwindow

import (
	"hyprpop/src/dto/state"
	"os"
	"os/exec"
	"path/filepath"
)

const ChromiumProfileDir = ".config/hypr/hyprpop/chromium/floatingChromium"

func createChromiumWindow(window *state.WindowConfig) (int, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return -1, err
	}
	profilePath := filepath.Join(home, ChromiumProfileDir, window.Name)

	cmd := exec.Command("chromium",
		"--app="+window.URL,
		"--class="+window.Name,
		"--user-data-dir="+profilePath,
	)

	err = cmd.Start()
	if err != nil {
		return -1, err
	}

	return cmd.Process.Pid, nil
}
