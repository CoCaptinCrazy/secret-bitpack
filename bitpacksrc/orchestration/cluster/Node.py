import uuid


class Node:

    def __init__(
        self,
        name,
        address
    ):
        self.id = str(uuid.uuid4())

        self.name = name

        self.address = address

        self.status = "online"

        self.capacity = {
            "cpu": 0,
            "memory": 0
        }