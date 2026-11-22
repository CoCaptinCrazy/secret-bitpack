class Service:

    def __init__(
        self,
        name
    ):
        self.name = name
        self.modules = []

    def add(
        self,
        module
    ):
        self.modules.append(
            module
        )