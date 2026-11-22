import time


class ResourceTracker:

    def __init__(self):
        self.stats = {}

    def track(
        self,
        module_id
    ):

        self.stats[module_id] = {
            "started": time.time(),
            "cpu": 0,
            "memory": 0
        }

    def get(
        self,
        module_id
    ):
        return self.stats.get(
            module_id
        )