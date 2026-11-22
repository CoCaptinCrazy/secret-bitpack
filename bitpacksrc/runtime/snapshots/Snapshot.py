import uuid
import time


class Snapshot:

    def __init__(
        self,
        module_id,
        state
    ):
        self.id = str(uuid.uuid4())

        self.module_id = module_id

        self.created = int(
            time.time()
        )

        self.state = state

    def to_dict(self):
        return {
            "id": self.id,
            "module_id": self.module_id,
            "created": self.created,
            "state": self.state
        }