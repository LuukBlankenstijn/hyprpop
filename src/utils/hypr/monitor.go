package hypr

import (
	"encoding/json"
	"fmt"
	"hyprpop/src/dto/state"
	"os/exec"
)

func getActiveMonitor() (*state.Monitor, error) {
	cmd := exec.Command("hyprctl", "monitors", "-j")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute hyprctl: %w", err)
	}

	var monitor []state.Monitor
	if err := json.Unmarshal(output, &monitor); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	for _, monitor := range monitor {
		if monitor.Focused {
			return &monitor, nil
		}
	}

	return nil, fmt.Errorf("no active monitor found")
}

func getMonitorByWorkspace(workspace *state.Workspace) (*state.Monitor, error) {
	cmd := exec.Command("hyprctl", "monitors", "-j")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute hyprctl: %w", err)
	}

	var monitors []state.Monitor
	if err := json.Unmarshal(output, &monitors); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	for _, monitor := range monitors {
		if monitor.Workspace == *workspace {
			return &monitor, nil
		}
	}

	return nil, fmt.Errorf("no active monitor found")
}
