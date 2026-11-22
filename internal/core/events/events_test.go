package events

import "testing"

func TestPublishDispatchesToMatchingSubscribers(t *testing.T) {
	bus := New()
	called := false
	bus.Subscribe("ready", func(event Event) {
		called = event.Payload == "ok"
	})
	bus.Publish(Event{Type: "ready", Payload: "ok"})
	if !called {
		t.Fatal("event was not dispatched")
	}
}
