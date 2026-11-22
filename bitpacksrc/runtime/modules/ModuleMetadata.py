import time


class ModuleMetadata:

    def __init__(
        self,
        blueprint
    ):
        self.blueprint = blueprint

        self.created = int(
            time.time()
        )

        self.labels = {}

    def label(
        self,
        key,
        value
    ):
        self.labels[key] = value