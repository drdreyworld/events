package events

import "testing"

func TestCreateEvent(t *testing.T) {
	event := CreateEvent("event_id_1", "event_data")
	if event == nil {
		t.Error("CreateEvent returns nil")
	}

	de, ok := event.(*DefaultEvent)
	if !ok {
		t.Error("CreateEvent returns not DefaultEvent")
	}

	if de.GetID() != "event_id_1" {
		t.Error("Missed event id")
	}

	if de.GetData() != "event_data" {
		t.Error("Missed event data")
	}
}
