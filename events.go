package events

import "sync"

func CreateEvents() Events {
	return Events{
		subscribes: map[string]map[int]Subscriber{},
	}
}

type Events struct {
	sync.Mutex
	subscribes map[string]map[int]Subscriber
}

func (e *Events) Subscribe(event string, subscriber Subscriber) {
	e.Lock()
	defer e.Unlock()

	if _, ok := e.subscribes[event]; !ok {
		e.subscribes[event] = map[int]Subscriber{}
	}

	e.subscribes[event][subscriber.SubscriberID()] = subscriber
}

func (e *Events) Unsubscribe(event string, subscriber Subscriber) {
	e.Lock()
	defer e.Unlock()

	if _, ok := e.subscribes[event]; ok {
		if _, ok := e.subscribes[event][subscriber.SubscriberID()]; ok {
			delete(e.subscribes[event], subscriber.SubscriberID())
		}
	}
}

func (e *Events) Publish(event Event) {
	e.Lock()
	defer e.Unlock()

	for id := range e.subscribes[event.GetID()] {
		e.subscribes[event.GetID()][id].Notify(event)
	}
}
