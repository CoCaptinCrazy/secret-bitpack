from datetime import datetime


class Event:
    def __init__(self, name, data=None):
        self.name = name
        self.data = data or {}
        self.timestamp = datetime.now()

    def __str__(self):
        return (
            f"Event("
            f"name={self.name}, "
            f"data={self.data}"
            f")"
        )