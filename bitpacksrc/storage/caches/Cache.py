import time


class Cache:

    def __init__(
        self,
        name,
        data=None
    ):
        self.name = name

        self.data = data or {}

        self.created = int(
            time.time()
        )

    def put(
        self,
        key,
        value
    ):
        self.data[key] = value

    def get(
        self,
        key,
        default=None
    ):
        return self.data.get(
            key,
            default
        )