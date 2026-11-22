from bitpack.core.events.EventBus import EventBus


class EventsCommand:

    def execute(self, args):

        bus = EventBus()

        print(
            bus.listeners
        )