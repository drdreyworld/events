package events

type Event interface {
	GetID() string
	GetData() interface{}
}

func CreateEvent(id string, data interface{}) Event {
	return &DefaultEvent{
		eventId:   id,
		eventData: data,
	}
}

type DefaultEvent struct {
	eventId   string
	eventData interface{}
}

func (e *DefaultEvent) GetID() string {
	return e.eventId
}

func (e *DefaultEvent) GetData() interface{} {
	return e.eventData
}
