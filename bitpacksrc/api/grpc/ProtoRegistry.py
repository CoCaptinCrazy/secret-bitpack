class ProtoRegistry:

    def __init__(self):
        self.services = []

    def register(
        self,
        service
    ):
        self.services.append(
            service
        )