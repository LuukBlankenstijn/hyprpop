package pubsub

import "fmt"

type Event struct {
	Type string
	Name string
}

func (e Event) ToString() string {
	return fmt.Sprintf("hyprpop:%s:%s", e.Type, e.Name)
}
