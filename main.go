package main

import (
	"hyprpop/project/core"
	"hyprpop/project/listeners/floatingWindow"
)

func main() {
	app, err := core.Initialize()
	if err != nil {
		panic(err)
	}

	app.RegisterListener(floatingwindow.StartListening)
	select {}
}
