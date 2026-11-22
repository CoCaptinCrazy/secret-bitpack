import hashlib
import time


class Layer:

    def __init__(
        self,
        name,
        instructions
    ):
        self.name = name

        self.instructions = instructions

        self.created = int(time.time())

        self.id = self._generate_id()

    def _generate_id(self):

        content = "\n".join(
            self.instructions
        )

        return hashlib.sha256(
            content.encode()
        ).hexdigest()

    def to_dict(self):
        return {
            "id": self.id,
            "name": self.name,
            "instructions": self.instructions,
            "created": self.created
        }