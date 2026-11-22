class LoadBalancer:

    def __init__(self):
        self.index = 0

    def next(
        self,
        service
    ):

        if not service.modules:
            return None

        module = service.modules[
            self.index %
            len(service.modules)
        ]

        self.index += 1

        return module