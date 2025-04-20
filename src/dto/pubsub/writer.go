package pubsub

type PubSubWriter interface {
	Publish(event Event)
}
