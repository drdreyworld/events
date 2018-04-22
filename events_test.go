package events

import (
	"testing"
	"fmt"
	"math/rand"
	"time"
)

type testSubscriber struct {
	subscriberId int
	handledEvents []Event
}

func (ts *testSubscriber) SubscriberID() int {
	return ts.subscriberId
}

func (ts *testSubscriber) Notify(event Event) {
	ts.handledEvents = append(ts.handledEvents, event)
}

func (ts *testSubscriber) TestHandledCount(t *testing.T, expected int) {
	t.Helper()

	if len(ts.handledEvents) < expected {
		t.Errorf("Subscriber%d handled events count less than expected", ts.subscriberId)
	}

	if len(ts.handledEvents) > expected {
		t.Errorf("Subscriber%d handled events count greater than expected", ts.subscriberId)
	}
}

func (ts *testSubscriber) TestHandledEvent(t *testing.T, index int, expectedId string, expectedData string) {
	t.Helper()

	if ts.handledEvents[index].GetID() != expectedId {
		t.Errorf("Subscriber%d event with invalid id\n", ts.subscriberId)
	}

	if ts.handledEvents[index].GetData() != expectedData {
		t.Errorf("Subscriber%d event with invalid data\n", ts.subscriberId)
	}
}

func TestCreateEvents(t *testing.T) {
	subscriber_1 := &testSubscriber{}
	subscriber_1.subscriberId = 1

	subscriber_2 := &testSubscriber{}
	subscriber_2.subscriberId = 2

	subscriber_3 := &testSubscriber{}
	subscriber_3.subscriberId = 3

	events := CreateEvents()
	events.Subscribe("test_event_1", subscriber_1)
	events.Subscribe("test_event_2", subscriber_2)
	events.Subscribe("test_event_1", subscriber_3)
	events.Subscribe("test_event_2", subscriber_3)

	events.Publish(CreateEvent("test_event_1", "test_event_1_data"))

	subscriber_1.TestHandledCount(t, 1)
	subscriber_2.TestHandledCount(t, 0)
	subscriber_3.TestHandledCount(t, 1)

	events.Publish(CreateEvent("test_event_2", "test_event_2_data"))

	subscriber_1.TestHandledCount(t, 1)
	subscriber_2.TestHandledCount(t, 1)
	subscriber_3.TestHandledCount(t, 2)

	subscriber_1.TestHandledEvent(t, 0, "test_event_1", "test_event_1_data")
	subscriber_2.TestHandledEvent(t, 0, "test_event_2", "test_event_2_data")
	subscriber_3.TestHandledEvent(t, 0, "test_event_1", "test_event_1_data")
	subscriber_3.TestHandledEvent(t, 1, "test_event_2", "test_event_2_data")

	events.Unsubscribe("test_event_1", subscriber_3)

	events.Publish(CreateEvent("test_event_1", "test_event_1_data"))

	subscriber_1.TestHandledCount(t, 2)
	subscriber_2.TestHandledCount(t, 1)
	subscriber_3.TestHandledCount(t, 2)
}

func TestEvents_ParallelComplex(t *testing.T) {
	scount := 100
	subscribers := []*testSubscriber{}

	for i := 0; i < scount; i++ {
		subscribers = append(subscribers, &testSubscriber{
			subscriberId:i,
		})
	}

	events := CreateEvents()

	t.Run("SubscribeParallel", func(t *testing.T) {
		for i := 0; i < scount; i++ {
			sc := subscribers[i]
			t.Run("Subscribe", func(t *testing.T) {
				t.Parallel()
				rand.Seed(time.Now().Unix())
				events.Subscribe(fmt.Sprintf("test_event_%d", i), sc)
				events.Publish(CreateEvent(
					fmt.Sprintf("test_event_%d", rand.Intn(scount)),
					fmt.Sprintf("test_event_data_%d", i),
				))
			})

			t.Run("Publish", func(t *testing.T) {
				t.Parallel()
				for j := 0; j < rand.Intn(scount); j++ {
					events.Publish(CreateEvent(
						fmt.Sprintf("test_event_%d", rand.Intn(scount)),
						fmt.Sprintf("test_event_data_%d", rand.Intn(scount)),
					))
				}
			})

			t.Run("Unsubscribe", func(t *testing.T) {
				t.Parallel()
				events.Unsubscribe(fmt.Sprintf("test_event_%d", i), sc)
			})
		}
	})
}

func TestEvents_SubscribeParallel(t *testing.T) {
	events := CreateEvents()
	for i := 0; i < 100; i++ {
		sc := &testSubscriber{
			subscriberId:i,
		}
		t.Run("Subscribe", func(t *testing.T) {
			t.Parallel()
			events.Subscribe(fmt.Sprintf("test_event_%d", i), sc)
		})
	}
}

func TestEvents_UnsubscribeParallel(t *testing.T) {
	events := CreateEvents()
	for i := 0; i < 100; i++ {
		sc := &testSubscriber{
			subscriberId:i,
		}
		t.Run("Unsubscribe", func(t *testing.T) {
			t.Parallel()
			events.Unsubscribe(fmt.Sprintf("test_event_%d", i), sc)
		})
	}
}

func TestEvents_PublishParallel(t *testing.T) {
	events := CreateEvents()
	for i := 0; i < 100; i++ {
		events.Subscribe(fmt.Sprintf("test_event_%d", i), &testSubscriber{
			subscriberId:i,
		})

		rand.Seed(time.Now().UnixNano())
		t.Run("Publish", func(t *testing.T) {
			t.Parallel()
			for j := 0; j < rand.Intn(100); j++ {
				events.Publish(CreateEvent(
					fmt.Sprintf("test_event_%d", rand.Intn(100)),
					fmt.Sprintf("test_event_data_%d", rand.Intn(100)),
				))
			}
		})
	}
}