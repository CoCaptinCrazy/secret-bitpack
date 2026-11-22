import uuid


class Bridge:

    def __init__(
        self,
        name,
        host_port,
        target_port
    ):
        self.id = str(uuid.uuid4())

        self.name = name

        self.host_port = host_port

        self.target_port = target_port