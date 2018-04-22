package events

type Subscriber interface {
	SubscriberID() int
	Notify(event Event)
}
