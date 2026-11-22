from bitpack.core.events.EventBus import EventBus


class EventDispatcher:

    def __init__(self):
        self.bus = EventBus()

    def subscribe(self, event_name, listener):
        self.bus.subscribe(
            event_name,
            listener
        )

    def dispatch(self, event):
        listeners = self.bus.listeners_for(
            event.name
        )

        for listener in listeners:
            listener.handle(event)