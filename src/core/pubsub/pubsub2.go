package pubsub

import (
	"hyprpop/src/dto/pubsub"
	"hyprpop/src/logging"
	"sync"
)

var (
	globalBus *EventBus
	once      sync.Once
)

type EventBus struct {
	mu          sync.RWMutex
	subscribers map[chan pubsub.Event]struct{}
}

func Initialize() {
	once.Do(func() {
		globalBus = &EventBus{
			subscribers: make(map[chan pubsub.Event]struct{}),
		}
	})
}

func Subscribe(buffersize int) chan pubsub.Event {
	if globalBus == nil {
		logging.Fatal("pubsub not initialized", nil)
	}
	return globalBus.subscribe(buffersize)
}

func Unsubscribe(ch chan pubsub.Event) {
	if globalBus == nil {
		return
	}
	globalBus.unsubscribe(ch)
}

func Publish(event pubsub.Event) {
	if globalBus == nil {
		return // Silently ignore if not initialized
	}
	globalBus.publish(event)
}

func (eb *EventBus) subscribe(bufferSize int) chan pubsub.Event {
	ch := make(chan pubsub.Event, bufferSize)
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.subscribers[ch] = struct{}{}
	return ch
}

func (eb *EventBus) unsubscribe(ch chan pubsub.Event) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	delete(eb.subscribers, ch)
	close(ch)
}

func (eb *EventBus) publish(event pubsub.Event) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	for ch := range eb.subscribers {
		select {
		case ch <- event:
		default:
			// Non-blocking send
		}
	}
}
