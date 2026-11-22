class Assembly:

    def __init__(
        self,
        name
    ):
        self.name = name
        self.modules = []

    def add_module(
        self,
        module
    ):
        self.modules.append(
            module
        )