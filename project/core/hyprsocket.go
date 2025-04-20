package core

import (
	"fmt"
	"hyprpop/project/dto/pubsub"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func Listen(writer pubsub.PubSubWriter) {
	socketPath, err := getSocketPath()
	if err != nil {
		panic(err)
	}
	listenEvents(socketPath, writer)
}

func listenEvents(path string, writer pubsub.PubSubWriter) {
	for {
		// Connect to the Unix domain socket
		conn, err := net.Dial("unix", path)
		if err != nil {
			fmt.Printf("Error connecting to socket: %v\n", err)
			return
		}

		// Read from the socket
		buffer := make([]byte, 1024)
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				fmt.Printf("Error reading from socket: %v\n", err)
				err := conn.Close()
				if err != nil {
					fmt.Printf("Error closing socket: %v\n", err)
				}
				break
			}

			event, valid := parseEvent(string(buffer[:n]))
			if !valid {
				continue
			}
			writer.Publish(*event)
		}
	}
}

func parseEvent(input string) (*pubsub.Event, bool) {
	parts := strings.Split(input, ">>")
	if len(parts) < 2 || parts[0] != "custom" || !strings.HasPrefix(parts[1], "hyprpop") {
		return nil, false
	}
	eventData := strings.ReplaceAll(parts[1], "hyprpop:", "")
	eventData = strings.TrimSuffix(eventData, "custom")
	eventData = strings.TrimSuffix(eventData, "\n")
	parts = strings.Split(eventData, ":")
	if len(parts) != 2 {
		return nil, false
	}
	var e = pubsub.Event{
		Type: parts[0],
		Name: parts[1],
	}
	return &e, true
}

func getSocketPath() (string, error) {
	// Find the Hyprland runtime directory
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		runtimeDir = fmt.Sprintf("/run/user/%d", os.Getuid())
	}

	// Find the Hyprland instance signature
	hyprDir := filepath.Join(runtimeDir, "hypr")
	instanceSignature := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	if instanceSignature == "" {
		// if the env variable is not set, try to find it manually
		instances, err := os.ReadDir(hyprDir)
		if err != nil {
			return "", fmt.Errorf("error reading Hypr directory: %v", err)
		}

		if len(instances) == 0 {
			return "", fmt.Errorf("no Hyprland instances found")
		}
		instanceSignature = instances[0].Name()
	}

	return filepath.Join(hyprDir, instanceSignature, ".socket2.sock"), nil
}
