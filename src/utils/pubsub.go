package utils

import (
	"hyprpop/src/dto/pubsub"
	"sync"
)

type PubSub struct {
	mu          sync.RWMutex
	subscribers map[chan pubsub.Event]struct{}
}

func NewPubSub() *PubSub {
	return &PubSub{
		subscribers: make(map[chan pubsub.Event]struct{}),
	}
}

func (ps *PubSub) Subscribe(bufferSize int) chan pubsub.Event {
	ch := make(chan pubsub.Event, bufferSize)
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.subscribers[ch] = struct{}{}
	return ch
}

func (ps *PubSub) UnSubscribe(ch chan pubsub.Event) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	delete(ps.subscribers, ch)
	close(ch)
}

func (ps *PubSub) Publish(event pubsub.Event) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for ch := range ps.subscribers {
		select {
		case ch <- event:
		default:
		}
	}
}
