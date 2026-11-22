class Blueprint:

    def __init__(
        self,
        name,
        version="1.0"
    ):
        self.name = name
        self.version = version

        self.layers = []
        self.metadata = {}

    def add_layer(self, layer):
        self.layers.append(layer)

    def set_metadata(
        self,
        metadata
    ):
        self.metadata = metadata