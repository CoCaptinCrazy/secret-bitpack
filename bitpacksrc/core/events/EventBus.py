class EventBus:
    _instance = None

    def __new__(cls):
        if cls._instance is None:
            cls._instance = super().__new__(cls)

            cls._instance.listeners = {}

        return cls._instance

    def subscribe(self, event_name, listener):
        self.listeners.setdefault(
            event_name,
            []
        ).append(listener)

    def emit(self, event):
        for listener in self.listeners.get(
            event.name,
            []
        ):
            listener.handle(event)