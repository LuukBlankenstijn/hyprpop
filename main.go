package main

import (
	"fmt"
	"hyprpop/src/core"
	floatingwindow "hyprpop/src/listeners/floatingWindow"
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

	defer cleanup()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-c
		cleanup()
		os.Exit(0)
	}()
	select {}
}

func cleanup() {
	cleanupOnce.Do(func() {
		fmt.Println("cleanup...")
		//TODO: deregister keybinds, kill windows
	})
}
