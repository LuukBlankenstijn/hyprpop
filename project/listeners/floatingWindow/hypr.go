package floatingwindow

import (
	"hyprwindow/project/dto/state"
	"os"
	"os/exec"
	"path/filepath"
)

const ChromiumProfileDir = ".config/hypr/hyprwindow/chromium"

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

	cmd.Stderr = nil
	cmd.Stdout = nil

	err = cmd.Start()
	if err != nil {
		return -1, err
	}

	return cmd.Process.Pid, nil
}
