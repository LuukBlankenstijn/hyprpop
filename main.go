package main

import (
	"hyprwindow/project/core"
)

func main() {
	_, err := core.Initialize()
	if err != nil {
		panic(err)
	}

	// create window manager
	// TODO: implement window manager
}
