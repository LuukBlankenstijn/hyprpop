package main

import (
	"hyprpop/src/core"
	hyprutils "hyprpop/src/hypr/utils"
	floatingwindow "hyprpop/src/listeners/floatingWindow"
	"hyprpop/src/state"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	cleanupOnce sync.Once
)

func main() {
	app, err := core.Initialize()
	if err != nil {
		panic(err)
	}

	app.RegisterListener(floatingwindow.StartListening)
	time.Sleep(1 * time.Second)

	defer cleanup(app.State)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-c
		cleanup(app.State)
		os.Exit(0)
	}()
	select {}
}

func cleanup(state *state.GlobalConfig) {
	cleanupOnce.Do(func() {
		_ = hyprutils.DeregisterAllKeybinds()
		hyprutils.KillAllWindows(*state)
	})
}
