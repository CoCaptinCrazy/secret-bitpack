import time


class BlueprintMetadata:

    def __init__(
        self,
        name,
        version,
        environment
    ):
        self.name = name
        self.version = version

        self.environment = environment

        self.created = int(
            time.time()
        )

    def to_dict(self):
        return {
            "name": self.name,
            "version": self.version,
            "environment": self.environment,
            "created": self.created
        }