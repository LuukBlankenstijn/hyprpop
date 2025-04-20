package hypr

import (
	"encoding/json"
	"fmt"
	"hyprpop/src/dto/state"
	"os/exec"
)

func GetActiveWorkSpace() (*state.Workspace, error) {
	cmd := exec.Command("hyprctl", "activeworkspace", "-j")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute hyprctl: %w", err)
	}

	var workspace state.Workspace
	if err := json.Unmarshal(output, &workspace); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &workspace, nil
}

func MoveWindowToWorkspace(window *state.Window, workspace string, silent bool) error {
	windowId := fmt.Sprintf("address:%s", window.Address)
	var command string
	if silent {
		command = "movetoworkspacesilent"
	} else {
		command = "movetoworkspace"
	}
	cmd := exec.Command(
		"hyprctl",
		"dispatch",
		command,
		workspace+",",
		windowId,
	)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
