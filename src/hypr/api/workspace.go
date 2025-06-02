package api

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
