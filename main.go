package main

import (
	"hyprpop/src/core"
	"hyprpop/src/listeners/floatingWindow"
	"time"
)

func main() {
	app, err := core.Initialize()
	if err != nil {
		panic(err)
	}

	app.RegisterListener(floatingwindow.StartListening)
	time.Sleep(1 * time.Second)
	select {}
}
