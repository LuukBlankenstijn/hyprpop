package main

import (
	"hyprpop/src/core"
	"hyprpop/src/listeners/floatingWindow"
)

func main() {
	app, err := core.Initialize()
	if err != nil {
		panic(err)
	}

	app.RegisterListener(floatingwindow.StartListening)
	select {}
}
