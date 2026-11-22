class EventListener:

    def handle(self, event):
        raise NotImplementedError(
            "Listeners must implement handle()"
        )