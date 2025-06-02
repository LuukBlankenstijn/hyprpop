package utils

import (
	"fmt"
	"hyprpop/src/dto/state"
	"os/exec"
)

func MoveWindowToWorkspace(window *state.Window, workspaceName string, silent bool) error {
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
		workspaceName+",",
		windowId,
	)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
